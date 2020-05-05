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

func metrics(w http.ResponseWriter, req *http.Request) {

	var body = "label_explosion{exploding_label=\"" + randomString(8) + "\"} 1\n" +
		"metric_explosion_" + randomString(8) + "{label=\"this_is_fine\"} 1\n"

	fmt.Fprintf(w, body)
}

func main() {

	http.HandleFunc("/metrics", metrics)

	http.ListenAndServe(":8080", nil)
}
