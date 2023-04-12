package obfsvpn

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
)

// Obfs4Config is an obfs4 specific configuration object
type Obfs4Config struct {
	NodeID     string `json:"node-id"`
	PrivateKey string `json:"private-key"`
	PublicKey  string `json:"public-key"`
	DRBGSeed   string `json:"drbg-seed"`
	IatMode    int    `json:"iat-mode"`
}

// ServerConfig is the configuration for the obfsvpn server
type ServerConfig struct {
	OpenvpnAddr string
	Obfs4Config Obfs4Config
	StateDir    string
	// Obfs4Network is either "tcp" or "kcp"
	Obfs4Network    string
	Obfs4ListenAddr string
	PortSeed        int64
	PortCount       uint
}

// HopServer is a obfsvpn server
type HopServer struct {
	cfg    ServerConfig
	logger *log.Logger
	debug  *log.Logger
	ctx    context.Context
	stop   context.CancelFunc
}

// NewHopServer returns a new HopServer
func NewHopServer(ctx context.Context, stop context.CancelFunc, cfg ServerConfig, logger, debug *log.Logger) *HopServer {
	return &HopServer{
		ctx:    ctx,
		stop:   stop,
		cfg:    cfg,
		debug:  debug,
		logger: logger,
	}
}

// Start starts the obfsvpn server
func (s *HopServer) Start() error {

	// TODO pass a "mode" ? (kcp)
	listenConfig, err := NewListenConfig(
		s.cfg.Obfs4Config.NodeID, s.cfg.Obfs4Config.PrivateKey, s.cfg.Obfs4Config.PublicKey,
		s.cfg.Obfs4Config.DRBGSeed,
		s.cfg.StateDir,
	)
	if err != nil {
		return fmt.Errorf("error creating listener from config: %w", err)
	}

	// We want a non-crypto RNG so that we can share a seed
	// #nosec G404
	r := rand.New(rand.NewSource(s.cfg.PortSeed))

	s.debug.Printf("DEBUG: %v", listenConfig)

	listeners := make([]net.Listener, s.cfg.PortCount)

	for i := 0; i < int(s.cfg.PortCount); i++ {
		portOffset := r.Intn(PortHopRange)
		addr := net.JoinHostPort(s.cfg.Obfs4ListenAddr, fmt.Sprint(portOffset+MinHopPort))
		listeners[i], err = listenConfig.Listen(s.ctx, s.cfg.Obfs4Network, addr)

		if err != nil {
			s.logger.Printf("Error binding to %s: %v", s.cfg.Obfs4ListenAddr, err)
		}
	}

	for _, ln := range listeners {
		go s.acceptLoop(ln)
	}

	<-s.ctx.Done()
	// Stop releases the signal handling and falls back to the default behavior,
	// so sending another interrupt will immediately terminate.
	s.stop()
	s.logger.Printf("shutting down…")
	for _, ln := range listeners {
		err := ln.Close()
		if err != nil {
			s.logger.Printf("error closing listener: %v", err)
		}
	}
	return nil
}

func (s *HopServer) acceptLoop(ln net.Listener) {
	s.logger.Printf("Listening on %s…", ln.Addr())

	for {
		conn, err := ln.Accept()
		if err != nil {
			s.debug.Printf("error accepting connection: %v", err)
			continue
		}
		s.debug.Printf("accepted connection from %v", conn.RemoteAddr().String())

		udpRemote, err := net.ResolveUDPAddr("udp", s.cfg.OpenvpnAddr)
		if err != nil {
			s.logger.Printf("Error binding to %s: %v", s.cfg.OpenvpnAddr, err)
		}

		udpConn, err := net.DialUDP("udp", nil, udpRemote)
		if err != nil {
			s.logger.Printf("error dialing to %s: %v", udpRemote, err)
		}

		go readTCPWriteUDP(conn, udpConn, s.debug, s.logger)

		go readUDPWriteTCP(conn, udpConn, s.debug, s.logger)
	}
}

func readTCPWriteUDP(tcpConn net.Conn, udpConn *net.UDPConn, debug, logger *log.Logger) error {
	lengthBuffer := make([]byte, 2)
	datagramBuffer := make([]byte, MaxUDPLen)
	for {
		udpBuffer, err := ReadTCPFrameUDP(tcpConn, datagramBuffer, lengthBuffer)
		if err != nil {
			debug.Printf("Reading/framing error: %v\n", err)
			break
		}

		_, err = udpConn.Write(udpBuffer)
		if err != nil {
			debug.Printf("Write err  %v\n", err)
			break
		}
	}
	return nil
}

func readUDPWriteTCP(tcpConn net.Conn, udpConn *net.UDPConn, debug, logger *log.Logger) error {
	datagramBuffer := make([]byte, MaxUDPLen)
	for {
		tcpBuffer, _, err := ReadUDPFrameTCP(udpConn, datagramBuffer)
		if err != nil {
			debug.Printf("Reading/framing error: %v", err)
			break
		}
		_, err = tcpConn.Write(tcpBuffer)
		if err != nil {
			debug.Printf("Write err on %v to %v: %v\n", tcpConn.LocalAddr(), tcpConn.RemoteAddr(), err)
			break
		}
	}

	return nil
}
