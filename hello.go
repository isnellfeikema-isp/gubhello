package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	guber "github.com/mailgun/gubernator"
	"net/http"
	"os"
	"time"
)

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func main() {
	var (
		gubEndpoint = flag.String("gubEndpoint", "gubernator:81", "")
		name      = flag.String("name", "get_hello", "")
		uniqueKey = flag.String("uniqueKey", "default", "")
		hits      = flag.Int64("hits", 1, "")
		limit     = flag.Int64("limit", 1, "")
		duration  = flag.Int64("duration", 1000, "ms")
	)

	flag.Parse()

	client, err := guber.DialV1Server(*gubEndpoint)
	checkErr(err)

	http.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
		rateLimit := &guber.RateLimitReq{
			Name:      *name,
			UniqueKey: *uniqueKey,
			Hits:      *hits,
			Limit:     *limit,
			Duration:  *duration,
			Algorithm: guber.Algorithm_LEAKY_BUCKET,
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
		// Now hit our cluster with the rate limits
		rateLimitResp, err := client.GetRateLimits(ctx, &guber.GetRateLimitsReq{
			Requests: []*guber.RateLimitReq{rateLimit},
		})
		if err != nil {
			fmt.Println(err)
		}
		cancel()

		for _, r := range rateLimitResp.Responses {
			if r.Error != "" {
				panic(r.Error)
			}
			spew.Dump(r)
			if r.Status == guber.Status_OVER_LIMIT {
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
	})

	fmt.Println(http.ListenAndServe(":80", nil))
}
