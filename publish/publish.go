package publish

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/hinshun/kit"
	"github.com/hinshun/kit/config"
	shell "github.com/ipfs/go-ipfs-api"
)

func Publish(pluginPaths []string) error {
	sh := shell.NewLocalShell()

	manifest := config.Manifest{
		Type: config.CommandManifest,
	}

	var nativePluginPath string
	for _, pluginPath := range pluginPaths {
		parts := strings.Split(pluginPath, "-")
		if len(parts) != 3 {
			return fmt.Errorf("expected plugin path to be name-GOOS-GOARCH")
		}

		if parts[1] == runtime.GOOS && parts[2] == runtime.GOARCH {
			nativePluginPath = pluginPath
		}

		f, err := os.Open(pluginPath)
		if err != nil {
			return err
		}

		digest, err := sh.Add(f)
		if err != nil {
			return err
		}

		manifest.Platforms = append(manifest.Platforms, config.Platform{
			OS:           parts[1],
			Architecture: parts[2],
			Digest:       digest,
		})
	}

	if nativePluginPath == "" {
		return fmt.Errorf("expected one of plugin path to be native")
	}

	constructor, err := kit.OpenConstructor(nativePluginPath)
	if err != nil {
		return err
	}

	cmd, err := constructor()
	if err != nil {
		return err
	}

	manifest.Usage = cmd.Usage()

	for _, arg := range cmd.Args() {
		manifest.Args = append(manifest.Args, config.Arg{
			Type:  arg.Type(),
			Usage: arg.Usage(),
		})
	}

	for _, flag := range cmd.Flags() {
		manifest.Flags = append(manifest.Flags, config.Flag{
			Name:  flag.Name(),
			Type:  flag.Type(),
			Usage: flag.Usage(),
		})
	}

	content, err := json.MarshalIndent(&manifest, "", "    ")
	if err != nil {
		return err
	}

	digest, err := sh.Add(bytes.NewReader(content))
	if err != nil {
		return err
	}

	fmt.Println(digest)
	return nil
}
