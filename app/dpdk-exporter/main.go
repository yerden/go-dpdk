package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yerden/go-dpdk/eal"
)

var (
	addr     = flag.String("addr", ":22017", "Endpoint for prometheus")
	interval = flag.Duration("interval", time.Second, "Interval between statistics reports")
)

func main() {
	// gracefully shutdown
	shutdownCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGSEGV)
	defer stop()

	// initialize the EAL
	n, err := eal.Init(os.Args)
	if err != nil {
		log.Panicf("invalid EAL arguments: %v", err)
	}
	// clean up the EAL
	defer eal.Cleanup()
	defer eal.StopLcores()

	os.Args[n], os.Args = os.Args[0], os.Args[n:]
	flag.Parse()

	// if the port is already in use, the app will not start
	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Panicf("net listen: %v", err)
	}
	defer ln.Close()

	var metrics *Metrics
	if err := eal.ExecOnMain(func(lc *eal.LcoreCtx) {
		var err error
		metrics, err = NewMetrics()
		log.Panicf("init metrics collecting: %v", err)
	}); err != nil {
		log.Panicf("exec on main lcore: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		metrics.StartCollecting(shutdownCtx)
	}()

	// export metrics
	http.Handle("/metrics", promhttp.Handler())
	srv := &http.Server{Handler: http.DefaultServeMux}
	go func() {
		if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			log.Panicf("http server: %v", err)
		}
	}()

	log.Println("all workers started")
	wg.Wait()
	log.Printf("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shut down http server: %v", err)
	}
	log.Println("all workers stopped")
}
