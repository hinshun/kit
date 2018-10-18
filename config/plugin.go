package config

type Plugin struct {
	Name    string  `json:"name"`
	Ref     string  `json:"ref"`
	Plugins Plugins `json:"plugins,omitempty"`
}

type Plugins []Plugin

func (p Plugins) Walk(names []string, fn func(Plugin) error) ([]string, error) {
	var (
		current = p
		i       int
		name    string
	)
	for i, name = range names {
		var found *Plugin
		for _, plugin := range current {
			if plugin.Name == name {
				found = &plugin
				break
			}
		}

		if found == nil {
			break
		}

		err := fn(*found)
		if err != nil {
			return nil, err
		}

		current = found.Plugins
	}

	for _, plugin := range current {
		err := fn(plugin)
		if err != nil {
			return nil, err
		}
	}

	return names[i:], nil
}
