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
	var client1 = http.Client{Timeout: 3 * time.Second}
	var client2 = http.Client{Timeout: 3 * time.Second}
	var client3 = http.Client{Timeout: 3 * time.Second}

	count := 0
	count1 := 0
	count2 := 0
	count3 := 0

	for range time.Tick(3 * time.Second) {

		resp1, err1 := client1.Get(fmt.Sprintf("%s/api/v1/some-data", baseAddress))
		resp2, err2 := client2.Get(fmt.Sprintf("%s/api/v1/some-data", baseAddress))
		resp3, err3 := client3.Get(fmt.Sprintf("%s/api/v1/some-data", baseAddress))

		if err1 != nil {
			c.Error(err1)
		}
		if err2 != nil {
			c.Error(err2)
		}
		if err3 != nil {
			c.Error(err3)
		}

		if resp1.Header.Get("lb-from") == "server1:8080" {
			count1++
		}
		if resp1.Header.Get("lb-from") == "server2:8080" {
			count2++
		}
		if resp1.Header.Get("lb-from") == "server3:8080" {
			count3++
		}
		if resp2.Header.Get("lb-from") == "server1:8080" {
			count1++
		}
		if resp2.Header.Get("lb-from") == "server2:8080" {
			count2++
		}
		if resp2.Header.Get("lb-from") == "server3:8080" {
			count3++
		}
		if resp3.Header.Get("lb-from") == "server1:8080" {
			count1++
		}
		if resp3.Header.Get("lb-from") == "server2:8080" {
			count2++
		}
		if resp3.Header.Get("lb-from") == "server3:8080" {
			count3++
		}

		fmt.Println("1 client: ", resp1.Header.Get("lb-from"))
		fmt.Println("2 client: ", resp2.Header.Get("lb-from"))
		fmt.Println("3 client: ", resp3.Header.Get("lb-from"))

		count++
		if count == 20 {
			break
		}
	}
	fmt.Println("Responses from server 1: ", count1)
	fmt.Println("Responses from server 2: ", count2)
	fmt.Println("Responses from server 3: ", count3)
}

func BenchmarkBalancer(b *testing.B) {
	// TODO: Реалізуйте інтеграційний бенчмарк для балансувальникка
	var client1 = http.Client{Timeout: 3 * time.Second}
	var client2 = http.Client{Timeout: 3 * time.Second}
	var client3 = http.Client{Timeout: 3 * time.Second}

	for i := 0; i < b.N; i++ {
		_, err1 := client1.Get(fmt.Sprintf("%s/api/v1/some-data", baseAddress))
		_, err2 := client2.Get(fmt.Sprintf("%s/api/v1/some-data", baseAddress))
		_, err3 := client3.Get(fmt.Sprintf("%s/api/v1/some-data", baseAddress))

		if err1 != nil {
			b.Error(err1)
		}
		if err2 != nil {
			b.Error(err2)
		}
		if err3 != nil {
			b.Error(err3)
		}
	}
}
