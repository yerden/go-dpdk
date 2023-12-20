package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yerden/go-dpdk/eal"
	"github.com/yerden/go-dpdk/ethdev"
)

var metricsEndpoint = flag.String("metrics", ":10010", "Specify listen address for Prometheus endpoint")
var fcMode FcModeFlag

func main() {
	n, err := eal.Init(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	defer eal.Cleanup()
	defer eal.StopLcores()

	os.Args[n], os.Args = os.Args[0], os.Args[n:]
	flag.Var(&fcMode, "flowctrl", "Specify Flow Control mode: none (default), rxpause, txpause, full")

	flag.Parse()
	reg := prometheus.NewRegistry()
	ethdev.RegisterTelemetryLSC("/ethdev/lsc")
	app, err := NewApp(reg)
	if err != nil {
		panic(err)
	}

	retCh := make(chan error, len(app.Work))

	for lcore, pq := range app.Work {
		eal.ExecOnLcoreAsync(lcore, retCh, LcoreFunc(pq, app.QCR))
	}

	// stats report
	go func() {
		ticker := time.NewTicker(*statsInt)
		defer ticker.Stop()

		for range ticker.C {
			app.Stats.Report()
		}
	}()

	go func() {
		for err := range retCh {
			if err == nil {
				continue
			}
			if e, ok := err.(*eal.ErrLcorePanic); ok {
				e.FprintStack(os.Stdout)
			}
			log.Println(err)
		}
	}()

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	srv := &http.Server{
		Addr:    *metricsEndpoint,
		Handler: mux,
	}
	log.Println(srv.ListenAndServe())
}
