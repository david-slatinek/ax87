#!/bin/bash

rm ./*.pem

## Create client private key
openssl genpkey -out client-key.pem -algorithm EC -pkeyopt ec_paramgen_curve:P-256

# View client private key
# openssl pkey -in client-key.pem -text -noout

# Create client Certificate Signing Request (CSR)
# For -config, check OpenSSL Cookbook - Unattended CSR Generation
openssl req -new -key client-key.pem -out client-req.pem -config client-config.cnf

# View CSR
# openssl req -text -in client-req.pem -noout

path="../ca-cert"

# Create a self-signed certificate
openssl x509 -req -in client-req.pem -extfile client-ext.cnf -days 60 -CA "$path"/ca-cert.pem -CAkey "$path"/ca-key.pem -CAcreateserial -out client-cert.pem

# View certificate
# openssl x509 -text -in client-cert.pem -noout

# Verify certificate
openssl verify -CAfile "$path"/ca-cert.pem client-cert.pem
