package main

import (
	"context"
	"flag"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Tornado9966/Lab3_Go/httptools"
	"github.com/Tornado9966/Lab3_Go/signal"
)

var (
	port       = flag.Int("port", 8090, "load balancer port")
	timeoutSec = flag.Int("timeout-sec", 3, "request timeout time in seconds")
	https      = flag.Bool("https", false, "whether backends support HTTPs")

	traceEnabled = flag.Bool("trace", false, "whether to include tracing information into responses")

	timeout = time.Duration(*timeoutSec) * time.Second

	serversPool        map[int]string
	healthyServersPool []int
)

func scheme() string {
	if *https {
		return "https"
	}
	return "http"
}

func health(dst string) bool {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	req, _ := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("%s://%s/health", scheme(), dst), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	if resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}

func forward(dst string, rw http.ResponseWriter, r *http.Request) error {
	ctx, _ := context.WithTimeout(r.Context(), timeout)
	fwdRequest := r.Clone(ctx)
	fwdRequest.RequestURI = ""
	fwdRequest.URL.Host = dst
	fwdRequest.URL.Scheme = scheme()
	fwdRequest.Host = dst

	resp, err := http.DefaultClient.Do(fwdRequest)
	if err == nil {
		for k, values := range resp.Header {
			for _, value := range values {
				rw.Header().Add(k, value)
			}
		}
		if *traceEnabled {
			rw.Header().Set("lb-from", dst)
		}
		log.Println("fwd", resp.StatusCode, resp.Request.URL)
		rw.WriteHeader(resp.StatusCode)
		defer resp.Body.Close()
		_, err := io.Copy(rw, resp.Body)
		if err != nil {
			log.Printf("Failed to write response: %s", err)
		}
		return nil
	} else {
		log.Printf("Failed to get response from %s: %s", dst, err)
		rw.WriteHeader(http.StatusServiceUnavailable)
		return err
	}
}

func hash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return (int)(h.Sum32())
}

func find(slice []int, key int) (int, bool) {
	for i, item := range slice {
		if item == key {
			return i, true
		}
	}
	return -1, false
}

func remove(slice []int, key int) []int {
	if i, exist := find(slice, key); exist {
		return append(slice[:i], slice[i+1:]...)
	}
	return slice
}

func checkServer(server string, key int) {
	healthy := health(server)
	log.Println(server, healthy)
	if healthy {
		if _, exist := find(healthyServersPool, key); !exist {
			healthyServersPool = append(healthyServersPool, key)
		}
	} else {
		healthyServersPool = remove(healthyServersPool, key)
	}
}

func chooseServer(addr string) string {
	hash := hash(addr)
	if _, exist := find(healthyServersPool, hash%len(serversPool)); exist {
		return serversPool[hash%len(serversPool)]
	} else {
		return serversPool[healthyServersPool[hash%len(healthyServersPool)]]
	}
}

func main() {
	flag.Parse()

	serversPool = make(map[int]string)

	serversPool[0] = "server1:8080"
	serversPool[1] = "server2:8080"
	serversPool[2] = "server3:8080"

	checkServer("server1:8080", 0)
	checkServer("server2:8080", 1)
	checkServer("server3:8080", 2)

	for key, server := range serversPool {
		server := server
		key := key
		go func() {
			for range time.Tick(10 * time.Second) {
				checkServer(server, key)
			}
		}()
	}

	frontend := httptools.CreateServer(*port, http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if len(healthyServersPool) != 0 {
			forward(chooseServer(r.RemoteAddr), rw, r)
		} else {
			log.Println("All servers are busy. Wait please.")
		}
	}))

	log.Println("Starting load balancer...")
	log.Printf("Tracing support enabled: %t", *traceEnabled)
	frontend.Start()
	signal.WaitForTerminationSignal()
}
