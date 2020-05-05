package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

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
	body := explodingMetric() + explodingLabels() + explodingLabelValues()
	fmt.Fprintf(w, body)
}

func main() {

	http.HandleFunc("/metrics", metrics)

	http.ListenAndServe(":8080", nil)
}
