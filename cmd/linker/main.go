package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hinshun/kit/config"
	"github.com/hinshun/kit/publish"
	shell "github.com/ipfs/go-ipfs-api"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("linker: requires exactly 1 arg for ipfs multiaddr")
	}

	if err := run(os.Args[1]); err != nil {
		fmt.Fprintf(os.Stderr, "linker: %s\n", err)
		os.Exit(1)
	}
}

func run(host string) error {
	pathsByPlugin := make(map[string][]string)

	err := filepath.Walk("./bin", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		filename := filepath.Base(path)
		parts := strings.Split(filename, "-")
		if len(parts) != 3 {
			return fmt.Errorf("expected plugin path '%s' to be name-GOOS-GOARCH", filename)
		}

		plugin := filepath.Join(filepath.Dir(path), parts[0])
		if plugin == "bin/kit" {
			return nil
		}

		pathsByPlugin[plugin] = append(pathsByPlugin[plugin], path)
		return nil
	})
	if err != nil {
		return err
	}

	var names []string
	for plugin := range pathsByPlugin {
		names = append(names, plugin)
	}
	sort.Strings(names)

	sh := shell.NewShell(host)

	pluginNamespace := config.Plugin{
		Name:     "plugin",
		Manifest: "/kit/plugin",
	}

	content, err := json.MarshalIndent(pluginNamespace, "", "    ")
	if err != nil {
		return err
	}

	digest, err := sh.Add(bytes.NewReader(content))
	if err != nil {
		return err
	}

	ldflags := []string{
		fmt.Sprintf(
			"-X github.com/hinshun/kit/content/kitstore.PluginDigest=%s",
			digest,
		),
	}

	for _, name := range names {
		paths, ok := pathsByPlugin[name]
		if !ok {
			return fmt.Errorf("did not find plugin '%s'", name)
		}

		digest, err := publish.Publish(sh, paths)
		if err != nil {
			return err
		}

		var varName []string
		for _, part := range strings.Split(name, "/")[1:] {
			varName = append(varName, strings.Title(part))
		}

		ldflags = append(ldflags,
			fmt.Sprintf(
				"-X github.com/hinshun/kit/content/kitstore.%sDigest=%s",
				strings.Join(varName, ""),
				digest,
			),
		)
	}

	fmt.Printf("%s", strings.Join(ldflags, " "))
	return nil
}
