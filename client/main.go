package main

import (
	"os"
	"path/filepath"
	"io/ioutil"
	"crypto/tls"
	"flag"
	"log"

	"github.com/reiver/go-telnet"
)


func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	certDir := filepath.Join(wd,"certs")
	var address string
	var port string
	flag.StringVar(&address, "addr", "", "address of server you want to connect")
	flag.StringVar(&certDir,"cert",certDir,"path to certs dir")
	flag.StringVar(&port, "p", "", "port number of server")
	flag.Parse()

	clientCert,err := ioutil.ReadFile(filepath.Join(certDir,"client.pem"))
	if err != nil {
		log.Fatal(err)
	}

	clientKey,err := ioutil.ReadFile(filepath.Join(certDir,"client.key"))
	if err != nil {
		log.Fatal(err)
	}

	certs, err := tls.X509KeyPair(clientCert, clientKey)
	if err != nil {
		log.Fatal(err)
	}

	tlsconfig := &tls.Config{
		Certificates:       []tls.Certificate{certs},
		InsecureSkipVerify: true,
	}

	caller := telnet.StandardCaller
	err = telnet.DialToAndCallTLS(address+":"+port, caller,tlsconfig)
	if err != nil {
		log.Fatal(err)
	}
}
