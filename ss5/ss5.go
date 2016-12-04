package ss5

import (
  "fmt"
  "net"
  "os"
  "strconv"
  "strings"
)

func StartServer(address string, port int) {
  checkError := func (err error) {
    if err != nil {
      fmt.Fprintf(os.Stderr, "[ERROR]: %s\n", err.Error())
  		os.Exit(1)
    }
  }
  portStr := strconv.Itoa(port)
  tcpAddr, err := net.ResolveTCPAddr("tcp", ":" + portStr)
  checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
  checkError(err)
  fmt.Printf("Server Start. Access http://%s:%s/\n", address, portStr)
  Listen(listener)
}

func Listen(listener *net.TCPListener) {
  for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

    go CreateResponse(conn)
	}
}

func CreateResponse(conn net.Conn) {
  request := ParseRequest(conn)

  fmt.Printf("connection start. requestPath=%s, remoteAddr=%s\n", request.path, conn.RemoteAddr())

  // TODO route by config
  response := FileServer(request)

  WriteResponse(conn, response)
  fmt.Printf("connection close.\n")
}

func ParseRequest(conn net.Conn) Request {
  requestStr := ""
  for {
    buffer := make([]byte, 1024)
    length, _ := conn.Read(buffer)
    requestStr += string(buffer[:length])
    if strings.HasSuffix(requestStr, "\r\n\r\n") {
      break
    }
  }
  rows := strings.Split(requestStr, "\r\n")
  params := strings.Split(rows[0], " ")
  request := Request{}
  request.method = params[0]
  request.path = params[1]
  request.version = params[2]
  // TODO read headers
  return request
}

func WriteResponse(conn net.Conn, response Response) {
  defer conn.Close()
  header := "HTTP/1.1 " + response.status
  header += "\r\n"
  header += "Content-Type: " + response.contentType
  header += "\r\n"
  header += "Content-Length: " + strconv.Itoa(len(response.body))
  header += "\r\n"
  // TODO add extra headers
  header += "\r\n"
  conn.Write([]byte(header))
  conn.Write(response.body)
  conn.Write([]byte("\r\n"))
}
