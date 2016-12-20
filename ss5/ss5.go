package ss5

import (
  "bytes"
  "fmt"
  "net"
  "os"
  "strconv"
  "strings"
  "time"
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

  response := FileServer(request)

  WriteResponse(conn, response)
}

func ParseRequest(conn net.Conn) (request Request) {

  requestBytes := ReadRequest(conn, 0)

  ParseRequestHeader(requestBytes, &request)

  ParseContentTypes(&request)

  ParseContentLength(&request)

  ParseRequestBody(conn, &request)

  return request
}

func ReadRequest(conn net.Conn, timeout int) (requestBytes []byte) {
  bufferLength := 1024
  for {
    buffer := make([]byte, bufferLength)
    if timeout > 0 {
      conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Millisecond))
    }
    length, _ := conn.Read(buffer)
    requestBytes = append(requestBytes, buffer...)
    if length < bufferLength {
      return requestBytes
    }
  }
}

func ParseRequestHeader(requestRaw []byte, request *Request) {
  byteNewLine := []byte("\r\n")
  requestRows := bytes.Split(requestRaw, byteNewLine)
  request.Headers = map[string]string{}
  for i, requestRowBytes := range requestRows {
    if len(requestRowBytes) == 0 {
      request.Body = bytes.Join(requestRows[i:], byteNewLine)
      break
    }
    requestRow := string(requestRowBytes)
    if i == 0 {
      params := strings.Split(requestRow, " ")
      request.Method = params[0]
      request.Path = params[1]
      request.Version = params[2]
    } else {
      keyValues := strings.Split(requestRow, ":")
      request.Headers[keyValues[0]] = strings.TrimSpace( strings.Join(keyValues[1:], ":") )
    }
  }
}

func ParseContentTypes(request *Request) {
  contentType, existsContentType := request.Headers["Content-Type"]
  if existsContentType {
    params := strings.Split(contentType, ";")
    request.ContentType = params[0]
    if len(params) <= 1 {
      return
    }
    for _, param := range params[1:] {
      keyValues := strings.Split(param, "=")
      if len(keyValues) >= 2 {
        switch strings.Trim(keyValues[0], " ") {
          case "boundary": request.Boundary = keyValues[1]
        }
      }
    }
  }
}

func ParseContentLength(request *Request) {
  contentLength, existsContentLength := request.Headers["Content-Length"]
  if existsContentLength {
    contentLengthValue, err := strconv.Atoi(contentLength)
    if err == nil {
      request.ContentLength = contentLengthValue
    }
  }
}

func ParseRequestBody(conn net.Conn, request *Request) {
  switch request.Method {
    case HTTP_METHOD_POST, HTTP_METHOD_PUT, HTTP_METHOD_DELETE:
    default: return
  }
  if request.ContentLength <= 0 {
    return
  }
  if len(request.Body) == 0 {
    request.Body = ReadRequest(conn, 100)
  }
  switch request.ContentType {
    case CONTENT_TYPE_MP_FD, CONTENT_TYPE_MP_MX, CONTENT_TYPE_MP_RL:
    default: return
  }

  // Separate by Boundary
  forms := bytes.Split(request.Body, []byte("--" + request.Boundary))
  trimmedForms := initByte(tailByte(forms))
  for _, form := range trimmedForms {
    requestForm := RequestForm{}

    // Separate by line
    rows := bytes.Split(form, []byte("\r\n"))
    for j, row := range rows {
      if len(row) == 0 {
        if j > 0 {
          requestForm.Body = bytes.Trim(bytes.Join(rows[j+1:], []byte("\r\n")), "\n")
          break
        } else {
          continue
        }
      }

      // Separate by colon
      headerValues := bytes.Split(row, []byte(":"))
      if len(headerValues) >= 2 {
        values := bytes.Join(headerValues[1:], []byte(":"))
        // Separate by semiColon
        params := bytes.Split(values, []byte(";"))
        headerValue := strings.TrimSpace(string(params[0]))
        switch strings.TrimSpace(string(headerValues[0])) {
          case "Content-Disposition": requestForm.ContentDisposition = headerValue
          case "Content-Type": requestForm.ContentType = headerValue
        }
        for _, p := range params[1:] {
          // Separate by equal
          keyValue := bytes.Split(p, []byte("="))
          if len(keyValue) < 2 {
            continue
          }
          value := strings.TrimSpace(string(keyValue[1]))
          switch strings.TrimSpace(string(keyValue[0])) {
            case "name": requestForm.Name = value
            case "filename": requestForm.FileName = value
          }
        }
      }
    }
    request.Forms = append(request.Forms, requestForm)
  }
}

func tailByte(input [][]byte) [][]byte {
  if len(input) >= 1 {
    return input[1:]
  } else {
    return input
  }
}

func initByte(input [][]byte) [][]byte {
  if len(input) >= 1 {
    return input[:len(input) - 1]
  } else {
    return input
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
