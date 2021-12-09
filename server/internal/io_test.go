package internal

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"go.uber.org/multierr"
)

func TestSend_success(t *testing.T) {
	ctx := context.Background()
	content := "test"
	send := func(b []byte) error {
		if string(b) != content {
			t.Errorf("Writer got %q; want %q", string(b), content)
		}
		return nil
	}

	errCh1 := make(chan error, 1)
	errCh2 := make(chan error, 1)
	go sendBuffer(ctx, strings.NewReader(content), errCh1, send)
	go sendBuffer(ctx, strings.NewReader(content), errCh2, send)
	if err := waitForSenders(ctx, errCh1, errCh2); err != nil {
		t.Error(err)
	}
}

func TestSend_writer_error(t *testing.T) {
	ctx := context.Background()
	send := func([]byte) error {
		return errors.New("writer error")
	}

	errCh1 := make(chan error, 1)
	errCh2 := make(chan error, 1)
	go sendBuffer(ctx, strings.NewReader("test"), errCh1, send)
	go sendBuffer(ctx, strings.NewReader("test"), errCh2, send)
	if err := waitForSenders(ctx, errCh1, errCh2); len(multierr.Errors(err)) != 2 {
		t.Errorf("Got error: %v; want 2 errors", err)
	}
}

func TestSend_context_done(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	send := func([]byte) error {
		// We should not block or get an error.
		time.Sleep(time.Minute * 10)
		return errors.New("error")
	}

	errCh1 := make(chan error, 1)
	errCh2 := make(chan error, 1)
	go sendBuffer(ctx, strings.NewReader("test"), errCh1, send)
	go sendBuffer(ctx, strings.NewReader("test"), errCh2, send)
	cancel()
	if err := waitForSenders(ctx, errCh1, errCh2); err != nil {
		t.Error(err)
	}
}
