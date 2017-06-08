package stats

import (
	"log"
	"net/http"
	"time"

	"github.com/datadog/datadog-go/statsd"
)

var client *statsd.Client

// Intialize the statsd package
func init() {
	var err error
	client, err = statsd.New("whatever")
	if err != nil {
		log.Fatal("Could not connect to statsd server : ", err)
	}
	log.Println("Statsd client initialized...")
}

func Submit() {
}

func Increment(key string) {
}

func Timer(key string, duration time.Duration) {
}

func IcrementHandler(fn http.HandlerFunc, key string) http.HandlerFunc {
	Increment(key)
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r)
	}
}

func GaugeHandler(fn http.HandlerFunc, key string) http.HandlerFunc {
	start := time.Now()
	defer func() {
		Timer(key, time.Since(start))
	}()

	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r)
	}
}
