package main

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"io/ioutil"

	"github.com/urvil38/color"
)

func Color(s string, a color.Attribute) string {
	return color.Wrap(s, a, color.Bold)
}

func chatManager(clients *clients, input chan message) {
	for {
		message := <-input
		clients.mu.Lock()
		for _, client := range clients.connections {
			if client.name != message.author {
				client.conn.Write(bytes.NewBufferString(Color("                                 "+message.author+": ", color.FgGreen) + Color(message.message, color.FgYellow) + "\n").Bytes())
			}
			if client.name == message.author {
				client.conn.Write(bytes.NewBufferString(Color("Me"+": ", color.FgGreen) + Color(message.message, color.FgYellow) + "\n").Bytes())
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
	var port, address string
	var v bool
	flag.StringVar(&port, "p", "8080", "port on which server is running")
	flag.StringVar(&address, "addr", "127.0.0.1", "address of server")
	flag.BoolVar(&v, "v", false, "version of server")
	flag.Parse()
	if v {
		fmt.Println("version 1.0.0")
		os.Exit(0)
	}
	input := make(chan message)
	var clients clients

	serverCert,err := ioutil.ReadFile("../certs/server.pem")
	if err != nil {
		panic("Enable to get server certificate :"+err.Error())
	}

	serverKey,err := ioutil.ReadFile("../certs/server.key")
	if err != nil {
		panic("Enable to get server certificate key :"+err.Error())
	}

	certs, err := tls.X509KeyPair(serverCert, serverKey)
	if err != nil {
		panic("failed to parse root certificate")
	}

	tlsconfig := &tls.Config{Certificates: []tls.Certificate{certs}, Rand: rand.Reader}

	l, err := tls.Listen("tcp", address+":"+port, tlsconfig)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("started listing on address \x1b[1;33m%s\x1b[0m and port \x1b[1;33m%s\x1b[0m\n", address, port)
	for i := 0 ; i < 3 ; i++ {
		go chatManager(&clients, input)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}

		go clients.handleConn(conn, input)
	}
}

func (cl *clients) handleConn(c net.Conn, input chan message) {
	b := make([]byte, 256)
	c.Write([]byte("Enter Your name: "))
	numBytes, err := c.Read(b)
	if err != nil || (numBytes-2) == 0 {
		c.Write([]byte(Color("You must need provide your name!\n", color.FgRed)))
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
			continue
		}
		s := string(bytes.Trim(b[:numBytes], "\n\r\x00"))
		if s == "bye" {
			c.Close()
			break
		}
		input <- message{message: s, author: con.name}
	}
}
