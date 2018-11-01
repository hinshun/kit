package local

import (
	"context"
	"os"
	"path/filepath"

	"github.com/hinshun/kit"
	"github.com/hinshun/kit/content"
	shell "github.com/ipfs/go-ipfs-api"
)

type store struct {
}

func NewStore() content.Store {
	return &store{}
}

func (s *store) Get(ctx context.Context, digest string) (string, error) {
	switch digest {
	case "/kit/init":
		digest = "QmRLyNYEx9Em1SnrN6n1yAL53Qa3eJXbVsAq4uhNBMXp1z"
	case "/kit/bootstrap":
		digest = "QmWhLFJfpdE6n8zgJVxNU8tfGjjKPdNTu1NTv3yg24sKNQ"
	case "/kit/plugin":
		digest = "QmetM7PMkuGJtwBS5Lw57cCfMDSYk4UHZRPGNzq9JtGLCP"
	case "/kit/plugin/add":
		digest = "QmXs1eXBt5BTW331cEL3d9Gd5sto3a5tnE17qZYjSVvdfi"
	case "/kit/plugin/rm":
		digest = "QmUPqj4NGi39jVjGP3YYWBQBQmNJX99g1mK3YCgZJaf5W3"
	case "/kit/plugin/publish":
		digest = "QmVSyD7JfKLaCKWZZqorG5wyUjUsZ6fVf3uRQngyJtd8BJ"
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

		sh := shell.NewLocalShell()
		err = sh.Get(digest, filename)
		if err != nil {
			return filename, err
		}
	}

	return filename, nil
}

// func SyncCommands(ctx context.Context, sh *shell.Shell, cfg *kit.Config) (refs []string, err error) {
// 	var pluginLock kit.ConfigLock

// 	lockPath := ".kit/store/plugin.lock"
// 	_, err = os.Stat(lockPath)
// 	if err != nil {
// 		if !os.IsNotExist(err) {
// 			return nil, err
// 		}
// 	} else {
// 		data, err := ioutil.ReadFile(lockPath)
// 		if err != nil {
// 			return nil, err
// 		}

// 		err = json.Unmarshal(data, &pluginLock)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	refLockByRef := make(map[string]string)
// 	for _, refLock := range pluginLock.RefLocks {
// 		refLockByRef[refLock.Ref] = refLock.Cid
// 	}

// 	var refLocks []kit.RefLock
// 	for _, ref := range cfg.Plugins.Refs() {
// 		if ref == "" {
// 			continue
// 		}

// 		c, err := SyncCommand(ctx, sh, refLockByRef, ref)
// 		if err != nil {
// 			return nil, err
// 		}

// 		refs = append(refs, ref)
// 		refLocks = append(refLocks, kit.RefLock{
// 			Ref: ref,
// 			Cid: c.String(),
// 		})
// 	}

// 	refIndex := 0
// 	for i, refLock := range pluginLock.RefLocks {
// 		found := false
// 		for j := refIndex; j < len(refs); j, refIndex = j+1, refIndex+1 {
// 			if refLock.Ref == refs[j] {
// 				found = true
// 				break
// 			}
// 		}

// 		if found {
// 			pluginLock.RefLocks = append(pluginLock.RefLocks[:i], pluginLock.RefLocks[i+1:]...)
// 		}
// 	}

// 	pluginLock.RefLocks = append(pluginLock.RefLocks, refLocks...)
// 	sort.SliceStable(pluginLock.RefLocks, func(i, j int) bool {
// 		return pluginLock.RefLocks[i].Ref < pluginLock.RefLocks[j].Ref
// 	})

// 	data, err := json.MarshalIndent(&pluginLock, "", "    ")
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = ioutil.WriteFile(lockPath, data, 0664)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return refs, nil
// }

// func SyncCommand(ctx context.Context, sh *shell.Shell, refLockByRef map[string]string, ref string) (cid.Cid, error) {
// 	dir := filepath.Join(".kit/store", NextToLast(ref))
// 	err := os.MkdirAll(dir, 0755)
// 	if err != nil {
// 		return cid.Cid{}, fmt.Errorf("failed to make kit store directory: %s", err)
// 	}

// 	filename := filepath.Join(dir, ref)
// 	expected, ok := refLockByRef[ref]
// 	if ok {
// 		_, err = os.Stat(filename)
// 		if err != nil {
// 			if !os.IsNotExist(err) {
// 				return cid.Cid{}, err
// 			}
// 			ok = false
// 		}
// 	}

// 	if !ok {
// 		err = sh.Get(ref, filename)
// 		if err != nil {
// 			return cid.Cid{}, fmt.Errorf("failed to get '%s' from ipfs: %s", ref, err)
// 		}
// 	}

// 	data, err := ioutil.ReadFile(filename)
// 	if err != nil {
// 		return cid.Cid{}, err
// 	}

// 	c, err := cid.Parse(util.Hash(data))
// 	if err != nil {
// 		return cid.Cid{}, err
// 	}

// 	if ok && expected != c.String() {
// 		fmt.Printf("plugin ref '%s' has cid '%s' mismatched with lock file cid '%s'\n", ref, c.String(), expected)

// 		err = os.Remove(filename)
// 		if err != nil {
// 			return cid.Cid{}, fmt.Errorf("failed to remove mismatched cid file '%s': %s", filename, err)
// 		}

// 		return SyncCommand(ctx, sh, refLockByRef, ref)
// 	}

// 	return c, nil
// }

func NextToLast(ref string) string {
	nextToLastLen := 2
	offset := len(ref) - nextToLastLen - 1
	return ref[offset : offset+nextToLastLen]
}
