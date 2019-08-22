package main

import (
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
		log.Fatal("environment variable FCD_PORT must contain a port number to run the REST endpoint")
	}
	dnsAer := cloudflare.NewDNSA()
	handler := fritzbox.NewUpdateHandler(api.NewWithAPIToken, dnsAer)
	s := newServer(port, handler)
	log.Fatal(s.ListenAndServe())
}

func getenvInt(key string) (int, error) {
	envStr := os.Getenv(key)
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
