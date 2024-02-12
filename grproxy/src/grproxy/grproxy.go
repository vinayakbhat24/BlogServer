package main

import (
    "log"
    "net/http"
    "net/http/httputil"
    "fmt"
    "github.com/samuel/go-zookeeper/zk"
    "time"
)

var i = 0

func main() {
fmt.Printf("Inside grproxy_main")

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        director := func(req *http.Request) {
            req = r
            req.URL.Scheme = "http"
            req.URL.Host = "nginx:80"
        }
        proxy := &httputil.ReverseProxy{Director: director}
        proxy.ServeHTTP(w, r)
    })

    http.HandleFunc("/library", func(w http.ResponseWriter, r *http.Request) {

	conn, _, err1 := zk.Connect([]string{"hbase:2181"}, time.Second)

	if err1 != nil {
		fmt.Printf(" connect zk error: %s ", err1)
	} 
	defer conn.Close()


	list, _, err2 := conn.Children("/go_servers")

	if err2 != nil {
		fmt.Printf(" get server list error: %s \n", err2)
		return
	}

	count := len(list)
		fmt.Printf(" count: %d \n", count)
	if count == 0 {
		fmt.Printf(" No servers available: %s \n")
		return
	}
	i = (i+1)%count


        director := func(req *http.Request) {
            req = r
            req.URL.Scheme = "http"
            req.URL.Host = list[i]+":9090"
        }

        proxy := &httputil.ReverseProxy{Director: director}
        proxy.ServeHTTP(w, r)
    })
    log.Fatal(http.ListenAndServe(":8080", nil))
}

