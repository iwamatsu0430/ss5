package ss5

import (
  "fmt"
  "os"
  "path/filepath"
)

func FileServer(request Request) (response Response) {
  // TODO read from config
  targetFile := "public"
  if (request.path == "/") {
    // TODO read from config
    targetFile += "/index.html"
  } else {
    targetFile += request.path
  }
  file, err := os.Open(targetFile)

  if err == nil {
    response.status = OK
    response.contentType = FindContentType(targetFile)
    for {
      bytes := make([]byte, 1024)
      _, err := file.Read(bytes)
      response.body = append(response.body, bytes...)
      if err != nil {
        break
      }
    }
  } else {
    // TODO route by config
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

  return response
}

func FindContentType(filePath string) string {
  ext := filepath.Ext(filePath)
  switch ext {
    case ".html", ".htm": return HTML
    case ".csv": return CSV
    case ".js": return JS
    case ".json": return JSON
    case ".jpg", ".jpeg": return JPG
    case ".png": return PNG
    case ".gif": return GIF
    default: return PLAIN
  }
}
