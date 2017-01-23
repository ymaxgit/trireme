package enforcer

import (
	"syscall"

	"github.com/aporeto-inc/trireme/crypto"
)

// AuthInfo keeps authentication information about a connection
type AuthInfo struct {
	LocalContext    []byte
	RemoteContext   []byte
	LocalContextID  string
	RemoteContextID string
	RemotePublicKey interface{}
	RemoteIP        string
	RemotePort      string
}

// initAuthInfo creates the authentication information for a connection
func initAuthInfo(s *AuthInfo) {

	nonse, _ := crypto.GenerateRandomBytes(32)
	s.LocalContext = nonse
}

// TCPConnection is information regarding TCP Connection
type TCPConnection struct {
	State TCPFlowState
	Auth  AuthInfo
}

// NewTCPConnection returns a TCPConnection information struct
func NewTCPConnection() *TCPConnection {

	c := &TCPConnection{
		State: TCPSynSend,
	}
	initAuthInfo(&c.Auth)
	return c
}

// UDPConnection stores information about a UDP connection
type UDPConnection struct {
	state   UDPFlowState
	Auth    AuthInfo
	addr    *syscall.SockaddrInet4
	packets [][]byte
}

// NewUDPConnection returns a UDPConnection information struct
func NewUDPConnection(dip []byte, dport uint16) *UDPConnection {

	c := &UDPConnection{
		packets: [][]byte{},
		addr: &syscall.SockaddrInet4{
			Port: int(dport),
			Addr: [4]byte{dip[0], dip[1], dip[2], dip[3]},
		},
		state: UDPSynSend,
	}
	initAuthInfo(&c.Auth)
	return c
}

// CachePacket caches the data packets while authentication is in progress
func (c *UDPConnection) CachePacket(p []byte) {

	c.packets = append(c.packets, p)
}

// TransmitCachedPackets will transmit all cached packets for this flow
func (c *UDPConnection) TransmitCachedPackets(fd int) {

	for _, p := range c.packets {
		err := syscall.Sendto(fd, p, 0, c.addr)
		if err != nil {
			//TODO: Log and continue
		}
	}
}
