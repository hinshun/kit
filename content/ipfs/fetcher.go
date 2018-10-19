package ipfs

import (
	"github.com/hinshun/kit/content"
	shell "github.com/ipfs/go-ipfs-api"
)

type fetcher struct {
	sh *shell.Shell
}

func NewFetcher(sh *shell.Shell) content.Fetcher {
	return &fetcher{
		sh: sh,
	}
}
