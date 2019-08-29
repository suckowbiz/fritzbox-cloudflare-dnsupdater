package main

import (
	"log"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestRESTEndpointsAvailableIntegration(t *testing.T) {
	assert := assert.New(t)
	port := 61978

	tests := []struct {
		endpoint string
		method   string
	}{
		{
			endpoint: "http://localhost:" + strconv.Itoa(port) + "/update",
			method:   http.MethodPost,
		},
	}

	server := newServer(port, func(writer http.ResponseWriter, request *http.Request) {

	})
	go func() {
		log.Fatal(server.ListenAndServe())
	}()

	// Avoid: Post http://localhost:61978/update: dial tcp 127.0.0.1:61978: connect: connection refused
	time.Sleep(time.Second * 5)

	for _, test := range tests {
		req, _ := http.NewRequest(test.method, test.endpoint, nil)
		res, err := http.DefaultClient.Do(req)
		if res != nil {
			//nolint
			defer res.Body.Close()
		}
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.NotEqual(http.StatusNotFound, res.StatusCode)
	}
}
