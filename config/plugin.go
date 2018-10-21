package config

type Plugin struct {
	Name     string  `json:"name"`
	Usage    string  `json:"usage,omitempty"`
	Manifest string  `json:"manifest,omitempty"`
	Plugins  Plugins `json:"plugins,omitempty"`
}

type Plugins []Plugin

func (p Plugins) Walk(names []string, fn func(plugin Plugin, depth int) error) (Plugin, error) {
	current := Plugin{Plugins: p}
	for depth, name := range names {
		found := false
		for _, plugin := range current.Plugins {
			if plugin.Name == name {
				found = true
				current = plugin
				break
			}
		}

		if !found {
			break
		}

		err := fn(current, depth)
		if err != nil {
			return Plugin{}, err
		}
	}

	return current, nil
}
