package config

import (
	"log"
	"sort"
	"strings"
)

type Plugin struct {
	Name     string  `json:"name"`
	Usage    string  `json:"usage,omitempty"`
	Manifest string  `json:"manifest,omitempty"`
	Plugins  Plugins `json:"plugins,omitempty"`
}

type Plugins []Plugin

func (p Plugins) Merge(plugins Plugins) Plugins {
	indexByName := make(map[string]int)
	for i, plugin := range p {
		indexByName[plugin.Name] = i
	}

	for _, plugin := range plugins {
		i, ok := indexByName[plugin.Name]
		if !ok {
			p = append(p, plugin)
			continue
		}

		plugin.Plugins = p[i].Plugins.Merge(plugin.Plugins)
		p[i] = plugin
	}

	return p
}

// Sort lexicographically sort plugins by name to produce a deterministic
// config.
func (p Plugins) Sort() {
	sort.SliceStable(p, func(i, j int) bool {
		return p[i].Name < p[j].Name
	})
}

func (p Plugins) Print() {
	p.printTree(0)
}

func (p Plugins) printTree(depth int) {
	for _, plugin := range p {
		var spaces []string
		for i := 0; i < depth; i++ {
			spaces = append(spaces, "\t")
		}

		log.Printf("%s%s\n", strings.Join(spaces, ""), plugin.Name)
		plugin.Plugins.printTree(depth + 1)
	}
}
