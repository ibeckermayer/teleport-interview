package main

import (
	"flag"
	"log"

	"github.com/ibeckermayer/teleport-interview/backend/internal/server"
)

func main() {
	port := flag.Int("port", 8000, "Port to serve the app from")
	certFilePath := flag.String("cert", "../certs/localhost.crt", "Relative path to a valid SSL cert")
	keyFilePath := flag.String("key", "../certs/localhost.key", "Relative path to the cert's private key")
	flag.Parse()

	cfg := server.NewConfig(*port, *certFilePath, *keyFilePath)
	srv := server.New(cfg)

	log.Fatal(srv.Run())
}
