package resty

import (
	"bytes"
	"io"
	"net/http"
)

// ... existing code ...

func (r *Request) prepareBody(req *http.Request, body interface{}) error {
	switch v := body.(type) {
	case io.Reader:
		// If it's a ReadSeeker, we can seek back to start
		if seeker, ok := v.(io.ReadSeeker); ok {
			req.GetBody = func() (io.ReadCloser, error) {
				_, err := seeker.Seek(0, io.SeekStart)
				if err != nil {
					return nil, err
				}
				return io.NopCloser(seeker), nil
			}
		} else {
			// Buffer the reader so it can be re-read
			buf, err := io.ReadAll(v)
			if err != nil {
				return err
			}
			req.Body = io.NopCloser(bytes.NewReader(buf))
			req.GetBody = func() (io.ReadCloser, error) {
				return io.NopCloser(bytes.NewReader(buf)), nil
			}
			return nil
		}
	}
	// ... existing logic for other types ...
	return nil
}