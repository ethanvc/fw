package internal

type RefWriter struct {
}

func NewRefWriter() *RefWriter {
	return &RefWriter{}
}

func (w *RefWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (w *RefWriter) Close() error {
	return nil
}
