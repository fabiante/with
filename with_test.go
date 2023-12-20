package with_test

import (
	"errors"
	"testing"

	"github.com/fabiante/with"
	"github.com/stretchr/testify/require"
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

		require.NoError(t, err)
		require.True(t, c.closed)
	})

	t.Run("closes if error is returned and returns err", func(t *testing.T) {
		innerErr := errors.New("some error")

		c := new(Closer)
		err := with.Closer(c, func() error {
			return innerErr
		})

		require.ErrorIs(t, err, innerErr)
		require.True(t, c.closed)
	})

	t.Run("returns close err close fails", func(t *testing.T) {
		closeErr := errors.New("error from closing")

		c := NewFailCloser(closeErr)
		err := with.Closer(c, func() error {
			return nil
		})

		require.ErrorIs(t, err, closeErr)
		require.True(t, c.closed)
	})

	t.Run("returns both original err and close error if close fails", func(t *testing.T) {
		innerErr := errors.New("some error")
		closeErr := errors.New("error from closing")

		c := NewFailCloser(closeErr)
		err := with.Closer(c, func() error {
			return innerErr
		})

		require.ErrorIs(t, err, innerErr)
		require.ErrorIs(t, err, closeErr)
		require.True(t, c.closed)
	})
}
