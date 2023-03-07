package main

import (
	"flag"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	addr = flag.String("addr", "http://localhost:3090", "API server address")
)

func TestOk(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	u, err := url.Parse(*addr)
	assert.NoError(t, err)

	u.Path = "/ok"
	resp, err := http.Get(u.String())
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
