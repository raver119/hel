package hel

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestServer_StartAsync(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("bar"))
	}).Methods(http.MethodGet)

	srv, err := NewServer(23456, r)
	require.NoError(t, err)

	require.NoError(t, srv.StartAsync())
	time.Sleep(2 * time.Second)

	// checking registered url
	resp, err := http.Get("http://localhost:23456/foo")
	require.NoError(t, err)

	b, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "bar", string(b))

	// and non-existent url must fail
	resp, err = http.Get("http://localhost:23456/boo")
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	require.NoError(t, srv.Stop())
}
