package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/segmentio/stats/v4"
	"github.com/segmentio/stats/v4/prometheus"

	"github.com/yerden/go-dpdk/eal"
)

var metricsEndpoint = flag.String("metrics", ":10010", "Specify listen address for Prometheus endpoint")

func main() {
	n, err := eal.Init(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	defer eal.Cleanup()
	defer eal.StopLcores()

	os.Args[n], os.Args = os.Args[0], os.Args[n:]
	flag.Parse()
	statsHandler := prometheus.DefaultHandler
	eng := stats.NewEngine("dpdk", statsHandler)
	app, err := NewApp(eng)
	if err != nil {
		panic(err)
	}

	retCh := make(chan error, len(app.Work))

	for lcore, pq := range app.Work {
		eal.ExecOnLcoreAsync(lcore, retCh, LcoreFunc(pq))
	}

	// stats report
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for t := range ticker.C {
			app.Stats.ReportAt(t)
		}
	}()

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", statsHandler)
		srv := &http.Server{
			Addr:    *metricsEndpoint,
			Handler: mux,
		}
		srv.ListenAndServe()
	}()

	err = <-retCh
	if e, ok := err.(*eal.ErrLcorePanic); ok {
		e.FprintStack(os.Stdout)
	}
	log.Println(err)
}
