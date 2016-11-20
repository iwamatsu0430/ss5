package ss5

import (
  "fmt"
  "net"
  "os"
  "strconv"
  "time"
)

func Start(address string, port int) {
  Init(address, port)
}

func Init(address string, port int) {
  portStr := strconv.Itoa(port)
  tcpAddr, err := net.ResolveTCPAddr("tcp", ":" + portStr)
  checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
  checkError(err)
  fmt.Printf("Server Start. Access http://%s:%s/\n", address, portStr)
  listen(listener)
}

func listen(listener *net.TCPListener) {
  for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

    go response(conn)
	}
}

func response(conn net.Conn) {
  fmt.Printf("connection start. remoteAddr=%s\n", conn.RemoteAddr())

  body := `<html>
<body>
<h1>Hello, World!</h1>
</body>
</html>`
  time.Sleep(5 * time.Second)
  output(conn, "200 OK", body)

  fmt.Printf("connection close.\n")
}

func output(conn net.Conn, status string, body string) {
  defer conn.Close()
  message := "HTTP/1.1 " + status
  message += "\r\n"
  message += "Content-Type: text/html"
  message += "\r\n"
  message += "Content-Length: " + strconv.Itoa(len(body))
  message += "\r\n"
  message += "\r\n"
  message += body
  message += "\r\n"
  conn.Write([]byte(message))
}

func checkError(err error) {
  if err != nil {
    fmt.Fprintf(os.Stderr, "[ERROR]: %s\n", err.Error())
		os.Exit(1)
  }
}
