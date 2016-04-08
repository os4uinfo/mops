package client

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func httpService(addr string) {
	http.HandleFunc("/stop", handleStop)

	fmt.Println("worker http service: " + addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handleStop(w http.ResponseWriter, r *http.Request) {
	if canDo(r.RemoteAddr) {
		close()
	} else {
		logError("not local addr, but stop worker!", r.RemoteAddr)
	}
}

func canDo(remoteAddr string) bool {
	return strings.HasPrefix(remoteAddr, "127.0.0.1")
}

