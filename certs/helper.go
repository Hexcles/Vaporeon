package certs

import (
	"crypto/x509"
	"errors"
	"os"
)

var ErrInvalidCA = errors.New("unable to add CA to pool; cert may be invalid")

func LoadCA(caFile string) (*x509.CertPool, error) {
	ca, err := os.ReadFile(caFile)
	if err != nil {
		return nil, err
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(ca) {
		return nil, ErrInvalidCA
	}
	return pool, nil
}
