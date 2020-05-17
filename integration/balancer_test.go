package integration

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

const baseAddress = "http://balancer:8090"

var client = http.Client{
	Timeout: 3 * time.Second,
}

func (s *MySuite) TestBalancer(c *C) {
	// TODO: Реалізуйте інтеграційний тест для балансувальникка.
	resp, err := client.Get(fmt.Sprintf("%s/api/v1/some-data", baseAddress))
	if err != nil {
		c.Error(err)
	}
	c.Logf("response from [%s]", resp.Header.Get("lb-from"))
}

func BenchmarkBalancer(b *testing.B) {
	// TODO: Реалізуйте інтеграційний бенчмарк для балансувальникка.
}
