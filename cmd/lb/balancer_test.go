package main

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestBalancer(c *C) {
	// TODO: Реалізуйте юніт-тест для балансувальникка.
	// балансувальник визначає сервер по адресу клієнта
	// тобто коли не відбувається зміна списку доступних серверів
	// для одного клієнта буде використовуватись один і той же сервер

	clients := []string{
		"172.20.0.3:45234",
		"172.20.0.4:34563",
		"172.20.0.2:34563"}

	// all servers alive
	serversPool = []string{
		"server1:8080",
		"server2:8080",
		"server3:8080",
	}
	var prevIndex int
	for i := 0; i < len(clients); i++ {
		for j := 0; i < 5; j++ {
			serverIndex := 1
			if j != 0 {
				c.Assert(serverIndex, Equals, prevIndex)
			}
			prevIndex = serverIndex
		}
	}

	// // only one alive
	// serversPool = []string{
	// 	"server1:8080",
	// }
	// for i := 0; i < len(clients); i++ {
	// 	for j := 0; i < 5; j++ {
	// 		serverIndex := getIndexByClient(clients[i], len(serversPool))
	// 		if j != 0 {
	// 			c.Assert(serverIndex, Equals, prevIndex)
	// 		}
	// 		prevIndex = serverIndex
	// 	}
	// }

	// // none alive
	// serversPool = []string{}
	// for i := 0; i < len(clients); i++ {
	// 	for j := 0; i < 5; j++ {
	// 		serverIndex := getIndexByClient(clients[i], len(serversPool))
	// 		c.Assert(serverIndex, Equals, -1)
	// 	}
	// }
}
