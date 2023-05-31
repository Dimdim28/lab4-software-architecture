package main

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestFindMinServer(t *testing.T) {
	assert := assert.New(t)

	t.Run("No healthy servers", func(t *testing.T) {
		serversPool = []*Server{
			{URL: "Server1", ConnCnt: 10, Healthy: false},
			{URL: "Server2", ConnCnt: 20, Healthy: false},
			{URL: "Server3", ConnCnt: 30, Healthy: false},
		}
		assert.Equal(-1, FindMinServer())
	})

	t.Run("All healthy servers", func(t *testing.T) {
		serversPool = []*Server{
			{URL: "Server1", ConnCnt: 10, Healthy: true},
			{URL: "Server2", ConnCnt: 20, Healthy: true},
			{URL: "Server3", ConnCnt: 30, Healthy: true},
		}
		assert.Equal(0, FindMinServer())
	})

	t.Run("Mixed healthy and unhealthy servers", func(t *testing.T) {
		serversPool = []*Server{
			{URL: "Server1", ConnCnt: 10, Healthy: false},
			{URL: "Server2", ConnCnt: 20, Healthy: true},
			{URL: "Server3", ConnCnt: 30, Healthy: true},
		}
		assert.Equal(1, FindMinServer())
	})

	t.Run("Minimum connection count", func(t *testing.T) {
		serversPool = []*Server{
			{URL: "Server1", ConnCnt: 10, Healthy: true},
			{URL: "Server2", ConnCnt: 5, Healthy: true},
			{URL: "Server3", ConnCnt: 30, Healthy: true},
		}
		assert.Equal(1, FindMinServer())
	})
}

func TestHealth(t *testing.T) {
	mockURL := "http://example.com/Health"
	httpmock.RegisterResponder(http.MethodGet, mockURL, httpmock.NewStringResponder(http.StatusOK, ""))

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	server := &Server{
		URL: "example.com",
	}

	result := Health(server)

	assert.True(t, result)
	assert.True(t, server.Healthy)

	server.Healthy = false // скинув перед некст тестом

	httpmock.RegisterResponder(http.MethodGet, mockURL, httpmock.NewStringResponder(http.StatusInternalServerError, ""))
	result2 := Health(server)

	assert.False(t, result2)
	assert.False(t, server.Healthy)
}
