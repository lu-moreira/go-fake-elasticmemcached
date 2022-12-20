package pkg

import (
	"fmt"
	"net"
	"strings"
)

type Commander interface {
	Execute(conn net.Conn, data []byte) (shouldQuit bool)
}

type commander struct {
	n NodeManager
}

func NewCommander(n NodeManager) Commander {
	return &commander{
		n: n,
	}
}

func (c commander) Execute(conn net.Conn, data []byte) (shouldQuit bool) {
	cmd := strings.TrimSpace(string(data))

	switch cmd {
	case "config get cluster", "get AmazonElasticCache:cluster":
		conn.Write(c.opClusterConfigResponse(c.n.Nodes(), fmt.Sprint(c.n.Version())))
	case "quit":
		conn.Close()
		shouldQuit = true
	}
	return
}

func (c commander) opClusterConfigResponse(nodes []string, version string) []byte {
	nodesStr := strings.Join(nodes, " ")
	sb := strings.Builder{}
	_, _ = sb.WriteString(fmt.Sprintf("CONFIG cluster 0 %d\r\n", len(nodesStr)))
	_, _ = sb.WriteString(fmt.Sprintf("%s\n", version))
	_, _ = sb.WriteString(fmt.Sprintf("%s\n\r\n", nodesStr))
	_, _ = sb.WriteString("END\r\n")

	return []byte(sb.String())
}
