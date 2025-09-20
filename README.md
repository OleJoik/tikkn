# tikkn

### Initialize

Create dev certificates for TLS termination

```
mkdir -p local/certs
openssl req -x509 -newkey rsa:4096 -nodes \
    -keyout local/certs/privkey.pem \
    -out local/certs/cert.pem -days 365 \
    -subj "/CN=localhost"
```
