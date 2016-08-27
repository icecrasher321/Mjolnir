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
  f, err := os.Create("results.html")
  checkError(err)
  rate := uint64(r) // per second
  duration := 1 * time.Second
  targeter := vegeta.NewStaticTargeter(vegeta.Target{
    Method: method,
    URL:    url,
    Body: body,
    Header: header,
  })
  attacker := vegeta.NewAttacker()

  w := bufio.NewWriter(f)
  defer w.Flush()
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

  reporter := vegeta.NewPlotReporter("Results", &results)
  reporter.Report(w)
}

func checkError(e error) {
  if e != nil {
    fmt.Println(e)
    panic(e)
  }
}
