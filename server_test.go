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

	srv, err := NewServer(23456, r, WithBindAddress("127.0.0.1"))
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

func TestServer_StartAsync_Blocking(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("bar"))
	}).Methods(http.MethodGet)

	srv, err := NewServer(34567, r, WithBlockingLaunch(true), WithBindAddress("127.0.0.1"))
	require.NoError(t, err)
	defer srv.Stop()

	require.NoError(t, srv.StartAsync())

	resp, err := http.Get("http://localhost:34567/foo")
	require.NoError(t, err)

	b, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "bar", string(b))
}
