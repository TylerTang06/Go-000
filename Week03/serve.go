package week03

import (
	"fmt"
	"net/http"

	"golang.org/x/sync/errgroup"
)

// 基于 errgroup 实现一个 http server 的启动和关闭

func serveApp() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(resp, "Hello golang")
	})

	return http.ListenAndServe("0.0.0.0:18080", mux)
}

func serveDebug() error {
	return http.ListenAndServe("127.0.0.1:18081", nil)
}

// HandleServes ...
func HandleServes() {
	g := new(errgroup.Group)
	g.Go(serveApp)
	g.Go(serveDebug)

	if err := g.Wait(); err != nil {
		fmt.Printf("Serve meet error=%+v\n", err)
	}
}
