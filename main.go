package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var (
	tempo                = flag.Int("tempo", 2, "define number of metric groups exposed in one request")
	samplesExplosion     = flag.Bool("explode-samples-scraped", true, "Use if you want to cause prometheus `scrape_samples_scraped` metric to constantly increase. It exposes increasing number of random metrics.")
	labelsNameExplosion  = flag.Bool("label-name-explosion", true, "Expose static number of metrics (defined by `-tempo`), but change label name on every resuest.")
	labelsValueExplosion = flag.Bool("label-value-explosion", true, "Expose static number of metrics (defined by `-tempo`), but change label value on every resuest.")
)

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
var counter = 0

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// This should cause prometheus `scrape_samples_scraped` metric to constantly increase
func explodingSamples() string {
	buffer := ""
	for i := 0; i < counter; i++ {
		c := strconv.Itoa(i)
		buffer += "eagle_sample_" + c + "{label=\"this_is_fine\"} 1\n"
	}
	counter++
	return buffer
}

func explodingLabels() string {
	return "eagle_labels{" + randomString(8) + "=\"this_is_fine\"} 1\n"
}

func explodingLabelValues() string {
	return "eagle_label_values{this_is_fine=\"" + randomString(8) + "\"} 1\n"
}

func metrics(w http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Host)
	body := ""
	for i := 0; i < *tempo; i++ {
		if *samplesExplosion {
			body += explodingSamples()
		}
		if *labelsNameExplosion {
			body += explodingLabels()
		}
		if *labelsValueExplosion {
			body += explodingLabelValues()
		}
	}
	fmt.Fprintf(w, body)
}

func main() {
	flag.Parse()

	http.HandleFunc("/metrics", metrics)

	http.ListenAndServe(":8080", nil)
}
