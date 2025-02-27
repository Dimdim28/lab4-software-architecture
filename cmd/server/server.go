package main

import (
	"context"
	"flag"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
	"io/ioutil"

	"github.com/Dimdim28/lab4-software-architecture/httptools"
	"github.com/Dimdim28/lab4-software-architecture/signal"
)

var (
	port       = flag.Int("port", 8080, "server port")
	delay      = flag.Int("delay", 0, "response delay in millseconds")
	healthInit = flag.Bool("health", true, "initial server health")
	debug      = flag.Bool("debug", false, "whether we can change server's health status")
	dbUrl      = flag.String("db-url", "db:8100", "hostname of database service")
	report		 = make(Report)
)

const (
	scheme = "http"
  teamName = "procrastinaiton"
)

type boolMutex struct {
	mu sync.Mutex
	v  bool
}

func (c *boolMutex) Inverse() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.v = !c.v
}

func (c *boolMutex) Get() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.v
}

func main() {
	flag.Parse()
	h := new(http.ServeMux)
	health := boolMutex{v: *healthInit}
	createTeam()

	h.Handle("/health", healthHandler(&health))
	if *debug {
		h.Handle("/inverse-health", healthInverseHandler(&health))
	}

	h.Handle("/api/v1/some-data", http.HandlerFunc(handleDefaultGet))
	h.Handle("/report", report)

	server := httptools.CreateServer(*port, h)
	server.Start()
	signal.WaitForTerminationSignal()
}

func healthHandler(health *boolMutex) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("content-type", "text/plain")
		if health.Get() {
			rw.WriteHeader(http.StatusOK)
			_, _ = rw.Write([]byte("OK"))
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
			_, _ = rw.Write([]byte("FAILURE"))
		}
	})
}

func healthInverseHandler(health *boolMutex) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		health.Inverse()
	})
}

func handleDefaultGet(rw http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(10)*time.Second)
	defer cancel()
	performDbRequest(ctx, rw, r, key)
}

func performDbRequest(ctx context.Context, rw http.ResponseWriter, r *http.Request, key string) {
	fwdRequest := r.Clone(ctx)
	fwdRequest.RequestURI = ""
	fwdRequest.URL.Host = *dbUrl
	fwdRequest.Host = *dbUrl
	fwdRequest.URL.Scheme = "http"
	fwdRequest.URL.Path = "/db/" + key

	resp, err := http.DefaultClient.Do(fwdRequest)
	if *delay > 0 && *delay < 300 {
		time.Sleep(time.Duration(*delay) * time.Millisecond)
	}
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if resp.StatusCode == http.StatusBadRequest {
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil && string(body) == "record does not exist\n" {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
	}
	
	report.Process(r)
	copyResponseDetails(rw, resp)
}

func copyResponseDetails(rw http.ResponseWriter, resp *http.Response) {
	rw.WriteHeader(resp.StatusCode)
	rw.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	rw.Header().Set("Content-Length", resp.Header.Get("Content-Length"))
	io.Copy(rw, resp.Body)
	resp.Body.Close()
}

func createTeam() {
	formData := url.Values{}
	formData.Set("value", time.Now().Format("2006-01-02"))

	resp, err := http.PostForm(scheme + "://" + *dbUrl + "/db/" + teamName, formData)
	if err != nil || resp.StatusCode != http.StatusOK {
		panic("Error occured when initializing DB")
	}
}