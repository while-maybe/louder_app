package stdlibapiadapter

import (
	"fmt"
	"net/http"
	"os"
)

// const (
// 	DefaultServerAddress = ":8080"
// )

func StartServer(router http.Handler, addr string) {
	// if addr == "" {
	// 	log.Println("Server address not specified. Using port", DefaultServerAddress)
	// 	addr = DefaultServerAddress
	// }

	err := http.ListenAndServe(addr, router)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
