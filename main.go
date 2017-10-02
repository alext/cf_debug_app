package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

func main() {
	addr := ":" + os.Getenv("PORT")
	fmt.Println("Listening on", addr)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("OK\n")) })
	http.HandleFunc("/env", env)
	http.HandleFunc("/request-info", requestInfo)
	http.HandleFunc("/boom", func(w http.ResponseWriter, r *http.Request) { http.Error(w, "Error", http.StatusInternalServerError) })
	http.HandleFunc("/missing", http.NotFound)
	http.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
		w.Write([]byte("OK\n"))
	})

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func env(w http.ResponseWriter, r *http.Request) {
	env := newEnvMap(os.Environ())
	data := make(map[string]interface{})
	data["ENV"] = expandJSONFields(env)

	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(output)
}

func requestInfo(w http.ResponseWriter, r *http.Request) {
	info, err := httputil.DumpRequest(r, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(info)
}

func expandJSONFields(env envMap) map[string]interface{} {
	data := make(map[string]interface{}, len(env))
	for key, value := range env {
		if value[0] == '{' || value[0] == '[' {
			var valueData interface{}
			err := json.Unmarshal([]byte(value), &valueData)
			if err != nil {
				data[key] = value
			} else {
				data[key] = valueData
			}
		} else {
			data[key] = value
		}
	}
	return data
}
