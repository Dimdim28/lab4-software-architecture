package integration

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const baseAddress = "http://balancer:8090"

var client = http.Client{
	Timeout: 3 * time.Second,
}

const teamName = "procrastinatioin"

func getData(key string) (*http.Response, error) {
	path := fmt.Sprintf("%s/api/v1/some-data", baseAddress)

	queryParams := url.Values{}
	queryParams.Set("key", key)
	path += "?" + queryParams.Encode()

	return client.Get(path)
}

type IntegrationTestSuite struct {
	suite.Suite
}

func TestBalancer(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) TestGetRequest() {
	if _, exists := os.LookupEnv("INTEGRATION_TEST"); !exists {
		s.T().Skip("Integration test is not enabled")
	}

	serverNum := 0
	for i := 0; i < 10; i++ {
		resp, err := getData(teamName)
		assert.NoError(s.T(), err)

		if i%3 == 0 {
			serverNum = 1
		} else if i%3 == 1 {
			serverNum = 2
		} else {
			serverNum = 3
		}
		assert.Equal(s.T(), fmt.Sprintf("server%d:8080", serverNum), resp.Header.Get("lb-from"))

		resp, err = getData("procrastination")
		assert.Equal(s, http.StatusNotFound, resp.StatusCode)
	}
}



func (s *IntegrationTestSuite) BenchmarkBalancer(b *testing.B) {
	if _, exists := os.LookupEnv("INTEGRATION_TEST"); !exists {
		s.T().Skip("Integration test is not enabled")
	}

	for i := 0; i < b.N; i++ {
		_, err := getData(teamName)
		assert.NoError(s.T(), err)
	}
}
