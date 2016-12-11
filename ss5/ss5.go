package ss5

import (
  "bytes"
  "fmt"
  "net"
  "os"
  "strconv"
  "strings"
  "github.com/BurntSushi/toml"
)

var config Config

func StartServer(address string, port int) {

  checkError := func (err error) {
    if err != nil {
      Exit(err.Error())
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

  // TODO route by config
  response := FileServer(request)

  WriteResponse(conn, response)
}

func ReadRequest(conn net.Conn) (requestBytes []byte) {
  bufferLength := 1024
  for {
    buffer := make([]byte, bufferLength)
    length, _ := conn.Read(buffer)
    requestBytes = append(requestBytes, buffer...)
    if length < bufferLength {
      return requestBytes
    }
  }
}

func ParseRequest(conn net.Conn) (request Request) {

  requestBytes := ReadRequest(conn)
  requestStr := string(requestBytes[:len(requestBytes)])
  rows := strings.Split(requestStr, "\r\n")

  // parse first line params
  firstLineParams := strings.Split(rows[0], " ")
  request.Method = firstLineParams[0]
  request.Path = firstLineParams[1]
  request.Version = firstLineParams[2]

  // parse request header
  headerStrs := rows[1:]
  request.Headers = map[string]string{}
  for i := range headerStrs {
    if headerStrs[i] == "" {
      break
    }
    keyValues := strings.Split(headerStrs[i], ":")
    request.Headers[keyValues[0]] = strings.TrimSpace( strings.Join(keyValues[1:], ":") )
  }

  // parse content type
  contentTypes := strings.Split(request.Headers["Content-Type"], ";")
  request.ContentType = contentTypes[0]

  // read request more
  isMore := false
  // TODO use const
  switch request.Method {
    case "POST", "PUT", "DELETE": isMore = true
  }
  switch request.ContentType {
    case "multipart/form-data": isMore = true
  }
  if isMore {
    request.Boundary = strings.Split(contentTypes[1], "=")[1]
    ParseRequestBody(conn, &request)
  }

  return request
}

func ParseRequestBody(conn net.Conn, request *Request) {
  request.Body = ReadRequest(conn)
  fmt.Println(string(request.Body[:len(request.Body)]))

  // parse request form
  forms := bytes.Split(request.Body, []byte(request.Boundary))
  newline := []byte("\r\n")
  colon := []byte(":")
  semiColon := []byte(";")
  equal := []byte("=")
  for i := range forms {
    form := RequestForm{}
    rows := bytes.Split(forms[i], newline)
    for j := range rows {
      if len(rows[j]) == 0 {
        if len(rows) >= j {
          form.Body = bytes.Join(rows[j:], newline)
        }
        break
      }
      keyValues := bytes.Split(rows[j], colon)
      key := keyValues[0]
      if len(keyValues) <= 1 {
        // Invalid format
        break
      }
      values := bytes.Split(bytes.Join(keyValues[1:], colon), semiColon)
      switch string(key[:len(key)]) {
        case "Content-Disposition": func() {
          form.ContentDisposition = string(values[0])
          if len(values) <= 1 {
            // break
          }
          otherParams := values[1:]
          for k := range otherParams {
            innerKeyValue := bytes.Split(otherParams[k], equal)
            if len(innerKeyValue) <= 1 {
              break
            }
            innerValue := string(bytes.Join(innerKeyValue[1:], equal))
            switch string(innerKeyValue[0]) {
              case "name": form.Name = innerValue
              case "filename": form.FileName = innerValue
            }
          }
        }()
        case "Content-Type": // form.ContentType =
      }
    }
  }
}

func WriteResponse(conn net.Conn, response Response) {
  defer conn.Close()
  header := "HTTP/1.1 " + response.Status
  header += "\r\n"
  header += "Content-Type: " + response.ContentType
  header += "\r\n"
  header += "Content-Length: " + strconv.Itoa(len(response.Body))
  header += "\r\n"
  for k, v := range response.Headers {
    header += k + ": " + v
    header += "\r\n"
  }
  header += "\r\n"
  conn.Write([]byte(header))
  conn.Write(response.Body)
  conn.Write([]byte("\r\n"))
}

func Exit(message string) {
  fmt.Fprintf(os.Stderr, "[ERROR]: %s\n", message)
  os.Exit(1)
}
