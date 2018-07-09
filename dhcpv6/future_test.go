package dhcpv6

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestResponseValue(t *testing.T) {
	m, err := NewMessage()
	require.NoError(t, err)
	r := NewResponse(m, nil)
	require.Equal(t, r.Value(), m)
	require.Equal(t, r.Error(), nil)
}

func TestResponseError(t *testing.T) {
	e := errors.New("Test error")
	r := NewResponse(nil, e)
	require.Equal(t, r.Value(), nil)
	require.Equal(t, r.Error(), e)
}

func TestSuccessFuture(t *testing.T) {
	m, err := NewMessage()
	require.NoError(t, err)
	f := NewSuccessFuture(m)

	val, err := f.Wait()
	require.NoError(t, err)
	require.Equal(t, val, m)
}

func TestFailureFuture(t *testing.T) {
	e := errors.New("Test error")
	f := NewFailureFuture(e)

	val, err := f.Wait()
	require.Equal(t, err, e)
	require.Equal(t, val, nil)
}

func TestThenSuccess(t *testing.T) {
	m, err := NewMessage()
	require.NoError(t, err)
	s, err := NewMessage()
	require.NoError(t, err)
	e := errors.New("Test error")

	f := NewSuccessFuture(m).
		Then(func(_ DHCPv6) Future {
			return NewSuccessFuture(s)
		}, func(_ error) Future {
			return NewFailureFuture(e)
		})

	val, err := f.Wait()
	require.NoError(t, err)
	require.NotEqual(t, val, m)
	require.Equal(t, val, s)
}

func TestThenFailure(t *testing.T) {
	m, err := NewMessage()
	require.NoError(t, err)
	s, err := NewMessage()
	require.NoError(t, err)
	e := errors.New("Test error")
	e2 := errors.New("Test error 2")

	f := NewFailureFuture(e).
		Then(func(_ DHCPv6) Future {
			return NewSuccessFuture(s)
		}, func(_ error) Future {
			return NewFailureFuture(e2)
		})

	val, err := f.Wait()
	require.Error(t, err)
	require.NotEqual(t, val, m)
	require.NotEqual(t, val, s)
	require.NotEqual(t, err, e)
	require.Equal(t, err, e2)
}

func TestOnSuccess(t *testing.T) {
	m, err := NewMessage()
	require.NoError(t, err)
	s, err := NewMessage()
	require.NoError(t, err)

	f := NewSuccessFuture(m).
		OnSuccess(func(_ DHCPv6) Future {
			return NewSuccessFuture(s)
		})

	val, err := f.Wait()
	require.NoError(t, err)
	require.NotEqual(t, val, m)
	require.Equal(t, val, s)
}

func TestOnSuccessForFailureFuture(t *testing.T) {
	m, err := NewMessage()
	require.NoError(t, err)
	e := errors.New("Test error")

	f := NewFailureFuture(e).
		OnSuccess(func(_ DHCPv6) Future {
			return NewSuccessFuture(m)
		})

	val, err := f.Wait()
	require.Error(t, err)
	require.Equal(t, err, e)
	require.NotEqual(t, val, m)
}

func TestOnFailure(t *testing.T) {
	m, err := NewMessage()
	require.NoError(t, err)
	s, err := NewMessage()
	require.NoError(t, err)
	e := errors.New("Test error")

	f := NewFailureFuture(e).
		OnFailure(func(_ error) Future {
			return NewSuccessFuture(s)
		})

	val, err := f.Wait()
	require.NoError(t, err)
	require.NotEqual(t, val, m)
	require.Equal(t, val, s)
}

func TestOnFailureForSuccessFuture(t *testing.T) {
	m, err := NewMessage()
	require.NoError(t, err)
	s, err := NewMessage()
	require.NoError(t, err)

	f := NewSuccessFuture(m).
		OnFailure(func(_ error) Future {
			return NewSuccessFuture(s)
		})

	val, err := f.Wait()
	require.NoError(t, err)
	require.NotEqual(t, val, s)
	require.Equal(t, val, m)
}

func TestWaitTimeout(t *testing.T) {
	m, err := NewMessage()
	require.NoError(t, err)
	s, err := NewMessage()
	require.NoError(t, err)
	f := NewSuccessFuture(m).OnSuccess(func(_ DHCPv6) Future {
		time.Sleep(1 * time.Second)
		return NewSuccessFuture(s)
	})
	val, err := f.WaitTimeout(50 * time.Millisecond)
	require.Error(t, err)
	require.Equal(t, err.Error(), "Timed out")
	require.NotEqual(t, val, m)
	require.NotEqual(t, val, s)
}
