package ss5

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
