#!/bin/bash

rm server-cert.pem server-key.pem server-req.pem

# Create server private key
openssl genpkey -out server-key.pem -algorithm EC -pkeyopt ec_paramgen_curve:P-256

# View server private key
# openssl pkey -in server-key.pem -text -noout

# Create server Certificate Signing Request (CSR)
# For -config, check OpenSSL Cookbook - Unattended CSR Generation
openssl req -new -key server-key.pem -out server-req.pem -config server-config.cnf

# View CSR
# openssl req -text -in server-req.pem -noout

# Create a self-signed certificate
openssl x509 -req -in server-req.pem -extfile server-ext.cnf -days 60 -CA ../../ca-cert/ca-cert.pem -CAkey ../../ca-cert/ca-key.pem -CAcreateserial -out server-cert.pem

# View certificate
# openssl x509 -text -in server-cert.pem -noout

# Verify certificate
openssl verify -CAfile ../../ca-cert/ca-cert.pem server-cert.pem
