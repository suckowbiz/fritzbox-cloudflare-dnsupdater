package main

import (
	"log"
	"net/http"
	"strconv"
	"testing"

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
		//nolint
		defer server.Close()
	}()

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
