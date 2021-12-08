package auth

import (
	"context"
	"errors"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
)

var adminEmails = map[string]bool{
	"admin@localhost": true,
}

func isAdmin(email string) bool {
	return adminEmails[email]
}

// GetPeerEmail returns the first email address found in SAN of the client
// certificate of the peer attached to the given request context.
func GetPeerEmail(ctx context.Context) (string, error) {
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

// IsAdmin returns whether the peer of the context is an admin.
func IsAdmin(ctx context.Context) (bool, error) {
	email, err := GetPeerEmail(ctx)
	if err != nil {
		return false, err
	}
	return isAdmin(email), nil
}

// CanManage returns whether the peer of the context can manage resources owned
// by owner: the peer needs to be either the owner or an admin.
func CanManage(ctx context.Context, owner string) (bool, error) {
	email, err := GetPeerEmail(ctx)
	if err != nil {
		return false, err
	}
	return isAdmin(email) || email == owner, nil
}
