package misc

import "fmt"

func SafeWithError(fn func()) (retErr error) {
	defer func() {
		if err := recover(); err != nil {
			if ev, ok := err.(error); ok {
				retErr = ev
			} else {
				retErr = fmt.Errorf("%v", err)
			}
		}
	}()
	fn()
	return nil
}

func SafeValueWithError[T any](fn func() T) (retVal T, retErr error) {
	defer func() {
		if err := recover(); err != nil {
			if ev, ok := err.(error); ok {
				retErr = ev
			} else {
				retErr = fmt.Errorf("%v", err)
			}
			var zero T
			retVal = zero
		}
	}()
	return fn(), nil
}
