package request

import "io"

type nopReadSeekCloser struct {
	io.ReadSeeker
}

func (s nopReadSeekCloser) Close() error {
	return nil
}

func NopSeekerCloser(r io.ReadSeeker) io.ReadSeekCloser {
	return nopReadSeekCloser{r}
}
