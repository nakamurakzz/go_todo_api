package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"golang.org/x/sync/errgroup"
)

// context.Context:
// 1. キャンセルやデッドラインの伝播
// 2. リクエストまたはトランザクションスコープのメタデータを関数やゴルーチン間で伝播
func run(ctx context.Context, l net.Listener) error {
	s := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello %s", r.URL.Path[1:])
		}),
	}

	// goルーチン間でエラーを伝播するためのコンテキストを作成
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if err := s.Serve(l); err != nil && err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})

	<-ctx.Done()
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
		return err
	}
	return eg.Wait()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <port>")
		os.Exit(1)
	}
	p := os.Args[1]
	l, err := net.Listen("tcp", ":"+p)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}

	if err := run(context.Background(), l); err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}
}
