package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	port        int
	config_file string
)

func init() {
	flag.IntVar(&port, "port", 18888, "http service port")
	flag.StringVar(&config_file, "config", "config.json", "config file")
}

func main() {
	var addr string = fmt.Sprintf("0.0.0.0:%d", port)
	start_http_service_on(addr)
}

func start_http_service_on(addr string) {
	log.Printf("[network] ======== serve start, at: %s", addr)
	var svr_mux = http.NewServeMux()
	svr_mux.HandleFunc("/", http_handler)

	server := &http.Server{Addr: addr, Handler: svr_mux}

	if err := server.ListenAndServe(); err == http.ErrServerClosed {
		log.Println("[network] serve closed")
	} else {
		log.Println("[network] ======== serve failed! ========", err)
	}
}

func http_handler(w http.ResponseWriter, r *http.Request) {
	for true {
		if r.Method != http.MethodGet {
			break
		}

		var query = r.URL.Query()
		if query == nil {
			break
		}

		var flag = query["go-get"]
		if flag == nil || len(flag) == 0 || flag[0] != "1" {
			break
		}

		var host = r.Host
		var path = r.URL.Path
		// TODO: get from config
		var repo_type = ""
		var repo_url = ""
		var resp = fmt.Sprintf(html_template, host, path, repo_type, repo_url)
		w.Write([]byte(resp))
		return
	}

	w.WriteHeader(403)
}

const (
	html_template = `
	<html>
		<head>
			<meta name="go-import" content="%s%s %s %s">
		</head>
		<body></body>
	<html>
	`
)
