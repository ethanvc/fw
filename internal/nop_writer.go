package internal

type NopWriter struct {
}

func NewNopWriter() *NopWriter {
	return &NopWriter{}
}

func (w *NopWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (w *NopWriter) Close() error {
	return nil
}
