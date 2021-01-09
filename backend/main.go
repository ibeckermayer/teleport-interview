package main

import (
	"flag"
	"log"
	"time"

	"github.com/ibeckermayer/teleport-interview/backend/internal/server"
)

func main() {
	port := flag.Int("port", 8000, "Port to serve the app from")
	certFilePath := flag.String("cert", "../certs/localhost.crt", "Relative path to a valid SSL cert")
	keyFilePath := flag.String("key", "../certs/localhost.key", "Relative path to the cert's private key")
	sessionTimeout := flag.String("sesh", "12h", "A parseable duration string (https://golang.org/pkg/time/#ParseDuration) specifying the absolute timeout value for user sessions")
	flag.Parse()

	timeout, err := time.ParseDuration(*sessionTimeout)
	if err != nil {
		log.Fatalf("failed to parse duration string for command line flag sesh=%v; see https://golang.org/pkg/time/#ParseDuration", *sessionTimeout)
	}

	cfg := server.NewConfig(*port, *certFilePath, *keyFilePath, timeout)
	srv := server.New(cfg)

	log.Fatal(srv.Run())
}
