package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	port        int
	config_file string
	config      Config
)

func init() {
	flag.IntVar(&port, "port", 18888, "http service port")
	flag.StringVar(&config_file, "config", "config.json", "config file")
}

func main() {
	flag.Parse()

	var addr string = fmt.Sprintf("0.0.0.0:%d", port)

	if config_data, err := ioutil.ReadFile(config_file); err != nil {
		log.Fatalf("load config [%s] failed.", config_file)
	} else if err = json.Unmarshal(config_data, &config); err != nil {
		log.Fatalf("parse config file [%s] failed.", config_file)
	}

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
		// real host from reverse proxy
		if header_vals, ex := r.Header["X-Forwarded-Host"]; ex && len(header_vals) == 1 {
			host = header_vals[0]
		}
		var path = r.URL.Path
		// TODO: get from config
		var repo_type, repo_url, exist = config.Get(path)
		if !exist {
			break
		}

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

/////////////////////////////////////////////////////////////
type Config struct {
	Repositories []struct {
		Path     string `json:"path"`
		RepoType string `json:"repo_type"`
		RepoUrl  string `json:"repo_url"`
	} `json:"repositories"`
}

func (cfg *Config) Get(path string) (string, string, bool) {
	for _, repo := range cfg.Repositories {
		if repo.Path == path {
			return repo.RepoType, repo.RepoUrl, true
		}
	}
	return "", "", false
}
