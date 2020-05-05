package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

var tempo = flag.Int("tempo", 2, "define number of metric groups exposed in one request")
var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func explodingMetric() string {
	return "eagle_" + randomString(8) + "{label=\"this_is_fine\"} 1\n"
}

func explodingLabels() string {
	return "eagle_labels{" + randomString(8) + "=\"this_is_fine\"} 1\n"
}

func explodingLabelValues() string {
	return "eagle_label_values{label=\"" + randomString(8) + "\"} 1\n"
}

func metrics(w http.ResponseWriter, req *http.Request) {
	body := ""
	for i := 0; i < *tempo; i++ {
		body += explodingMetric() + explodingLabels() + explodingLabelValues()
	}
	fmt.Fprintf(w, body)
}

func main() {
	flag.Parse()

	http.HandleFunc("/metrics", metrics)

	http.ListenAndServe(":8080", nil)
}
