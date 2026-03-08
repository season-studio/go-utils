package ioex

import "fmt"

func Flush(w any) (retErr error) {
	defer func() {
		if err := recover(); err != nil {
			if rerr, ok := err.(error); ok {
				retErr = rerr
			} else {
				retErr = fmt.Errorf("%v", err)
			}
		}
	}()
	if flusher, ok := w.(interface{ Flush() error }); ok {
		return flusher.Flush()
	}
	return nil
}
