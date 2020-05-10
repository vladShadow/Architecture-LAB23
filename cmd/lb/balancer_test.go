package main

import (
	"math"
	"math/rand"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

func (s *MySuite) TestBalancer(c *C) {
	// TODO: Реалізуйте юніт-тест для балансувальникка.
	// перевірка збалансованості навантаження на сервери
	iterations := 1000
	serversAvailable := len(serversPool)
	serversLoads := make([]int, len(serversPool), serversAvailable)

	for i := 0; i < iterations; i++ {
		builder := make([]byte, 10)
		for l := range builder {
			builder[l] = byte(rand.Intn(255))
		}
		clientAddr := string(builder)
		serverIdx := indexOf(serversPool, serversList[getIndexByClient(clientAddr, serversAvailable)])
		serversLoads[serverIdx]++
	}

	expectedLoad := 1.0 / float64(serversAvailable)
	for _, number := range serversLoads {
		er := math.Abs(expectedLoad - float64(number)/float64(iterations))
		c.Assert(er > 0.05, Equals, true)
	}

}
