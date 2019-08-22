package main

import (
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRESTEndpointsAvailableIntegration(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		endpoint string
		method   string
	}{
		{
			endpoint: "http://localhost:1978/update",
			method:   http.MethodPost,
		},
	}

	server := newServer(1978, func(writer http.ResponseWriter, request *http.Request) {

	})
	go func() {
		log.Fatal(server.ListenAndServe())
		defer server.Close()
	}()

	for _, test := range tests {
		req, _ := http.NewRequest(test.method, test.endpoint, nil)
		res, err := http.DefaultClient.Do(req)
		assert.NotNil(res)
		assert.NoError(err)
		assert.NotEqual(http.StatusNotFound, res.StatusCode)
	}
}
