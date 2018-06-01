package ipfs

import (
	"context"
	"fmt"
	offline "gx/ipfs/QmWM5HhdG5ZQNyHQ5XhMdGmV9CvLpFynQfGpTxN2MEM7Lc/go-ipfs-exchange-offline"
	datastore "gx/ipfs/QmXRKBQA4wXP7xWbFiZsR1GP4HV6wMDQ1aWFxZZ4uBcPX9/go-datastore"
	blockstore "gx/ipfs/QmaG4DZ4JaqEfvPWt5nPPgoTzhc1tr1T3f4Nu9Jpdm8ymY/go-ipfs-blockstore"
	"gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/blockservice"
	"gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/core"
	"gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/merkledag"
	"gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/pin"
	"gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/repo/config"
	"gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/repo/fsrepo"

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

func NewInmemoryNode(ctx context.Context) *core.IpfsNode {
	n := &core.IpfsNode{}

	ds := datastore.NewMapDatastore()
	bs := blockstore.NewBlockstore(ds)
	gcl := blockstore.NewGCLocker()
	bserv := blockservice.New(bs, offline.Exchange(bs))
	dserv := merkledag.NewDAGService(bserv)

	n.Blockstore = blockstore.NewGCBlockstore(bs, gcl)
	n.Pinning = pin.NewPinner(ds, dserv, dserv)
	n.DAG = dserv
	return n
}
