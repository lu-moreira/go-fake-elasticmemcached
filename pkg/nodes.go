package pkg

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

type NodeManager interface {
	// Start and Stop, maybe need to be in a different contract
	Start(n int) error
	Stop() error
	//
	Nodes() []string
	Version() int
}

type dockerNodeManager struct {
	containers []string
	nodes      []string
	version    int
}

func NewDockerNodeManager() NodeManager {
	return &dockerNodeManager{
		containers: make([]string, 0),
		nodes:      make([]string, 0),
		version:    0,
	}
}

func (nm *dockerNodeManager) Start(numberOfNodes int) error {
	log.Printf("starting %d memcached docker containers\n", numberOfNodes)

	curr_nodes := len(nm.containers) // will serve to us to determine the port

	// new containers and nodes
	containers := make([]string, numberOfNodes)
	nodes := make([]string, numberOfNodes)
	for idx := 0; idx < numberOfNodes; idx++ {
		remotePort := fmt.Sprintf("112%d", idx+curr_nodes)
		cmd := exec.Command("docker", "run", "-d", "-p", fmt.Sprintf("%s:11211", remotePort), "memcached")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return err
		}
		if err := cmd.Start(); err != nil {
			return err
		}

		res, err := ioutil.ReadAll(stdout)
		if err != nil {
			return err
		}
		if err := cmd.Wait(); err != nil {
			return err
		}

		containers[idx] = string([]rune(strings.TrimSpace(string(res)))[:12]) // we don't need the full hash
		nodes[idx] = fmt.Sprintf("localhost|127.0.0.1|%s", remotePort)
		log.Printf("created memcached container with hash id %s and addr %s\n", containers[idx], fmt.Sprintf("127.0.0.1:%s", remotePort))
	}

	nm.containers = append(nm.containers, containers...)
	nm.nodes = append(nm.nodes, nodes...)
	nm.version++ // update the version
	return nil
}

func (nm dockerNodeManager) Stop() error {
	stopCmd := append([]string{"stop"}, nm.containers...)
	log.Printf("stopping memcached containers with command: docker %s\n", strings.Join(stopCmd, " "))

	cmd := exec.Command("docker", stopCmd...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}

	_, err = ioutil.ReadAll(stdout)
	if err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}

	log.Println("containers stopped")
	return nil
}

func (nm dockerNodeManager) Nodes() []string {
	return nm.nodes
}

func (nm dockerNodeManager) Version() int {
	return nm.version
}
