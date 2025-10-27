package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	port := os.Getenv("MARKET_TEST_PORT")
	if port == "" {
		port = "40100"
	}
	addr := "0.0.0.0" + ":" + port

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	fmt.Println("Listening on", addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}

		fmt.Println("Client connected:", conn.RemoteAddr())
		go func(c net.Conn) {
			defer c.Close()
			for {
				c.Write([]byte("marketflow tick\n"))
				// эмулируем потоковые данные
				// time.Sleep(time.Second)
			}
		}(conn)
	}
}
