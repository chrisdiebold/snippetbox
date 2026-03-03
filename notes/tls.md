# TLS and HTTPS

HTTPS is essentially HTTP transmitted over a TLS (Transport Layer Security) connection. The advantage to this is that HTTPS traffic is encrypted and signed, which helps ensure its privacy and integrity during transit.

For production servers I recommend using Let’s Encrypt to create your TLS certificates, but for development purposes the simplest thing to do is to generate your own self-signed certificate.

A self-signed certificate is the same as a normal TLS certificate, except that it isn’t cryptographically signed by a trusted certificate authority. This means that your web browser will display a warning the first time it’s used, but it will nonetheless encrypt HTTPS traffic correctly and is fine for development and testing purposes.

Handily, the crypto/tls package in Go’s standard library includes a generate_cert.go tool that we can use to easily create our own self-signed certificate.

To run the generate_cert.go tool, you’ll need to know the location on your computer where the source code for the Go standard library is installed. If you’re using Linux, macOS or FreeBSD and followed the official install instructions, then the generate_cert.go file should be located under `/usr/local/go/src/crypto/tls`.


```bash
go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
```

Behind the scenes, this generate_cert.go command works in two stages:
1. First it generates a 2048-bit RSA key pair, which is a cryptographically secure public key and private key.
2. It then stores the private key in a key.pem file, and generates a self-signed TLS certificate for the host localhost containing the public key — which it stores in a cert.pem file. Both the private key and certificate are encoded in PEM format, which is the standard format used by most TLS implementations.


As an alternative to the generate_cert.go tool, you might want to consider using [mkcert](https://github.com/FiloSottile/mkcert) to generate the TLS certificates. Although this requires some extra setup, it has the advantage that the generated certificates are locally trusted — meaning that you can use them for testing and development without getting security warnings in your web browser.

A big plus of using HTTPS is that Go will automatically upgrade the connection to use HTTP/2 if the client supports it.