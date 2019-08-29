package main

import (
	"fmt"
	"fritzbox-cloudflare-dnsupdater/internal/cloudflare"
	"fritzbox-cloudflare-dnsupdater/internal/fritzbox"
	"log"
	"net/http"
	"os"
	"strconv"

	api "github.com/cloudflare/cloudflare-go"
)

func main() {
	port, err := getenvInt("FCD_PORT")
	if err != nil {
		log.Fatalf("webserver startup requires a valid port number. Please set environment variable 'FCD_PORT' "+
			"accordingly. %v", err)
	}
	dnsa := cloudflare.NewDNSA()
	handler := fritzbox.NewUpdateHandler(api.NewWithAPIToken, dnsa)
	server := newServer(port, handler)
	log.Fatal(server.ListenAndServe())
}

func getenvInt(key string) (int, error) {
	envStr := os.Getenv(key)
	if envStr == "" {
		return 0, fmt.Errorf("'%s' must not be empty", key)
	}
	envInt, err := strconv.Atoi(envStr)
	if err != nil {
		return 0, err
	}
	return envInt, nil
}

func newServer(port int, handler func(http.ResponseWriter, *http.Request)) *http.Server {
	log.Printf("starting server at 0.0.0.0:%d", port)
	http.HandleFunc("/update", handler)
	return &http.Server{Addr: ":" + strconv.Itoa(port)}
}
