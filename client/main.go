package main

import (
	"io/ioutil"
	"crypto/tls"
	"flag"
	"log"

	"github.com/reiver/go-telnet"
)


func main() {
	var address string
	var port string
	flag.StringVar(&address, "addr", "", "address of server you want to connect")
	flag.StringVar(&port, "p", "", "port number of server")
	flag.Parse()

	clientCert,err := ioutil.ReadFile("../certs/client.pem")
	if err != nil {
		panic(err)
	}

	clientKey,err := ioutil.ReadFile("../certs/client.key")
	if err != nil {
		panic(err)
	}

	certs, err := tls.X509KeyPair(clientCert, clientKey)
	if err != nil {
		panic(err)
	}

	tlsconfig := &tls.Config{
		Certificates:       []tls.Certificate{certs},
		InsecureSkipVerify: true,
	}
	caller := telnet.StandardCaller
	log.Fatal(telnet.DialToAndCallTLS(address+":"+port, caller, tlsconfig))
}
