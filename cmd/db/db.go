package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Dimdim28/lab4-software-architecture/datastore"
	"github.com/Dimdim28/lab4-software-architecture/httptools"
	"github.com/Dimdim28/lab4-software-architecture/signal"
)

var (
	port = flag.Int("port", 8100, "server port")
	db *datastore.Db
)

func main() {
	flag.Parse()
	var err error
	db, err = datastore.NewDb("./out")
	if err != nil {
		panic(err)
	}
	startServer()
	signal.WaitForTerminationSignal()
}

func startServer() {
	handler := http.NewServeMux()
	handler.HandleFunc("/db/", handleDb)
	server := httptools.CreateServer(*port, handler)
	server.Start()
}

func handleDb(rw http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/db/")
	switch r.Method {
	case http.MethodGet:
		t := r.URL.Query().Get("type")
		switch t {
		case "", "string":
			sendResponse(rw, getString(key))
		case "int64":
			sendResponse(rw, getInt64(key))
		default:
			http.Error(rw, "Wrong type of data", http.StatusBadRequest)
		}
	case http.MethodPost:
		t := r.URL.Query().Get("type")
		value := r.FormValue("value")
		switch t {
		case "", "string":
			sendResponse(rw, putString(key, value))
		case "int64":
			sendResponse(rw, putInt64(key, value))
		default:
			http.Error(rw, "Wrong type of data", http.StatusBadRequest)
		}
	default:
		http.Error(rw, "Wrong method", http.StatusMethodNotAllowed)
	}
}

func sendResponse(rw http.ResponseWriter, data interface{}, err error) {
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	} else if err := json.NewEncoder(rw).Encode(data); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func getString(key string) (interface{}, error) {
	value, err := db.Get(key)
	if err != nil {
		return nil, err
	}
	return struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}{key, value}, nil
}

func getInt64(key string) (interface{}, error) {
	value, err := db.GetInt64(key)
	if err != nil {
		return nil, err
	}
	return struct {
		Key   string `json:"key"`
		Value int64  `json:"value"`
	}{key, value}, nil
}

func putString(key, value string) error {
	if value == "" {
		return fmt.Errorf("The value is empty")
	}
	return db.Put(key, value)
}

func putInt64(key, value string) error {
	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fmt.Errorf("This data type can not be converted")
	}
	return db.PutInt64(key, i)
}
