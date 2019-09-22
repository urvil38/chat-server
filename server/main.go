package main

import (
	"path/filepath"
	"strings"
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"

	"github.com/fatih/color"
)

func chatManager(clients *clients, input chan message) {
	for {
		message := <-input
		clients.mu.Lock()
		for _, client := range clients.connections {
			if client.name != message.author {
				client.conn.Write(bytes.NewBufferString(color.GreenString("                                 "+message.author+": ") + color.YellowString(message.message) + "\n").Bytes())
			}
			if client.name == message.author {
				client.conn.Write(bytes.NewBufferString(color.GreenString("Me: ") + color.YellowString(message.message) + "\n").Bytes())
			}
		}
		clients.mu.Unlock()
	}
}

type connection struct {
	name string
	conn net.Conn
}

type message struct {
	author  string
	message string
}

type clients struct {
	connections []connection
	mu          sync.Mutex
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	certDir := filepath.Join(wd,"certs")
	var port, address string
	var v bool
	flag.StringVar(&port, "p", "8080", "port on which server is running")
	flag.StringVar(&address, "addr", "127.0.0.1", "address of server")
	flag.StringVar(&certDir,"cert",certDir,"path to certs dir")
	flag.BoolVar(&v, "v", false, "version of server")
	flag.Parse()
	if v {
		fmt.Println("version 1.0.0")
		os.Exit(0)
	}

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	input := make(chan message)
	var clients clients

	serverCert, err := ioutil.ReadFile(filepath.Join(certDir,"server.pem"))
	if err != nil {
		log.Fatal("Enable to get server certificate :" + err.Error())
	}

	serverKey, err := ioutil.ReadFile(filepath.Join(certDir,"server.key"))
	if err != nil {
		log.Fatal("Enable to get server certificate key :" + err.Error())
	}

	certs, err := tls.X509KeyPair(serverCert, serverKey)
	if err != nil {
		log.Fatal(err)
	}

	tlsconfig := &tls.Config{Certificates: []tls.Certificate{certs}, Rand: rand.Reader}

	l, err := tls.Listen("tcp", address+":"+port, tlsconfig)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("started listing on address \x1b[1;33m%s\x1b[0m and port \x1b[1;33m%s\x1b[0m\n", address, port)
	for i := 0; i < 3; i++ {
		go chatManager(&clients, input)
	}

	go acceptConnection(l,&clients,input)
	<-c
	log.Println("Recevied SIGINT signal")
	log.Println("shutting down server")
	closeConnection(&clients)
	os.Exit(0)
}

func acceptConnection(l net.Listener, c *clients, input chan message) {
	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}

		go handleConn(conn, c, input)
	}
}

func closeConnection(c *clients) {
	c.mu.Lock()
	for _, c := range c.connections {
		c.conn.Close()
	}
	c.mu.Unlock()
}

func handleConn(c net.Conn, cl *clients,input chan message) {
	b := make([]byte, 256)
	c.Write([]byte("Enter Your name: "))
	numBytes, err := c.Read(b)
	if err != nil || (numBytes-2) == 0 {
		c.Write([]byte(color.RedString("You must need provide your name!\n")))
		c.Close()
	}
	name := string(bytes.Trim(b[:numBytes], "\n\r\x00"))

	con := connection{conn: c, name: name}

	cl.mu.Lock()
	cl.connections = append(cl.connections, con)
	cl.mu.Unlock()

	for {
		b := make([]byte, 2000)
		numBytes, err := c.Read(b)
	
		if err != nil || numBytes == 0 {
			c.Close()
			return
		}

		s := string(bytes.Trim(b[:numBytes], "\n\r\x00"))
		if len(strings.TrimSpace(s)) <= 0 {
			continue
		}
		
		if s == "bye" {
			c.Close()
			return
		}

		input <- message{message: s, author: con.name}
	}
}
