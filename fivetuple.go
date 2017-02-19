package turn

import (
	"net"
	"time"
)

type Address struct {
	IP   net.IP
	Port int
}

// FiveTuple represents 5-tuple (client's IP address, client's port, server IP
// address, server port, transport protocol) from RFC.
//
// On the  client, the 5-tuple uses the client's host transport address; on the
// server, the 5-tuple uses the client's server-reflexive transport
// address.
type FiveTuple struct {
	Client   Address  // host on client and reflexive on server
	Server   Address  // same on client and server
	Protocol Protocol // note: not same as REQUESTED-TRANSPORT
}

type Allocation struct {
	FiveTuple FiveTuple
	Key       []byte
	ExpiresAt time.Time
}
