// Package client exposes a socks5 proxy that uses obfs4 to communicate with the server,
// with an optional kcp wire transport.
package client

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	"0xacab.org/leap/obfsvpn"
	"github.com/xtaci/kcp-go"
)

type clientState string

const (
	starting clientState = "STARTING"
	running  clientState = "RUNNING"
	stopping clientState = "STOPPING"
	stopped  clientState = "STOPPED"
)

type Obfs4Config struct {
	Obfs4Cert   string
	Obfs4Remote string
}

type HopClient struct {
	kcp             bool
	ProxyAddr       string
	newObfs4Conn    chan net.Conn
	obfs4Conns      []net.Conn
	obfs4Endpoints  []*Obfs4Config
	obfs4Dialer     *obfsvpn.Dialer
	EventLogger     EventLogger
	state           clientState
	ctx             context.Context
	stop            context.CancelFunc
	openvpnConn     *net.UDPConn
	openvpnAddr     *net.UDPAddr
	openvpnAddrLock sync.RWMutex
	outLock         sync.Mutex
	minHopSeconds   uint
	hopJitter       uint
}

func NewHopClient(ctx context.Context, stop context.CancelFunc, kcp bool, proxyAddr string, obfs4Endpoints []*Obfs4Config, minHopSeconds uint, hopJitter uint) *HopClient {
	return &HopClient{
		kcp:            kcp,
		ProxyAddr:      proxyAddr,
		obfs4Endpoints: obfs4Endpoints,
		newObfs4Conn:   make(chan net.Conn),
		minHopSeconds:  minHopSeconds,
		hopJitter:      hopJitter,
		ctx:            ctx,
		stop:           stop,
		state:          stopped,
	}
}

func (c *HopClient) Start() (bool, error) {
	defer func() {
		c.state = stopped
		c.log("Start function ended")
	}()

	if c.IsStarted() {
		c.error("Cannot start proxy server, already running")
		return false, ErrAlreadyRunning
	}

	c.state = starting

	var err error

	obfs4Endpoint := c.obfs4Endpoints[0]

	c.obfs4Dialer, err = obfsvpn.NewDialerFromCert(obfs4Endpoint.Obfs4Cert)
	if err != nil {
		return false, fmt.Errorf("could not dial obfs4 remote: %w", err)
	}

	if c.kcp {
		c.obfs4Dialer.DialFunc = func(network, address string) (net.Conn, error) {
			c.log("Dialing kcp://%s", address)
			return kcp.Dial(address)
		}
	}

	obfs4Conn, err := c.obfs4Dialer.Dial("tcp", obfs4Endpoint.Obfs4Remote)
	if err != nil {
		c.error("Could not dial obfs4 remote: %v", err)
	}

	c.obfs4Conns = []net.Conn{obfs4Conn}

	// We want a non-crypto RNG so that we can share a seed
	// #nosec G404
	rand.Seed(time.Now().UnixNano())

	c.state = running

	proxyAddr, err := net.ResolveUDPAddr("udp", c.ProxyAddr)
	if err != nil {
		return false, fmt.Errorf("cannot resolve UDP addr: %w", err)
	}

	c.openvpnConn, err = net.ListenUDP("udp", proxyAddr)
	if err != nil {
		return false, fmt.Errorf("error accepting udp connection: %w", err)
	}

	go c.hop()

	go c.readUDPWriteTCP()

	go c.readTCPWriteUDP()

	<-c.ctx.Done()

	return true, nil
}

func (c *HopClient) hop() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		// #nosec G404
		sleepSeconds := rand.Intn(int(c.hopJitter)) + int(c.minHopSeconds)

		time.Sleep(time.Duration(sleepSeconds) * time.Second)

		// #nosec G404
		i := rand.Intn(len(c.obfs4Endpoints))
		obfs4Endpoint := c.obfs4Endpoints[i]

		host, port, err := net.SplitHostPort(obfs4Endpoint.Obfs4Remote)
		if err != nil {
			c.error("Could not split obfs4 remote: %v", err)
			continue
		}
		remoteAddrs, err := net.DefaultResolver.LookupHost(c.ctx, host)
		if err != nil {
			c.error("Could not lookup obfs4 remote: %v", err)
			continue
		}

		if len(remoteAddrs) <= 0 {
			c.error("Could not lookup obfs4 remote: %v", err)
			continue
		}

		newRemote := net.JoinHostPort(remoteAddrs[0], port)

		c.log("HOPPING to %+v", newRemote)

		obfs4Dialer, err := obfsvpn.NewDialerFromCert(obfs4Endpoint.Obfs4Cert)
		if err != nil {
			c.error("Could not dial obfs4 remote: %v", err)
			return
		}

		if c.kcp {
			obfs4Dialer.DialFunc = func(network, address string) (net.Conn, error) {
				c.log("Dialing kcp://%s", address)
				return kcp.Dial(address)
			}
		}

		c.log("Dialing new remote: %v", newRemote)
		newObfs4Conn, err := obfs4Dialer.Dial("tcp", newRemote)
		if err != nil {
			c.error("Could not dial obfs4 remote: %v", err)
		}
		c.log("Dialed new remote")

		c.outLock.Lock()
		c.obfs4Conns = append([]net.Conn{newObfs4Conn}, c.obfs4Conns...)
		c.outLock.Unlock()

		c.newObfs4Conn <- newObfs4Conn

		// If we wait sleepSeconds here to clean up the previous connection, we can guarantee that the
		// connection list will not grow unbounded
		go func() {
			time.Sleep(time.Duration(sleepSeconds) * time.Second)

			c.cleanupOldConn()
		}()
	}
}

