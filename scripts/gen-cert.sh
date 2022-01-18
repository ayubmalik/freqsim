#!/usr/bin/env bash

scriptdir="$(cd $(dirname "${0}") && pwd)"
cd ${scriptdir}

rm *.pem

subject="/C=UK/ST=Greater Manchester/L=Manchester/O=3H Services Ltd/OU=Development/CN=api.cloud-dev.net/emailAddress=ayub.malik@gmail.com"

# 1. Generate CA's private key and self-signed certificate
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout ca.key -out ca.crt -subj "${subject}"

echo "CA's self-signed certificate"
openssl x509 -in ca.crt -noout -text

# # # 2. Generate web server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout server.key -out server.csr -subj "${subject}"

# # # 3. Use CA's private key to sign web server's CSR and get back the signed certificate
openssl x509 -req -days 365 -CA ca.crt -CAkey ca.key -CAcreateserial -in server.csr -out server.crt -extfile <(printf "subjectAltName=DNS:localhost")

echo "Server's signed certificate"
openssl x509 -in server.crt -noout -text
