#!/bin/bash
set -ex

# Create the CA certs.
openssl req -x509                                     \
  -newkey rsa:4096                                    \
  -nodes                                              \
  -days 3650                                          \
  -keyout ca_key.pem                                  \
  -out ca_cert.pem                                    \
  -subj /CN=test-ca/                                  \
  -config ./openssl.cnf                               \
  -extensions test_ca

# Generate the server certs.
openssl genrsa -out server_key.pem 4096
openssl req -new                                    \
  -key server_key.pem                               \
  -out server_csr.pem                               \
  -subj /CN=test-server/                            \
  -config ./openssl.cnf                             \
  -reqexts test_server
openssl x509 -req           \
  -in server_csr.pem        \
  -CAkey ca_key.pem         \
  -CA ca_cert.pem           \
  -days 3650                \
  -set_serial 1000          \
  -out server_cert.pem      \
  -extfile ./openssl.cnf    \
  -extensions test_server
openssl verify -verbose -CAfile ca_cert.pem  server_cert.pem

# Generate two client certs.
openssl genrsa -out client1_key.pem 4096
openssl req -new                                    \
  -key client1_key.pem                              \
  -out client1_csr.pem                              \
  -subj /CN=test-client1/                           \
  -config ./openssl.cnf                             \
  -reqexts test_client1
openssl x509 -req           \
  -in client1_csr.pem       \
  -CAkey ca_key.pem         \
  -CA ca_cert.pem           \
  -days 3650                \
  -set_serial 1001          \
  -out client1_cert.pem     \
  -extfile ./openssl.cnf    \
  -extensions test_client1
openssl verify -verbose -CAfile ca_cert.pem  client1_cert.pem

openssl genrsa -out client2_key.pem 4096
openssl req -new                                    \
  -key client2_key.pem                              \
  -out client2_csr.pem                              \
  -subj /CN=test-client2/                           \
  -config ./openssl.cnf                             \
  -reqexts test_client2
openssl x509 -req           \
  -in client2_csr.pem       \
  -CAkey ca_key.pem         \
  -CA ca_cert.pem           \
  -days 3650                \
  -set_serial 1002          \
  -out client2_cert.pem     \
  -extfile ./openssl.cnf    \
  -extensions test_client2
openssl verify -verbose -CAfile ca_cert.pem  client2_cert.pem

rm *_csr.pem
