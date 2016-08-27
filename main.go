package main

import (
  "net/http"
  "github.com/icecrasher321/mjolnir"
)

func main() {
 var h http.Header
 var bod []byte
  mjolnir.Hammer("http://volist.herokuapp.com", "GET", h, bod, 100) // (URL, MethodType, Header, Body, Rate (ps))
}
