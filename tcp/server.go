package tcp

import (
	"context"
	"fmt"
	"github.com/StellarisJAY/mnet/interface/network"
	"github.com/StellarisJAY/mnet/route"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
)

type Server struct {
	port  int
	conns sync.Map
	proto network.Protocol

	nextConnId uint32
	ctx        context.Context
	cancel     context.CancelFunc

	listener net.Listener

	startHook network.ServerStartHook
	closeHook network.ServerCloseHook

	router network.Router

	closeChan chan interface{}
}

func MakeServer(port int, protocol network.Protocol) *Server {
	s := new(Server)
	s.port = port
	s.proto = protocol
	s.conns = sync.Map{}
	s.nextConnId = 0
	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.closeChan = make(chan interface{})
	s.router = route.MakeMapRouter()
	return s
}

func (s *Server) AddRoute(typeCode byte, handler network.Handler) {
	s.router.Register(typeCode, handler)
}

func (s *Server) Start() (error, chan interface{}) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(s.port))
	if err != nil {
		return fmt.Errorf("start server error: %v", err), nil
	}
	s.listener = listener

	// listen to os signals
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	go func() {
		sig := <-signals
		switch sig {
		case syscall.SIGHUP, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT:
			// use context.CancelFunc to stop server
			s.cancel()
		}
	}()

	go func() {
		for {
			conn, acceptErr := s.listener.Accept()
			if acceptErr != nil {
				break
			}
			// accepted connection, generate conn ID and bind protocol
			connection := MakeTcpConnection(conn, atomic.AddUint32(&s.nextConnId, 1), s.proto, s.router, false)
			s.conns.Store(connection, true)
			// start connection's IO loop
			go func() {
				connection.Start()
				s.conns.Delete(connection)
			}()
		}
	}()
	// call start hook after starting server
	if s.startHook != nil {
		s.startHook()
	}
	log.Println("TCP server started, listening: ", s.port)

	s.router.StartWorkers()
	go func() {
		// wait for close signal
		select {
		case <-s.ctx.Done():
			s.router.Close()
			s.gracefulClose()
		}
	}()
	return nil, s.closeChan
}

// Shutdown server gracefully
func (s *Server) gracefulClose() {
	log.Println("shutting down tcp server...")
	// call shutdown hook
	if s.closeHook != nil {
		s.closeHook()
	}
	// close all client connections
	s.conns.Range(func(key, value interface{}) bool {
		connection := key.(*Connection)
		connection.Close()
		return true
	})
	// close listener
	_ = s.listener.Close()
	s.closeChan <- 1
}

func (s *Server) Close() {
	s.cancel()
}

func (s *Server) Protocol() network.Protocol {
	return s.proto
}

func (s *Server) SetProtocol(p network.Protocol) {
	s.proto = p
}

func (s *Server) SetServerStartHook(hook network.ServerStartHook) {
	s.startHook = hook
}

func (s *Server) SetServerCloseHook(hook network.ServerCloseHook) {
	s.closeHook = hook
}
