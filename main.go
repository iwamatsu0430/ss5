package main

import (
  "flag"
  "strconv"
  "./ss5"
)

func main() {
  flag.Parse()
  server := flag.Arg(0)
  portStr := flag.Arg(1)
  port, _ := strconv.Atoi(portStr)
  ss5.Start(server, port)
}
