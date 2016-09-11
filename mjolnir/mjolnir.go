package mjolnir

import (
  "fmt"
  "time"
  "net/http"
  "bufio"
  "os"
  vegeta "github.com/tsenart/vegeta/lib"
)

func Hammer(url string, method string, header http.Header, body []byte, r int) {
  f1, err := os.Create("graph.html")
  checkError(err)
  f2, err := os.Create("text.txt")
  checkError(err)
  rate := uint64(r) // per second
  duration := 5 * time.Second
  targeter := vegeta.NewStaticTargeter(vegeta.Target{
    Method: method,
    URL:    url,
    Body: body,
    Header: header,
  })
  attacker := vegeta.NewAttacker()

  w1 := bufio.NewWriter(f1)
  defer w1.Flush()
  w2 := bufio.NewWriter(f2)
  defer w2.Flush()

  var metrics vegeta.Metrics
  var results vegeta.Results

for res := range attacker.Attack(targeter, rate, duration) {
    metrics.Add(res)
    results.Add(res)
  }
  metrics.Close()
  results.Close()

  fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
  fmt.Printf("Mean: %s\n", metrics.Latencies.Mean)

  graphReporter := vegeta.NewPlotReporter(url, &results)
  graphReporter.Report(w1)
  textReporter := vegeta.NewTextReporter(&metrics)
  textReporter.Report(w2)
  
}

func checkError(e error) {
  if e != nil {
    fmt.Println(e)
    panic(e)
  }
}

