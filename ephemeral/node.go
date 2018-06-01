package ephemeral

import (
	"context"
	"fmt"
	"gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/core"
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
			return nil, fmt.Errorf("failed to init ephemeral node: %s", err)
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
