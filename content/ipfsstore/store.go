package ipfsstore

import (
	"context"
	"os"
	"path/filepath"

	"github.com/hinshun/kit"
	"github.com/hinshun/kit/content"
	"github.com/hinshun/kit/content/localstore"
	shell "github.com/ipfs/go-ipfs-api"
)

type store struct {
}

func NewStore() content.Store {
	return &store{}
}

func (s *store) Get(ctx context.Context, digest string) (string, error) {
	dir := filepath.Join(os.Getenv("HOME"), kit.KitDir, "store", localstore.NextToLast(digest))

	// resp, err := http.Get(fmt.Sprintf("https://ipfs.io/ipfs/%s", digest))
	// if err != nil {
	// 	return filename, err
	// }
	// defer resp.Body.Close()

	// data, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return filename, err
	// }

	// err = ioutil.WriteFile(filename, data, 0664)
	// if err != nil {
	// 	return filename, err
	// }

	sh := shell.NewLocalShell()

	filename := filepath.Join(dir, digest)
	err := sh.Get(digest, filename)
	if err != nil {
		return filename, err
	}

	return filename, nil
}
