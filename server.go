package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	tlsConfig, err := tlsConfig()
	if err != nil {
		log.Fatal(err)
	}

	// create a server listening on 0.0.0.0:9000 using the privided config.
	srv := &http.Server{
		Addr:      ":9000",
		Handler:   index(),
		TLSConfig: tlsConfig,
	}

	// run the server with the appropriate keys.
	log.Fatal(srv.ListenAndServeTLS("server.pem", "server.key"))
}

// HTTP handler used by the server.
// If you go to `https://localhost:9000` you'll be prompted for the TLS certs.
// Only if the CA validates the config will you be able to see the message in the <h1>.
func index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := `
		<h1>Congrats you passed TLS auth.</h1>
		`
		w.Write([]byte(msg))
	}
}

// construct a tls config that requires client certs and uses the SCEP CA to validate them.
func tlsConfig() (*tls.Config, error) {
	caPEMBytes, err := ioutil.ReadFile("ca/depot/ca.pem")
	if err != nil {
		return nil, err
	}
	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caPEMBytes)

	tlsConfig := &tls.Config{
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  caPool,
	}
	return tlsConfig, nil
}
