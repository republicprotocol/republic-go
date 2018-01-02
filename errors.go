package dht

import "fmt"

type ErrIndexOutOfRange string

func NewErrIndexOutOfRange(i int) error {
	return ErrIndexOutOfRange(fmt.Sprintf("index %d is out of range", i))
}

func (err ErrIndexOutOfRange) Error() string {
	return string(err)
}

type ErrFullBucket string

func NewErrFullBucket() error {
	return ErrFullBucket("bucket is full")
}

func (err ErrFullBucket) Error() string {
	return string(err)
}
