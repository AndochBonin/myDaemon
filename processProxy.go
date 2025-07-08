package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

func handleTunneling(w http.ResponseWriter, r *http.Request) {
	dest_conn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	client_conn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	go transfer(dest_conn, client_conn)
	go transfer(client_conn, dest_conn)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

func handleHTTP(w http.ResponseWriter, req *http.Request) {
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func RunTLSProxy(certFile, keyFile string) {
	server := &http.Server{
		Addr: ":666",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Host == "localhost:666" {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				fmt.Fprintln(w, "<!DOCTYPE html><head><title>myDaemon</title></head><body><h1>hi</h1></body></html>")
			}
			if currentProcess != nil && isBlacklisted(r.URL.Hostname(), currentProcess.Program.WebHostBlocklist) {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				fmt.Fprintf(w, "%s", getBlockedURLPage(r.URL.Hostname()))
				return
			} else if r.Method == http.MethodConnect {
				handleTunneling(w, r)
			} else {
				handleHTTP(w, r)
			}
		}),
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	log.Fatal(server.ListenAndServeTLS(certFile, keyFile))
}

func isBlacklisted(targetHost string, blocklist []string) bool {
	for _, hostName := range blocklist {
		if strings.Contains(targetHost, hostName) {
			return true
		}
	}
	return false
}

func getBlockedURLPage(host string) string {
	blockedHTML := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>myDaemon</title>
		</head>
		<body>
			<h1>%s is not allowed right now</h1>
		</body>
		</html>`, host)
	return blockedHTML
}
