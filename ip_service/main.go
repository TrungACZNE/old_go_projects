package main

import (
	"fmt"
	"net/http"
	"strings"
)

func extractIP(addr string) string {
	if p := strings.Index(addr, ":"); p >= 0 {
		return addr[:p]
	} else {
		return addr
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	if val, ok := r.Header["X-Forwarded-For"]; ok && len(val) > 0 && len(val[0]) > 0 && val[0] != "unknown" {
		fmt.Fprintf(w, val[0])
	} else if val, ok := r.Header["X-Real-Ip"]; ok && len(val) > 0 && len(val[0]) > 0 && val[0] != "unknown" {
		fmt.Fprintf(w, val[0])
	} else {
		fmt.Fprintf(w, extractIP(r.RemoteAddr))
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe("127.0.0.1:20101", nil)
}
