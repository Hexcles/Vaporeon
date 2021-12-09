package auth

import (
	"crypto/tls"
	"crypto/x509"
	"testing"

	"golang.org/x/net/context"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
)

var a = NewEmailAuth()

func contextWithPeerEmail(ctx context.Context, email string) context.Context {
	return peer.NewContext(ctx, &peer.Peer{
		AuthInfo: credentials.TLSInfo{
			State: tls.ConnectionState{
				PeerCertificates: []*x509.Certificate{{
					EmailAddresses: []string{email},
				}},
			},
		},
	})
}

func TestGetPeerID_error(t *testing.T) {
	bgCtx := context.Background()
	contexts := []context.Context{
		bgCtx,
		// No email in SAN.
		peer.NewContext(bgCtx, &peer.Peer{
			AuthInfo: credentials.TLSInfo{
				State: tls.ConnectionState{
					PeerCertificates: []*x509.Certificate{{}},
				},
			},
		}),
		// No cert.
		peer.NewContext(bgCtx, &peer.Peer{
			AuthInfo: credentials.TLSInfo{},
		}),
		// No TLSInfo.
		peer.NewContext(bgCtx, &peer.Peer{}),
	}
	for _, ctx := range contexts {
		if got, err := a.GetPeerID(ctx); got != "" || err == nil {
			t.Errorf(`GetPeerID(empty context) = %q, %v; want "", non-nil`, got, err)
		}

	}
}

func TestGetPeerID_success(t *testing.T) {
	want := "test@example.com"
	ctx := contextWithPeerEmail(context.Background(), want)
	if got, err := a.GetPeerID(ctx); got != want || err != nil {
		t.Errorf(`GetPeerID(fake context) = %q, %v; want %q, nil`, got, err, want)
	}
}

func TestCanShutdown(t *testing.T) {
	ctx := context.Background()
	if ok, err := a.CanShutdown(ctx); ok || err == nil {
		t.Errorf("CanShutdown(empty context) = %v, %v; want false, non-nil", ok, err)
	}
	if ok, err := a.CanShutdown(contextWithPeerEmail(ctx, "guest@localhost")); ok || err != nil {
		t.Errorf("CanShutdown(guest) = %v, %v; want false, nil", ok, err)
	}
	if ok, err := a.CanShutdown(contextWithPeerEmail(ctx, "admin@localhost")); !ok || err != nil {
		t.Errorf("CanShutdown(admin) = %v, %v; want true, nil", ok, err)
	}
}

func TestCanManage(t *testing.T) {
	ctx := context.Background()
	owner := "guest@localhost"
	if ok, err := a.CanManage(ctx, owner); ok || err == nil {
		t.Errorf("CanManage(empty context, guest) = %v, %v; want false, non-nil", ok, err)
	}
	if ok, err := a.CanManage(contextWithPeerEmail(ctx, "guest2@localhost"), owner); ok || err != nil {
		t.Errorf("CanManage(guest2, guest) = %v, %v; want false, nil", ok, err)
	}
	if ok, err := a.CanManage(contextWithPeerEmail(ctx, "guest@localhost"), owner); !ok || err != nil {
		t.Errorf("CanManage(guest, guest) = %v, %v; want true, nil", ok, err)
	}
	if ok, err := a.CanManage(contextWithPeerEmail(ctx, "admin@localhost"), owner); !ok || err != nil {
		t.Errorf("CanManage(admin, guest) = %v, %v; want true, nil", ok, err)
	}
}
