package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/holger-hoffmann/dev.httpfileserver/service"
)

func main() {
	var (
		bindAddress = flag.String("bind-address", "0.0.0.0", "Address the server will bind to.")
		port        = flag.String("port", "8080", "The port the server will listen on.")
		dir         = flag.String("dir", "/tmp", "The directory from where the files are served and where to they are uploaded.")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		flag.PrintDefaults()
		return
	}
	flag.Parse()

	if os.Getenv("PORT") != "" {
		*port = os.Getenv("PORT")
	}

	fmt.Printf("dev.httpfileserver - Starting.\n")
	fmt.Printf("dev.httpfileserver -   bind-address: %s\n", *bindAddress)
	fmt.Printf("dev.httpfileserver -   port:         %s\n", *port)
	fmt.Printf("dev.httpfileserver -   dir:          %s\n", *dir)

	ctx := context.Background()

	s, err := service.New(ctx, *dir)
	if err != nil {
		log.Printf("dev.httpfileserver - Failed to initialize service: %s", err)
		log.Fatalf("dev.httpfileserver - Exiting.")
	}

	srv := &http.Server{
		Addr:    net.JoinHostPort(*bindAddress, *port),
		Handler: s.GetHandler(),
	}
	errChan := make(chan error, 10)
	go func() {
		fmt.Printf("dev.httpfileserver - Listening on %s.\n", net.JoinHostPort(*bindAddress, *port))
		err := srv.ListenAndServe()
		if err != nil {
			errChan <- err
			return
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case err := <-errChan:
			if err != nil {
				log.Printf("dev.httpfileserver - Error while serving, exiting: %s", err)
				os.Exit(1)
			}
		case <-signalChan:
			fmt.Println()
			log.Printf("dev.httpfileserver - Shutdown signal received, stopping gracefully.")
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			srv.Shutdown(ctx)
			log.Printf("dev.httpfileserver - Exiting.")
			os.Exit(0)
		}
	}
}
