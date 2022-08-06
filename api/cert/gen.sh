#!/bin/bash

# Create Certificate Authority (CA) private key
openssl genpkey -out ca-key.pem -algorithm EC -pkeyopt ec_paramgen_curve:P-256 -aes-256-cbc

# View CA private key
# openssl pkey -in ca-key.pem -text -noout

# Create CA certificate
openssl req -new -x509 -days 365 -key ca-key.pem -out ca-cert.pem

# Create server key
openssl genpkey -out server-key.pem -algorithm EC -pkeyopt ec_paramgen_curve:P-256

# Create server Certificate Signing Request (CSR)
openssl req -new -key server-key.pem -out server-req.pem

# View CSR
# openssl req -text -in server-req.pem -noout

# Create a self-signed certificate
openssl x509 -req -in server-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem

# View certificate
# openssl x509 -text -in server-cert.pem -noout

# Verify certificate
# openssl verify -CAfile ca-cert.pem server-cert.pem
