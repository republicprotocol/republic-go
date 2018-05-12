package client

import (
	"context"
	"fmt"
	"sync"

	"github.com/republicprotocol/republic-go/identity"
	"google.golang.org/grpc"
)

// Conn is an atomically reference counted grpc.ClientConn. It must only be
// created by a call to Dial, which must be accompanied by exactly one call to
// conn.Close. It is safe for concurrent use.
type Conn struct {
	*grpc.ClientConn

	rcMu *sync.Mutex
	rc   int64
}

// Dial creates a client connection to the given multiaddress. A context can be
// used to cancel or expire the pending connection. Once this function returns,
// the cancellation and expiration of the Context will do nothing. Users must
// call Conn.Close to terminate all the pending operations after this function
// returns. Users must call Conn.Clone to create copies of the connection.
func Dial(ctx context.Context, multiAddress identity.MultiAddress) (*Conn, error) {
	host, err := multiAddress.ValueForProtocol(identity.IP4Code)
	if err != nil {
		return nil, err
	}
	port, err := multiAddress.ValueForProtocol(identity.TCPCode)
	if err != nil {
		return nil, err
	}
	clientConn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%s", host, port), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &Conn{
		ClientConn: clientConn,
		rcMu:       new(sync.Mutex),
		rc:         1,
	}, nil
}

// Clone the connection. Every call to Conn.Clone must be accompanied by
// exactly one call to Conn.Close after the connection is no longer needed.
func (conn *Conn) Clone() *Conn {
	conn.rcMu.Lock()
	defer conn.rcMu.Unlock()

	conn.rc++
	return conn
}

// Close the connection. If there are no other references to the connection, as
// created by Conn.Clone, then the connection will be closed and all pending
// operations will be terminated.
func (conn *Conn) Close() error {
	conn.rcMu.Lock()
	defer conn.rcMu.Unlock()

	conn.rc--
	if conn.rc == 0 {
		return conn.ClientConn.Close()
	}
	return nil
}
