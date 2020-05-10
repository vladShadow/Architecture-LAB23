package main

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var clients = []string{
	"client/address/1",
	"another/client/address",
	"last/but/not/least/client/address",
}

func (s *MySuite) TestBalancer(c *C) {
	// TODO: Реалізуйте юніт-тест для балансувальникка.

	serversList = []string{
		"server1:8080",
		"server2:8080",
		"server3:8080",
	}

	serversPool = []string{
		"server1:8080",
		"server2:8080",
		"server3:8080",
	}

	// all servers alive
	var prevIndex int
	for i := 0; i < len(clients); i++ {
		for j := 0; i < 5; j++ {
			serverIndex := getIndexByClient(clients[i], len(serversPool))
			if j != 0 {
				c.Assert(serverIndex, Equals, prevIndex)
			}
			prevIndex = serverIndex
		}
	}
	// only one alive
	serversPool = []string{
		"server1:8080",
	}
	for i := 0; i < len(clients); i++ {
		for j := 0; i < 5; j++ {
			serverIndex := getIndexByClient(clients[i], len(serversPool))
			if j != 0 {
				c.Assert(serverIndex, Equals, prevIndex)
			}
			prevIndex = serverIndex
		}
	}
	// none alive
	serversPool = []string{}
	for i := 0; i < len(clients); i++ {
		for j := 0; i < 5; j++ {
			serverIndex := getIndexByClient(clients[i], len(serversPool))
			c.Assert(serverIndex, Equals, -1)
		}
	}
}
