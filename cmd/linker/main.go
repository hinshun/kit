package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hinshun/kit/publish"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "linker: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
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

	var ldflags []string
	for _, name := range names {
		paths, ok := pathsByPlugin[name]
		if !ok {
			return fmt.Errorf("did not find plugin '%s'", name)
		}

		digest, err := publish.Publish(paths)
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
