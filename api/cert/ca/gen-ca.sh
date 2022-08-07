#!/bin/bash

rm ca-cert.pem ca-cert.srl ca-key.pem

# Create Certificate Authority (CA) private key
openssl genpkey -out ca-key.pem -algorithm EC -pkeyopt ec_paramgen_curve:P-256 -aes-256-cbc

# View CA private key
# openssl pkey -in ca-key.pem -text -noout

# Create CA certificate
openssl req -new -x509 -days 365 -config ca-config.cnf -key ca-key.pem -out ca-cert.pem

# View certificate
# openssl x509 -in ca-cert.pem -noout -text
