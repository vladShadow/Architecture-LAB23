package main

import (
	"context"
	"flag"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/vladShadow/Architecture-LAB23/httptools"
	"github.com/vladShadow/Architecture-LAB23/signal"
)

var (
	port       = flag.Int("port", 8090, "load balancer port")
	timeoutSec = flag.Int("timeout-sec", 3, "request timeout time in seconds")
	https      = flag.Bool("https", false, "whether backends support HTTPs")

	traceEnabled = flag.Bool("trace", false, "whether to include tracing information into responses")
)

var (
	timeout     = time.Duration(*timeoutSec) * time.Second
	serversList = []string{
		"server1:8080",
		"server2:8080",
		"server3:8080",
	}
	serversPool = []string{}
	poolMutex   sync.Mutex
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

func indexOf(arr []string, str string) int {
	for i, n := range arr {
		if str == n {
			return i
		}
	}
	return -1
}

func getIndexByClient(addr string) int {
	len := len(serversPool)
	if len == 0 {
		log.Println("Failed to process the request: All servers are dead")
		return -1
	}
	poolIdx := hash(addr) % len
	idx := indexOf(serversList, serversPool[poolIdx])
	return idx
}

func hash(str string) int {
	temp := fnv.New32a()
	temp.Write([]byte(str))
	return int(temp.Sum32())
}

func main() {
	flag.Parse()

	// TODO: Використовуйте дані про стан сервреа, щоб підтримувати список тих серверів, яким можна відправляти ззапит.
	for _, server := range serversList {
		// ітерація по всіх серверах та підтримка пулу лише доступних серверів
		server := server
		go func() {
			for range time.Tick(10 * time.Second) {
				serverAvailable := health(server)
				idx := indexOf(serversPool, server)
				if serverAvailable && idx == -1 {
					serversPool = append(serversPool, server)
				}
				if !serverAvailable && idx != -1 {
					lastIdx := len(serversPool) - 1
					serversPool[idx] = serversPool[lastIdx]
					serversPool[lastIdx] = ""
				}
				log.Println(server, serverAvailable)
			}
		}()
	}

	frontend := httptools.CreateServer(*port, http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// TODO: Рееалізуйте свій алгоритм балансувальника.
		serverIndex := getIndexByClient(r.RemoteAddr)
		// індекс у повному списку серверів
		log.Println("serverIndex ", serverIndex)
		forward(serversList[serverIndex], rw, r)
	}))

	log.Println("Starting load balancer...")
	log.Printf("Tracing support enabled: %t", *traceEnabled)
	frontend.Start()
	signal.WaitForTerminationSignal()
}
