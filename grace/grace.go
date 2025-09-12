package grace

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func ServeHTTP(addr string, handler http.Handler) error {
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	return ServeServer(server)
}

func ServeHTTPS(addr, certFile, keyFile string, handler http.Handler) error {
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	return ServeServerTLS(server, certFile, keyFile)
}

func ServeServer(server *http.Server) error {
	go func() {
		log.Printf("Starting HTTP server on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()
	return waitForShutdown(server)
}

func ServeServerTLS(server *http.Server, certFile, keyFile string) error {
	go func() {
		log.Printf("Starting HTTPS server on %s", server.Addr)
		if err := server.ListenAndServeTLS(certFile, keyFile); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTPS server error: %v", err)
		}
	}()
	return waitForShutdown(server)
}

func waitForShutdown(server *http.Server) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown signal received...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced shutdown: %v", err)
		return err
	}

	log.Println("Server gracefully stopped")
	return nil
}
