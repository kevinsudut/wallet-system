package handlertemplate

import "fmt"

type ErrReader struct{}

func (e ErrReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("foo")
}
