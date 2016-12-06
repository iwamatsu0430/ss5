package ss5

import (
  "os"
  "path/filepath"
)

func FileServer(request Request) (response Response) {
  var targetPath string
  if (request.path == "/") {
    targetPath = config.Public.Path + "/" + config.Public.Index
  } else {
    targetPath = config.Public.Path + "/" + request.path
  }

  err := LoadFile(targetPath, &response.body)
  if err != nil {
    response.status = NOT_FOUND
    targetPath = config.Public.Path + "/" + config.Public.NotFound
    err := LoadFile(targetPath, &response.body)
    if err != nil {
      // TODO to const
      targetPath = "resources/views/defaults/404.html"
      err := LoadFile(targetPath, &response.body)
      if err != nil {
        // TODO break
      }
    }
  } else {
    response.status = OK
  }
  response.contentType = FindContentType(targetPath)

  return response
}

func LoadFile(filePath string, fileBytes *[]byte) error {
  file, err := os.Open(filePath)
  if err == nil {
    for {
      bytes := make([]byte, 1024)
      _, err := file.Read(bytes)
      *fileBytes = append(*fileBytes, bytes...)
      if err != nil {
        break
      }
    }
  }
  return err
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