func (c *HopClient) cleanupOldConn() {

	if len(c.obfs4Conns) > 1 {
		c.outLock.Lock()
		defer c.outLock.Unlock()
		connToClose := c.obfs4Conns[len(c.obfs4Conns)-1]
		c.log("Cleaning up old connection to %v", connToClose.RemoteAddr())

		err := connToClose.Close()
		if err != nil {
			c.log("Error closing obfs4 connection to %v: %v", connToClose.RemoteAddr(), err)
		}

		// Remove the connection from our tracking list
		c.obfs4Conns = c.obfs4Conns[:len(c.obfs4Conns)-1]
	}
}

func (c *HopClient) readUDPWriteTCP() {
	datagramBuffer := make([]byte, obfsvpn.MaxUDPLen)
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		tcpBuffer, newOpenvpnAddr, err := obfsvpn.ReadUDPFrameTCP(c.openvpnConn, datagramBuffer)
		if err != nil {
			c.error("Read err from %v: %v", c.openvpnConn.LocalAddr(), err)
			continue
		}

		if newOpenvpnAddr != c.openvpnAddr {
			c.openvpnAddrLock.Lock()
			c.openvpnAddr = newOpenvpnAddr
			c.openvpnAddrLock.Unlock()
		}

		func() {
			// Always write to the first connection in our list because it will be most up to date
			func() {
				c.outLock.Lock()
				defer c.outLock.Unlock()
				_, err := c.obfs4Conns[0].Write(tcpBuffer)
				if err != nil {
					c.log("Write err from %v to %v: %v", c.obfs4Conns[0].LocalAddr(), c.obfs4Conns[0].RemoteAddr(), err)
					return
				}
			}()
		}()
	}
}

func (c *HopClient) readTCPWriteUDP() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		fromTCP := make(chan []byte, 2048)

		handleObfs4Conn := func(conn net.Conn) {
			datagramBuffer := make([]byte, obfsvpn.MaxUDPLen)
			lengthBuffer := make([]byte, 2)
			for {
				udpBuffer, err := obfsvpn.ReadTCPFrameUDP(conn, datagramBuffer, lengthBuffer)
				if err != nil {
					c.error("Reading/framing error: %v", err)
					return
				}

				fromTCP <- udpBuffer
			}
		}

		go func() {
			for {
				newObfs4Conn := <-c.newObfs4Conn

				go handleObfs4Conn(newObfs4Conn)
			}
		}()

		go handleObfs4Conn(c.obfs4Conns[0])

		for {
			tcpBytes := <-fromTCP

			c.openvpnAddrLock.RLock()
			_, err := c.openvpnConn.WriteToUDP(tcpBytes, c.openvpnAddr)
			c.openvpnAddrLock.RUnlock()
			if err != nil {
				c.error("Write err from %v to %v: %v", c.openvpnConn.LocalAddr(), c.openvpnConn.RemoteAddr(), err)
				return
			}
		}
	}
}

func (c *HopClient) Close() error {
	// TODO: implement
	return nil
}

func (c *HopClient) Stop() (bool, error) {
	if !c.IsStarted() {
		return false, ErrNotRunning
	}

	if err := c.Close(); err != nil {
		c.error("error while stopping: %v", err)
		return false, err
	}

	c.state = stopped

	return true, nil
}

func (c *HopClient) log(format string, a ...interface{}) {
	if c.EventLogger != nil {
		c.EventLogger.Log(string(c.state), fmt.Sprintf(format, a...))
		return
	}
	if format == "" {
		log.Println(a...)
		return
	}
	log.Printf(format+"\n", a...)
}

func (c *HopClient) error(format string, a ...interface{}) {
	if c.EventLogger != nil {
		c.EventLogger.Error(fmt.Sprintf(format, a...))
		return
	}
	if format == "" {
		log.Println(a...)
		return
	}
	log.Printf(format+"\n", a...)
}

func (c *HopClient) IsStarted() bool {
	return c.state != stopped
}
