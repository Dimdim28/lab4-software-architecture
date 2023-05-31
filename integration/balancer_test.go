package integration

import (
	"fmt"
	"net/http"
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
		resp, err := client.Get(fmt.Sprintf("%s/api/v1/some-data", baseAddress))
		assert.NoError(s.T(), err)

		if i%3 == 0 {
			serverNum = 1
		} else if i%3 == 1 {
			serverNum = 2
		} else {
			serverNum = 3
		}
		assert.Equal(s.T(), fmt.Sprintf("server%d:8080", serverNum), resp.Header.Get("lb-from"))
	}
}

func (s *IntegrationTestSuite) BenchmarkBalancer(b *testing.B) {
	if _, exists := os.LookupEnv("INTEGRATION_TEST"); !exists {
		s.T().Skip("Integration test is not enabled")
	}

	for i := 0; i < b.N; i++ {
		_, err := client.Get(fmt.Sprintf("%s/api/v1/some-data", baseAddress))
		assert.NoError(s.T(), err)
	}
}
