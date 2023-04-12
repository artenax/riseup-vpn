package obfsvpn

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	pt "git.torproject.org/pluggable-transports/goptlib.git"
)

const transportName = "obfs4"

// Server is a obfsvpn server
type Server struct {
	cfg    ServerConfig
	logger *log.Logger
	debug  *log.Logger
	ctx    context.Context
	stop   context.CancelFunc
}

// NewServer returns a new Server
func NewServer(ctx context.Context, stop context.CancelFunc, cfg ServerConfig, logger, debug *log.Logger) *Server {
	return &Server{
		ctx:    ctx,
		stop:   stop,
		cfg:    cfg,
		debug:  debug,
		logger: logger,
	}
}

// Start starts the obfsvpn server
func (s *Server) Start() error {

	var network string
	if os.Getenv("KCP") == "1" {
		network = "kcp"
	} else {
		network = "tcp"
	}

	listenConfig, err := NewListenConfig(
		s.cfg.Obfs4Config.NodeID, s.cfg.Obfs4Config.PrivateKey, s.cfg.Obfs4Config.PublicKey,
		s.cfg.Obfs4Config.DRBGSeed,
		s.cfg.StateDir,
	)

	if err != nil {
		s.logger.Fatalf("Error creating listener from config: %v", err)
	}

	ln, err := listenConfig.Listen(s.ctx, network, s.cfg.Obfs4ListenAddr)

	if err != nil {
		s.logger.Fatalf("error binding to %s: %v", s.cfg.Obfs4ListenAddr, err)
	}

	go func() {
		<-s.ctx.Done()
		// Stop releases the signal handling and falls back to the default behavior,
		// so sending another interrupt will immediately terminate.
		s.stop()
		s.logger.Printf("shutting down…")
		err := ln.Close()
		if err != nil {
			s.logger.Printf("error closing listener: %v", err)
		}
	}()

	orAddr, err := net.ResolveTCPAddr("tcp", s.cfg.OpenvpnAddr)
	if err != nil {
		s.logger.Fatalf("error resolving TCP addr %s: %v", s.cfg.Obfs4ListenAddr, err)
	}

	info := &pt.ServerInfo{
		OrAddr: orAddr,
	}

	s.logger.Printf("Listening on %s…", ln.Addr())

	for {
		conn, err := ln.Accept()
		if err != nil {
			s.debug.Printf("error accepting connection: %v", err)
			continue
		}
		s.debug.Printf("accepted connection from %v", conn.RemoteAddr())
		go proxyConn(s.ctx, info, conn, s.logger, s.debug)
	}
}

// proxyConn is a connection to the obfs4 client that we have accepted. we will dial to the remote contained in info
func proxyConn(ctx context.Context, info *pt.ServerInfo, obfsConn net.Conn, logger, debug *log.Logger) {
	defer func() {
		err := obfsConn.Close()
		if err != nil {
			debug.Printf("Error closing connection: %v", err)
		}
	}()

	// FIXME scrub ips in other than debug mode!
	log.Println("Dialing:", info.OrAddr)
	log.Println("Obfs4 client:", obfsConn.RemoteAddr().String())

	/*
		    TODO(atanarjuat):

		    in the case of Tor, pt.DialOr returns a *net.TCPConn after dialing info.OrAddr.
		    in the vpn case (or any transparent proxy really), we do use
		    the pt.DialOr method to simply get a dialer to our upstream VPN remote.

		    keeping this terminology is a bit stupid and slightly confusing, instead
		    we could get the clearConn just by doing:

		    s, err := net.DialTCP("tcp", nil, info.ExtendedOrAddr)
		    if err != nil {
			    return nil, err
		    }

		    that is precisely what the code in ptlib is doing.

		    We also aspire at being a generic PT at some point, so perhaps it's better to keep the usage
		    of DialOr. Maybe not, and so we don't need to keep the confusing info struct.

		    Another argument in favor of doing our own dialing in here is that
		    we could have an UDP dialer. We need to think of a good way
		    to configure the server to distinguish between udp and tcp
		    upstream flows.
	*/

	// we'll refer to the connection to the usptream node as "clearConn", as opposed to the obfuscated conn.
	// but for sure openvpn or whatever protocol you wrap has its own layer of encryption :)

	clearConn, err := pt.DialOr(info, obfsConn.RemoteAddr().String(), transportName)

	if err != nil {
		logger.Printf("error dialing remote: %v", err)
		return
	}

	if err = CopyLoop(clearConn, obfsConn); err != nil {
		debug.Printf("%s - closed connection: %s", "obfsvpn", err.Error())
	} else {
		debug.Printf("%s - closed connection", "obfsvpn")
	}
}

// CopyLoop is a standard copy loop. We don't care too much who's client and
// who's server
func CopyLoop(left net.Conn, right net.Conn) error {

	fmt.Println("--> Entering copy loop.")

	if left == nil {
		fmt.Fprintln(os.Stderr, "--> Copy loop has a nil connection (left).")
		return errors.New("copy loop has a nil connection (left)")
	}

	if right == nil {
		fmt.Fprintln(os.Stderr, "--> Copy loop has a nil connection (right).")
		return errors.New("copy loop has a nil connection (right)")
	}

	// Note: right is always the pt connection.
	lockL := make(chan bool)
	lockR := make(chan bool)
	errChan := make(chan error)

	go CopyLeftToRight(left, right, lockL, errChan)
	go CopyRightToLeft(left, right, lockR, errChan)

	leftUp := true
	rightUp := true

	var copyErr error

	for leftUp || rightUp {
		select {
		case <-lockL:
			leftUp = false
		case <-lockR:
			rightUp = false
		case copyErr = <-errChan:
			log.Println("Error while copying")
		}
	}

	// XXX better to defer?
	err := left.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error closing left connection: ", err.Error())
	}
	err = right.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error closing right connection: ", err.Error())
	}

	return copyErr
}

// TODO check for data races

func CopyLeftToRight(l net.Conn, r net.Conn, ll chan bool, errChan chan error) {
	_, e := io.Copy(r, l)
	ll <- true
	if e != nil {
		errChan <- e
	}
}

func CopyRightToLeft(l net.Conn, r net.Conn, lr chan bool, errChan chan error) {
	_, e := io.Copy(l, r)
	lr <- true
	if e != nil {
		errChan <- e
	}
}
