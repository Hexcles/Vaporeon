package auth

import (
	"context"
	"errors"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
)

var adminEmails = []string{"admin@localhost"}

// EmailAuth is a simple gRPC authenticator based on the email address in the
// client certificate with a list of hard-coded admins.
//
// The type is read-only after created and can be copied.
type EmailAuth struct {
	adminEmails map[string]bool
}

// NewEmailAuth creates a new EmailAuth.
func NewEmailAuth() EmailAuth {
	e := EmailAuth{adminEmails: make(map[string]bool)}
	for _, email := range adminEmails {
		e.adminEmails[email] = true
	}
	return e
}

func (e EmailAuth) isAdmin(email string) bool {
	return e.adminEmails[email]
}

// GetPeerID returns the first email address found in SAN of the client
// certificate of the peer attached to the given request context.
func (e EmailAuth) GetPeerID(ctx context.Context) (string, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return "", errors.New("auth: peer not found in context")
	}
	tlsInfo, ok := p.AuthInfo.(credentials.TLSInfo)
	if !ok {
		return "", errors.New("auth: peer does not contain TLS info")
	}
	if len(tlsInfo.State.PeerCertificates) == 0 {
		return "", errors.New("auth: no certificate in TLS info")
	}
	cert := tlsInfo.State.PeerCertificates[0]
	if len(cert.EmailAddresses) == 0 {
		return "", errors.New("auth: no email address found in leaf cert")
	}
	return cert.EmailAddresses[0], nil
}

// CanShutdown returns whether the peer can shut down the server: the peer
// needs to be an admin.
func (e EmailAuth) CanShutdown(ctx context.Context) (bool, error) {
	email, err := e.GetPeerID(ctx)
	if err != nil {
		return false, err
	}
	return e.isAdmin(email), nil
}

// CanManage returns whether the peer of the context can manage resources
// owned by owner: the peer needs to be either the owner or an admin.
func (e EmailAuth) CanManage(ctx context.Context, owner string) (bool, error) {
	email, err := e.GetPeerID(ctx)
	if err != nil {
		return false, err
	}
	return e.isAdmin(email) || email == owner, nil
}
