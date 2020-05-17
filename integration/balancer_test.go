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

var clients = []http.Client{
	{Timeout: 3 * time.Second},
	{Timeout: 3 * time.Second},
	{Timeout: 3 * time.Second},
	{Timeout: 3 * time.Second},
	{Timeout: 3 * time.Second},
}
var requests = []int{0, 0, 0, 0, 0}

var serversList = []string{
	"server1:8080",
	"server2:8080",
	"server3:8080",
}

func indexOf(arr []string, str string) int {
	for i, n := range arr {
		if str == n {
			return i
		}
	}
	return -1
}
func (s *MySuite) TestBalancer(c *C) {
	// TODO: Реалізуйте інтеграційний тест для балансувальникка.
	counter := 0
	for range time.Tick(3 * time.Second) {
		for i := 0; i < len(clients); i++ {
			res, err := clients[i].Get(fmt.Sprintf("%s/api/v1/some-data", baseAddress))
			if err != nil {
				c.Error(err)
			}
			index := indexOf(serversList, res.Header.Get("lb-from"))
			if index != -1 {
				requests[index]++
			}
			fmt.Println("Server \"", res.Header.Get("lb-from"), "\" gets request from client number ", i)
		}

		counter++
		if counter >= 10 {
			break
		}
	}
}

func BenchmarkBalancer(b *testing.B) {
	// TODO: Реалізуйте інтеграційний бенчмарк для балансувальникка
	for k := 0; k < b.N; k++ {
		for i := 0; i < len(clients); i++ {
			_, err := clients[i].Get(fmt.Sprintf("%s/api/v1/some-data", baseAddress))
			if err != nil {
				b.Error(err)
			}
		}
	}
}
