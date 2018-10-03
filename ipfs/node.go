package ipfs

import (
	"context"
	"fmt"
	datastore "github.com/ipfs/go-datastore"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	offline "github.com/ipfs/go-ipfs-exchange-offline"
	"github.com/ipfs/go-ipfs/blockservice"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/merkledag"
	"github.com/ipfs/go-ipfs/pin"
	"github.com/ipfs/go-ipfs/repo/config"
	"github.com/ipfs/go-ipfs/repo/fsrepo"

	"os"
)

func NewNode(ctx context.Context, bootstrap []string) (*core.IpfsNode, error) {
	dir := fmt.Sprintf("%s/.kit/repo", os.Getenv("HOME"))
	if !fsrepo.IsInitialized(dir) {
		cfg, err := config.Init(os.Stdout, 2048)
		if err != nil {
			return nil, err
		}
		cfg.Bootstrap = bootstrap

		err = fsrepo.Init(dir, cfg)
		if err != nil {
			return nil, err
		}
	}

	repo, err := fsrepo.Open(dir)
	if err != nil {
		return nil, err
	}

	return core.NewNode(ctx, &core.BuildCfg{
		Repo:   repo,
		Online: true,
	})
}

func NewInMemoryNode() *core.IpfsNode {
	dstore := datastore.NewMapDatastore()
	bstore := blockstore.NewBlockstore(dstore)
	bserv := blockservice.New(bstore, offline.Exchange(bstore))
	dserv := merkledag.NewDAGService(bserv)

	return &core.IpfsNode{
		Blockstore: blockstore.NewGCBlockstore(bstore, blockstore.NewGCLocker()),
		Pinning:    pin.NewPinner(dstore, dserv, dserv),
		DAG:        dserv,
	}
}
