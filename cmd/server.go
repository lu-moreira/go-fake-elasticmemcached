package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"

	"github.com/lu-moreira/go-fake-elasticmemcached/pkg"
)

type Server interface {
	Start() error
	Stop()
	SetupShutDown() chan struct{}
}

type server struct {
	addr     string
	listener net.Listener
	quit     chan struct{}
	wg       sync.WaitGroup
	cmder    pkg.Commander
	nm       pkg.NodeManager
}

func NewServer(addr string, cmder pkg.Commander, nm pkg.NodeManager) Server {
	return &server{
		addr:  addr,
		quit:  make(chan struct{}),
		cmder: cmder,
		nm:    nm,
	}
}

func (s *server) Start() error {
	l, err := net.Listen(SERVER_NETWORK, s.addr)
	if err != nil {
		return err
	}
	s.listener = l

	s.wg.Add(1)
	go s.startListener()

	log.Printf("server listening at %s\n", s.addr)
	return nil
}

// Stop is used by SetupShutDown() but is exported for clients have a way to force it
func (s *server) Stop() {
	close(s.quit)
	s.listener.Close()
	s.wg.Wait()
}

func (s *server) SetupShutDown() chan struct{} {
	out := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint // We wait...

		// close connections
		s.Stop()

		// close all nodes
		if err := s.nm.Stop(); err != nil {
			log.Printf("node manager shutdown Error: %v", err)
		}
		close(out)
	}()
	return out
}

func (s *server) startListener() {
	defer s.wg.Done()

	for {
		select {
		case <-s.quit:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				select {
				case <-s.quit:
					return
				default:
					log.Println("accept error", err)
					return
				}
			}
			s.wg.Add(1)
			go func() {
				s.handleConnection(conn)
				s.wg.Done()
			}()
		}
	}
}

func (s *server) handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		select {
		case <-s.quit:
			return
		default:
			// read operation
			data, err := s.readOp(conn)
			if err != nil {
				log.Println(err)
				return
			}

			// execute action on data
			shouldQuit := s.cmder.Execute(conn, data)
			if shouldQuit {
				return
			}
		}
	}
}

func (s *server) readOp(conn net.Conn) ([]byte, error) {
	data, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return data, nil
}
