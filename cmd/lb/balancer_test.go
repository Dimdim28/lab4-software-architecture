package main

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

type Server struct {
	URL     string
	ConnCnt int
	Healthy bool
}

var serversPool []Server

func TestFindMinServer(t *testing.T) {
	assert := assert.New(t)

	t.Run("No healthy servers", func(t *testing.T) {
		serversPool = []Server{
			{URL: "Server1", ConnCnt: 10, Healthy: false},
			{URL: "Server2", ConnCnt: 20, Healthy: false},
			{URL: "Server3", ConnCnt: 30, Healthy: false},
		}
		assert.Equal(-1, FindMinServer())
	})

	t.Run("All healthy servers", func(t *testing.T) {
		serversPool = []Server{
			{URL: "Server1", ConnCnt: 10, Healthy: true},
			{URL: "Server2", ConnCnt: 20, Healthy: true},
			{URL: "Server3", ConnCnt: 30, Healthy: true},
		}
		assert.Equal(0, FindMinServer())
	})

	t.Run("Mixed healthy and unhealthy servers", func(t *testing.T) {
		serversPool = []Server{
			{URL: "Server1", ConnCnt: 10, Healthy: false},
			{URL: "Server2", ConnCnt: 20, Healthy: true},
			{URL: "Server3", ConnCnt: 30, Healthy: true},
		}
		assert.Equal(1, FindMinServer())
	})

	t.Run("Minimum connection count", func(t *testing.T) {
		serversPool = []Server{
			{URL: "Server1", ConnCnt: 10, Healthy: true},
			{URL: "Server2", ConnCnt: 5, Healthy: true},
			{URL: "Server3", ConnCnt: 30, Healthy: true},
		}
		assert.Equal(1, FindMinServer())
	})
}

// func Test_findMinServer(t *testing.T) {
// 	tests := []struct {
// 		server1 *Server
// 		server2 *Server
// 		server3 *Server
// 	}{
// 		{
// 			server1: {URL: "server1:8080", ConnCnt: 0, Healthy: true},
// 			server2: {URL: "server2:8080", ConnCnt: 0, Healthy: true},
// 			server3: {URL: "server3:8080", ConnCnt: 0, Healthy: true},
// 		},
// 		{
// 			server1: {URL: "server1:8080", ConnCnt: 5, Healthy: true},
// 			server2: {URL: "server2:8080", ConnCnt: 9, Healthy: true},
// 			server3: {URL: "server3:8080", ConnCnt: 3 Healthy: true},
// 		},
// 		{
// 			server1: {URL: "server1:8080", ConnCnt: 1, Healthy: false},
// 			server2: {URL: "server2:8080", ConnCnt: 4, Healthy: true},
// 			server3: {URL: "server3:8080", ConnCnt: 5, Healthy: true},
// 		},
// 		{
// 			server1: {URL: "server1:8080", ConnCnt: 5, Healthy: false},
// 			server2: {URL: "server2:8080", ConnCnt: 5, Healthy: false},
// 			server3: {URL: "server3:8080", ConnCnt: 5, Healthy: false},
// 		}
// 	}

// 	for _, tc := range tests {
// 		t.Run("test", func(t *testing.T) {
// 			assert.Equal()
// 		})
// 	}
// }
