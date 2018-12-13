package main

import (
	"net"
	"github.com/astaxie/beego/logs"
)

const HELLO_WORLD = `HTTP/1.1 200 OK
Server: openresty/1.11.2.5
Date: Thu, 22 Nov 2018 08:03:03 GMT
Content-Type: text/html; charset=utf-8
Content-Length: 12
Last-Modified: Sun, 30 Sep 2018 03:38:39 GMT
Connection: keep-alive
ETag: "5bb0453f-b54"
Access-Control-Allow-Origin: *
Access-Control-Allow-Headers: reqid, nid, host, x-real-ip, x-forwarded-ip, event-type, event-id, accept, content-type,x-token
Cache-Control: public
Accept-Ranges: bytes


hello world


`

func main() {

	listen, err := net.Listen("tcp", ":1024")
	if err != nil {
		panic(err)
	}
	defer listen.Close()

	logs.Debug("waiting for clients")
	for {
		con, err := listen.Accept()
		if err != nil {
			continue
		}

		go handleConnections(con)
	}
}

func handleConnections(con net.Conn) {
	buffer := make([]byte, 2048)
	con.Write([]byte(HELLO_WORLD))
	defer con.Close()
	for {
		n, err := con.Read(buffer)
		if err != nil {
			return
		}
		logs.Info(con.RemoteAddr(), ", receive data string:", string(buffer[:n]))
	}
}
