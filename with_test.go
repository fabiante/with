package with_test

import (
	"errors"
	"testing"

	"github.com/fabiante/with"
)

type Closer struct {
	err    error
	closed bool
}

func NewFailCloser(err error) *Closer {
	c := new(Closer)
	c.err = err
	return c
}

func (c *Closer) Close() error {
	c.closed = true
	return c.err
}

func Test(t *testing.T) {
	t.Run("closes if no error is returned", func(t *testing.T) {
		c := new(Closer)
		err := with.Closer(c, func() error {
			return nil
		})

		requireNoError(t, err)
		requireTrue(t, c.closed)
	})

	t.Run("closes if error is returned and returns err", func(t *testing.T) {
		innerErr := errors.New("some error")

		c := new(Closer)
		err := with.Closer(c, func() error {
			return innerErr
		})

		requireErrorIs(t, err, innerErr)
		requireTrue(t, c.closed)
	})

	t.Run("returns close err close fails", func(t *testing.T) {
		closeErr := errors.New("error from closing")

		c := NewFailCloser(closeErr)
		err := with.Closer(c, func() error {
			return nil
		})

		requireErrorIs(t, err, closeErr)
		requireTrue(t, c.closed)
	})

	t.Run("returns both original err and close error if close fails", func(t *testing.T) {
		innerErr := errors.New("some error")
		closeErr := errors.New("error from closing")

		c := NewFailCloser(closeErr)
		err := with.Closer(c, func() error {
			return innerErr
		})

		requireErrorIs(t, err, innerErr)
		requireErrorIs(t, err, closeErr)
		requireTrue(t, c.closed)
	})
}

func requireNoError(t *testing.T, e error) {
	t.Helper()
	if e != nil {
		t.Fatalf("unexpected error: %v", e)
	}
}

func requireTrue(t *testing.T, b bool) {
	t.Helper()
	if !b {
		t.Fatal("expected false to be true")
	}
}

func requireErrorIs(t *testing.T, actual, expected error) {
	t.Helper()
	if !errors.Is(actual, expected) {
		t.Fatalf("expected error %v to be %v", actual, expected)
	}
}
