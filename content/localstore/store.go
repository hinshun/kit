package localstore

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hinshun/kit"
	"github.com/hinshun/kit/content"
)

type store struct {
	content.Store
}

func NewStore(s content.Store) content.Store {
	return &store{
		Store: s,
	}
}

func (s *store) Get(ctx context.Context, digest string) (string, error) {
	if len(digest) < 3 {
		return "", fmt.Errorf("digest '%s' must be at least 3 characters", digest)
	}

	dir := filepath.Join(os.Getenv("HOME"), kit.KitDir, "store", NextToLast(digest))
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return "", err
	}

	filename := filepath.Join(dir, digest)
	_, err = os.Stat(filename)
	if err != nil {
		if !os.IsNotExist(err) {
			return filename, err
		}

		return s.Store.Get(ctx, digest)
	}

	return filename, nil
}

func NextToLast(ref string) string {
	nextToLastLen := 2
	offset := len(ref) - nextToLastLen - 1
	return ref[offset : offset+nextToLastLen]
}
