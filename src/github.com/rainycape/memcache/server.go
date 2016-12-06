package memcache

import (
	"github.com/consistent"
	"net"
	"strings"
	//"fmt"
)

// Servers is the interface used to manage a set of memcached
// servers.
//
// Implementations must be safely accessible from multiple
// goroutines.
type Servers interface {
	// PickServer selects one server from the ones by
	// managed by the Servers instance, based on the
	// given key.
	PickServer(key string) (*Addr, error)
	// Servers returns all the servers managed by the
	// Servers instance.
	Servers() ([]*Addr, error)
}

// ServerList is an implementation of the Servers interface.
// To initialize a ServerList use NewServerList.
type ServerList struct {
	addrs []*Addr
	chash *consistent.Consistent
}

// NewServerList returns a new ServerList with the given servers.
// All servers have the same weight. To give a server more weight,
// list it multiple times.
//
// NewServerList returns an error if any of the received addresses
// is not valid or fails to resolve, but it doesn't try to connect
// to the provided servers.
func NewServerList(servers ...string) (*ServerList, error) {

	// Create new consistent hash map
	chash := consistent.New()

	addrs := make([]*Addr, len(servers))
	for i, server := range servers {
		// Add server to consistent hash map
		//fmt.Println("Added server ", server, " to consistent hash")
		if strings.Contains(server, "/") {
			addr, err := net.ResolveUnixAddr("unix", server)
			if err != nil {
				return nil, err
			}
			addrs[i] = NewAddr(addr)
			chash.Add(addr.String())
		} else {
			tcpaddr, err := net.ResolveTCPAddr("tcp", server)
			if err != nil {
				return nil, err
			}
			addrs[i] = NewAddr(tcpaddr)
			chash.Add(tcpaddr.String())
		}
	}
	return &ServerList{addrs: addrs, chash: chash}, nil
}

func (s *ServerList) PickServer(key string) (*Addr, error) {
	if len(s.addrs) == 0 {
		return nil, ErrNoServers
	}

	server, err := s.chash.Get(key)
	if err != nil {
		return nil, err
	}

	server_index := 0

	for key, val := range s.addrs {
		if val.String() == server {
			server_index = key
		}
	}

	return s.addrs[server_index], err
}

func (s *ServerList) Servers() ([]*Addr, error) {
	return s.addrs, nil
}
