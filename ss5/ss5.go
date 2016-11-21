package ss5

import (
  "fmt"
  "net"
  "os"
  "path/filepath"
  "strconv"
  "strings"
)

type Request struct {
  method string
  path string
  version string
}

type Response struct {
  status string
  contentType string
  body []byte
}

const (
  OK = "200 OK"
  NOT_FOUND = "404 Not Found"
)

const (
  PLAIN = "text/plain"
  HTML = "text/html"
  PNG = "image/png"
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
  request := parseRequest(conn)

  fmt.Printf("connection start. requestPath=%s, remoteAddr=%s\n", request.path, conn.RemoteAddr())

  // TODO read rule
  targetFile := "public"
  if (request.path == "/") {
    targetFile += "/index.html"
  } else {
    targetFile += request.path
  }
  file, err := os.Open(targetFile)

  response := Response{}
  if err == nil {
    response.status = OK
    response.contentType = findContentType(targetFile)
    for {
      bytes := make([]byte, 1024)
      _, err := file.Read(bytes)
      response.body = append(response.body, bytes...)
      if err != nil {
        break
      }
    }
  } else {
    fmt.Printf("ERROR! %s\n", err)
    response.status = NOT_FOUND
    response.contentType = HTML
    response.body = []byte(`<html>
<head>
  <meta charset="utf-8">
  <title>404 Not Found</title>
</head>
<body>
  <p>404 Not Found! 😇</p>
</body>
</html>`)
  }

  write(conn, response)
  fmt.Printf("connection close.\n")
}

func parseRequest(conn net.Conn) Request {
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

func findContentType(filePath string) string {
  ext := filepath.Ext(filePath)
  switch ext {
    case ".html": return HTML
    case ".png": return PNG
    default: return PLAIN
  }
}

func write(conn net.Conn, response Response) {
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

func checkError(err error) {
  if err != nil {
    fmt.Fprintf(os.Stderr, "[ERROR]: %s\n", err.Error())
		os.Exit(1)
  }
}
