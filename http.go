package aweFunc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func RequestPostJSON[T any](ctx context.Context, url string, body any, out *T) error {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	var rdr io.Reader = http.NoBody
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("RequestPostJSON: marshal body: %w", err)
		}
		rdr = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, rdr)
	if err != nil {
		return fmt.Errorf("RequestPostJSON: make request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("RequestPostJSON: do request: %w", err)
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		slurp, _ := io.ReadAll(io.LimitReader(res.Body, 4096))
		_, _ = io.Copy(io.Discard, res.Body)
		return fmt.Errorf("RequestPostJSON: status %d: %s", res.StatusCode, strings.TrimSpace(string(slurp)))
	}

	if res.StatusCode == http.StatusNoContent || out == nil {
		_, _ = io.Copy(io.Discard, res.Body)
		return nil
	}

	if err := json.NewDecoder(res.Body).Decode(out); err != nil {
		return fmt.Errorf("RequestPostJSON: decode: %w", err)
	}

	return nil
}
