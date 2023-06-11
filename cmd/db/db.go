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
	db   *datastore.Db
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
			data, err := getString(key)
			sendResponse(rw, data, err)
		case "int64":
			data, err := getInt64(key)
			sendResponse(rw, data, err)
		default:
			http.Error(rw, "ERROR! Unknown data type", http.StatusBadRequest)
		}
	case http.MethodPost:
		t := r.URL.Query().Get("type")
		value := r.FormValue("value")
		switch t {
		case "", "string":
			err := putString(key, value)
			sendResponse(rw, nil, err)
		case "int64":
			err := putInt64(key, value)
			sendResponse(rw, nil, err)
		default:
			http.Error(rw, "ERROR! Unknown data type", http.StatusBadRequest)
		}
	default:
		http.Error(rw, "ERROR! Method not allowed", http.StatusMethodNotAllowed)
	}
}

func sendResponse(rw http.ResponseWriter, data interface{}, err error) {
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	} else if data != nil {
		if err := json.NewEncoder(rw).Encode(data); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
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
		return fmt.Errorf("ERROR! Can't save empty value")
	}
	return db.Put(key, value)
}

func putInt64(key, value string) error {
	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fmt.Errorf("ERROR! Can't convert value to the given type")
	}
	return db.PutInt64(key, i)
}
