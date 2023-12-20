package with

import (
	"errors"
	"io"
)

type closerErr struct {
	closeErr error
	innerErr error
}

func (err *closerErr) Error() string {
	return errors.Join(err.closeErr, err.innerErr).Error()
}

func (err *closerErr) Unwrap() []error {
	return []error{
		err.closeErr,
		err.innerErr,
	}
}

func Closer(closer io.Closer, fn func() error) (err error) {
	defer func() {
		r := closer.Close()

		if r != nil {
			err = &closerErr{
				closeErr: r,
				innerErr: err,
			}
		}
	}()

	return fn()
}
