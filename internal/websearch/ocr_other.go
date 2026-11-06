//go:build !windows

package websearch

import (
	"context"
	"errors"
)

func screenText(ctx context.Context, data, image string) (string, error) {
	return "", errors.New("Windows screenshot text recognition is required")
}
