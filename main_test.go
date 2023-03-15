package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestResult(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return run(ctx)
	})
	in := "massage"
	rsp, err := http.Get("http://localhost:8080/" + in)
	if err != nil {
		t.Errorf("failed to get: %+v", err)
	}
	defer rsp.Body.Close()
	got, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("failed to read: %+v", err)
	}

	// HTTPサーバの戻り値をチェック
	want := fmt.Sprintf("Hello %s", in)
	fmt.Println(want)
	fmt.Println(string(got))

	if string(got) != want {
		t.Errorf("want %s, but %s", want, got)
	}

	cancel()
	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
}
