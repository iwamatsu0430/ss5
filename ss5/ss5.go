package ss5

import (
  "fmt"
  "net"
  "os"
  "time"
  "strconv"
  "strings"
  "github.com/BurntSushi/toml"
)

var config Config

func StartServer(address string, port int) {
  checkError := func (err error) {
    if err != nil {
      fmt.Fprintf(os.Stderr, "[ERROR]: %s\n", err.Error())
  		os.Exit(1)
    }
  }

  _, err := toml.DecodeFile("config.toml", &config)
  checkError(err)

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

  fmt.Printf("======================= connected. requestPath=%s, remoteAddr=%s =======================\n", request.path, conn.RemoteAddr())

  // TODO route by config
  response := FileServer(request)

  WriteResponse(conn, response)
}

func ParseRequest(conn net.Conn) (request Request) {

  var requestBytes []byte
  bufferLength := 1024
  for {
    buffer := make([]byte, bufferLength)
    conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
    length, _ := conn.Read(buffer)
    requestBytes = append(requestBytes, buffer...)
    fmt.Println("Length: ", length)
    if length < bufferLength {
      break
    }
  }

  requestStr := string(requestBytes[:len(requestBytes)])
  // headerBodies := strings.Split(requestStr, "\r\n\r\n")

  // firstLineParams := strings.Split(rows[0], " ")
  // request.method = firstLineParams[0]
  // request.path = firstLineParams[1]
  // request.version = firstLineParams[2]
  strings.Split(",", " ")
  fmt.Printf("\n%s\n\n", requestStr)
  request.path = "/"

  // requestStr := ""
  // for {
  //   buffer := make([]byte, 1024)
  //   conn.SetReadDeadline(time.Now().Add(10 * time.Millisecond))
  //   length, _ := conn.Read(buffer)
  //   requestStr += string(buffer[:length])
  //   if length == 0 {
  //     break
  //   }
  //   // if strings.HasSuffix(requestStr, "\r\n\r\n") {
  //   //   break
  //   // }
  // }
  // rows := strings.Split(requestStr, "\r\n")
  //
  // firstLineParams := strings.Split(rows[0], " ")
  // request.method = firstLineParams[0]
  // request.path = firstLineParams[1]
  // request.version = firstLineParams[2]
  //
  // headers := map[string]string{}
  // for i, _ := range rows[1:] {
  //   if rows[i] == "" {
  //     break
  //   }
  //   keyValues := strings.Split(rows[i], ":")
  //   headers[keyValues[0]] = strings.Join(keyValues[1:], ":")
  //   fmt.Printf("header = %s, values = %s\n", keyValues[0], strings.Join(keyValues[1:], ":"))
  // }

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
