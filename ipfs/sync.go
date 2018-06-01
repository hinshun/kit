package ipfs

import (
	"context"
	"fmt"
	"gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/core"
	coreiface "gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/core/coreapi/interface"
	"gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/core/coreunix"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

func SyncCommands(ctx context.Context, api coreiface.CoreAPI, hashes []string) (path []string, err error) {
	n := NewInmemoryNode(ctx)

	var paths []string
	for _, hash := range hashes {
		path, err := SyncCommand(ctx, api, n, hash)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}
	return paths, nil
}

func SyncCommand(ctx context.Context, api coreiface.CoreAPI, n *core.IpfsNode, hash string) (filename string, err error) {
	path, err := api.Name().Resolve(ctx, hash)
	if err != nil {
		return "", err
	}
	logrus.Debugf("resolved %s to %s", hash, path)

	filename = fmt.Sprintf("%s/.kit%s", os.Getenv("HOME"), path)
	stat, err := os.Stat(filename)
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}

	if os.IsNotExist(err) {
		logrus.Debugf("downloading %s from ipfs...", path)
		err = DownloadIPFSObject(ctx, api, path, filename)
		if err != nil {
			return "", err
		}
	} else {
		logrus.Debugf("found %s, verifying hash...", filename)
		r, err := os.Open(filename)
		if err != nil {
			return "", err
		}

		key, err := coreunix.Add(n, r)
		if err != nil {
			return "", err
		}

		if key != stat.Name() {
			logrus.Debugf("local hash mismatch '%s', downloading %s from ipfs...", key, path)
			err = DownloadIPFSObject(ctx, api, path, filename)
			if err != nil {
				return "", err
			}
		}
	}

	return filename, nil
}

func DownloadIPFSObject(ctx context.Context, api coreiface.CoreAPI, path coreiface.Path, output string) error {
	readCloser, err := api.Unixfs().Cat(ctx, path)
	if err != nil {
		return err
	}
	defer readCloser.Close()

	data, err := ioutil.ReadAll(readCloser)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(output, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
