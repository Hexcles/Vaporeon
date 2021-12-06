package internal

import (
	"context"
	"io"

	"go.uber.org/multierr"
)

const streamBufferSize = 1024 * 1024 // 1 MiB

func sendBuffer(ctx context.Context, r io.Reader, errCh chan<- error, send func([]byte) error) {
	buffer := make([]byte, streamBufferSize)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			n, err := r.Read(buffer)
			// It's important to send any read bytes first
			// before checking the error as io.Reader allows
			// returning io.EOF along with the last read.
			if n > 0 {
				if err := send(buffer[:n]); err != nil {
					errCh <- err
					return
				}
			}
			if err != nil {
				errCh <- err
				return
			}
		}
	}
}

func waitForSenders(ctx context.Context, errCh1, errCh2 <-chan error) error {
	var err1, err2 error
Loop:
	for {
		select {
		case <-ctx.Done():
			return nil
		case err1 = <-errCh1:
			if err2 != nil {
				break Loop
			}
		case err2 = <-errCh2:
			if err1 != nil {
				break Loop
			}
		}
	}
	// EOF is expected.
	if err1 == io.EOF {
		err1 = nil
	}
	if err2 == io.EOF {
		err2 = nil
	}
	return multierr.Combine(err1, err2)
}
