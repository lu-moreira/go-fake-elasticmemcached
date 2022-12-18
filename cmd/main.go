package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/lu-moreira/go-fake-elasticmemcached/pkg"
)

const (
	SERVER_NETWORK = "tcp"
	SERVER_HOST    = "localhost"
	SERVER_PORT    = "11210"
)

func main() {
	numnodes := flag.Int("numnodes", 1, "Number of nodes")
	flag.Parse()

	nm := pkg.NewDockerNodeManager()
	err := nm.Start(*numnodes)
	if err != nil {
		log.Fatal(err)
	}

	addr := fmt.Sprintf("%s:%s", SERVER_HOST, SERVER_PORT)
	cmder := pkg.NewCommander(nm)
	s := NewServer(addr, cmder, nm)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

	shutchan := s.SetupShutDown()
	<-shutchan
	log.Println("done")
}
