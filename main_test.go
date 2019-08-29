package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestRESTEndpointsAvailableIntegration(t *testing.T) {
	port := 61978

	tests := []struct {
		endpoint string
		method   string
	}{
		{
			endpoint: "http://localhost:" + strconv.Itoa(port) + "/update",
			method:   http.MethodGet,
		},
	}

	server := newServer(port, func(writer http.ResponseWriter, request *http.Request) {

	})
	go func() {
		log.Fatal(server.ListenAndServe())
	}()

	// Avoid: http://localhost:61978/update: dial tcp 127.0.0.1:61978: connect: connection refused
	waitForConnection(t, "http://localhost:"+strconv.Itoa(port))

	for _, test := range tests {
		req, _ := http.NewRequest(test.method, test.endpoint, nil)
		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.NotEqual(t, http.StatusNotFound, res.StatusCode)
		_ = res.Body.Close()
	}
}

func waitForConnection(t *testing.T, url string) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)
	for {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("%s", err.Error())
			if strings.Contains(err.Error(), "connection refused") {
				continue
			}
			require.NoError(t, err)
		}
		_ = resp.Body.Close()
		break
	}
}
