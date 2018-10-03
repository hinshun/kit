package ipfs

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hinshun/kit"
	cid "github.com/ipfs/go-cid"
	util "github.com/ipfs/go-ipfs-util"
	"github.com/ipfs/go-ipfs/core"
	coreiface "github.com/ipfs/go-ipfs/core/coreapi/interface"
)

func SyncCommands(ctx context.Context, cfg *kit.Config, api coreiface.CoreAPI, hashes []string) (path []string, err error) {
	n := NewInMemoryNode()

	var paths []string
	for _, hash := range hashes {
		path, err := SyncCommand(ctx, cfg, api, n, hash)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}
	return paths, nil
}

func SyncCommand(ctx context.Context, cfg *kit.Config, api coreiface.CoreAPI, n *core.IpfsNode, key string) (filename string, err error) {
	// p, err := api.Name().Resolve(ctx, hash)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to resolve '%s': %s", hash, err)
	// }
	// fmt.Printf("resolved %s to %s\n", hash, p)

	c, err := cid.Parse(key)
	if err != nil {
		return "", err
	}
	p := coreiface.IpfsPath(c)

	filename = filepath.Join(cfg.RootDir, ".kit", p.String())
	stat, err := os.Stat(filename)
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}

	if os.IsNotExist(err) {
		fmt.Printf("downloading %s from ipfs...\n", p)
		err = WriteIPFSBlockToFile(ctx, api, p, filename)
		if err != nil {
			return "", err
		}
	} else {
		fmt.Printf("found %s, verifying hash...\n", filename)
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return "", err
		}

		expected, err := cid.Parse(util.Hash(data))
		if err != nil {
			return "", err
		}

		if expected.String() != stat.Name() {
			fmt.Printf("local hash mismatch '%s', downloading %s from ipfs...\n", expected, p)
			err = WriteIPFSBlockToFile(ctx, api, p, filename)
			if err != nil {
				return "", err
			}
		}
	}

	return filename, nil
}

func WriteIPFSBlockToFile(ctx context.Context, api coreiface.CoreAPI, p coreiface.Path, output string) error {
	reader, err := api.Block().Get(ctx, p)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(output), 0700)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(output, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
