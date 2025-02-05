package cluster

import (
	"errors"
	"io"
	"time"

	"github.com/hashicorp/memberlist"
	"github.com/stathat/consistent"
)

// Node is the interface that wraps the basic methods of a gossip node in the cluster.
type Node interface {
	ShouldProcess(key string) (string, bool)
	MemberList() []string
	NodeAddr() string
}

// gossipNode is the basic struct of a gossip node in the cluster, which implements the Node interface.
type gossipNode struct {
	hashRing *consistent.Consistent
	addr     string
}

func (gnd *gossipNode) ShouldProcess(key string) (string, bool) {
	addr, _ := gnd.hashRing.Get(key)
	return addr, addr == gnd.addr
}

func (gnd *gossipNode) MemberList() []string {
	return gnd.hashRing.Members()
}

func (gnd *gossipNode) NodeAddr() string {
	return gnd.addr
}

// New creates a new gossip node in the cluster.
// addr is the address of the node, and cluster is the address of the cluster.
func New(addr, cluster string) (Node, error) {
	if addr == "" {
		return nil, errors.New("addr is empty")
	}
	if cluster == "" {
		cluster = addr
	}

	conf := memberlist.DefaultLANConfig()
	conf.Name = addr
	conf.BindAddr = addr
	conf.LogOutput = io.Discard
	lst, err := memberlist.Create(conf)
	if err != nil {
		return nil, err
	}
	// memberlist automatically handles information transfer and status synchronization between nodes.
	_, err = lst.Join([]string{cluster})
	if err != nil {
		return nil, err
	}
	hashRing := consistent.New()
	hashRing.NumberOfReplicas = 256 // default is 20

	go func() {
		for {
			members := lst.Members()
			nodes := make([]string, len(members))
			for i, member := range members {
				nodes[i] = member.Name
			}
			hashRing.Set(nodes)
			time.Sleep(3 * time.Second)
		}
	}()

	return &gossipNode{hashRing: hashRing, addr: addr}, nil
}
