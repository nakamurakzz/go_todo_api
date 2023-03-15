package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/sync/errgroup"
)

// context.Context:
// 1. キャンセルやデッドラインの伝播
// 2. リクエストまたはトランザクションスコープのメタデータを関数やゴルーチン間で伝播
func run(ctx context.Context) error {
	s := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello", r.URL.Path[1:])
		}),
	}

	// goルーチン間でエラーを伝播するためのコンテキストを作成
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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
	err := run(context.Background())

	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}
}
