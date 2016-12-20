package ss5

type Request struct {
  Method string
  Path string
  Version string
  Headers map[string]string
  ContentType string
  ContentLength int
  Boundary string
  Body []byte
  Forms []RequestForm
}

type RequestForm struct {
  Name string
  FileName string
  ContentDisposition string
  ContentType string
  Body []byte
}

type Response struct {
  Status string
  ContentType string
  Headers map[string]string
  Body []byte
}

type Config struct {
  Public PublicConfig
}

type PublicConfig struct {
  Path      string
  Index     string
  NotFound  string
}
