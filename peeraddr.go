package turn

import (
	"net"

	"github.com/ernado/stun"
)

// PeerAddress implements XOR-PEER-ADDRESS attribute.
//
// The XOR-PEER-ADDRESS specifies the address and port of the peer as
// seen from the TURN server. (For example, the peer's server-reflexive
// transport address if the peer is behind a NAT.)
//
// https://trac.tools.ietf.org/html/rfc5766#section-14.3
type PeerAddress struct {
	IP   net.IP
	Port int
}

func (a PeerAddress) String() string {
	return stun.XORMappedAddress(a).String()
}

// AddTo adds XOR-PEER-ADDRESS to message.
func (a PeerAddress) AddTo(m *stun.Message) error {
	return (stun.XORMappedAddress)(a).AddToAs(m, stun.AttrXORPeerAddress)
}

// AddTo decodes XOR-PEER-ADDRESS from message.
func (a *PeerAddress) GetFrom(m *stun.Message) error {
	return (*stun.XORMappedAddress)(a).GetFromAs(m, stun.AttrXORPeerAddress)
}
