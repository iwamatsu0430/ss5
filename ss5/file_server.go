package ss5

import (
  "os"
  "path/filepath"
)

func FileServer(request Request) (response Response) {
  var targetPath string
  if (request.Path == "/") {
    targetPath = config.Public.Path + "/" + config.Public.Index
  } else {
    targetPath = config.Public.Path + "/" + request.Path
  }

  err := LoadFile(targetPath, &response.Body)
  if err != nil {
    response.Status = HTTP_STATUS_404
    targetPath = config.Public.Path + "/" + config.Public.NotFound
    err := LoadFile(targetPath, &response.Body)
    if err != nil {
      targetPath = SS5_DEFAULT_404
      err := LoadFile(targetPath, &response.Body)
      if err != nil {
        Exit("ss5 default 404 page not found!")
      }
    }
  } else {
    response.Status = HTTP_STATUS_200
  }
  response.ContentType = FindContentType(targetPath)

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
    case ".html", ".htm": return CONTENT_TYPE_HTML
    case ".csv": return CONTENT_TYPE_CSV
    case ".js": return CONTENT_TYPE_JS
    case ".json": return CONTENT_TYPE_JSON
    case ".jpg", ".jpeg": return CONTENT_TYPE_JPG
    case ".png": return CONTENT_TYPE_PNG
    case ".gif": return CONTENT_TYPE_GIF
    default: return CONTENT_TYPE_PLAIN
  }
}
