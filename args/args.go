package args

type StringArg struct {
	name  string
	usage string
	path  *string
}

func NewStringArg(name, usage string, path *string) *StringArg {
	return &StringArg{name, usage, path}
}

func (a *StringArg) Type() string {
	return a.name
}

func (a *StringArg) Usage() string {
	return a.usage
}

func (a *StringArg) Set(v string) error {
	*a.path = v
	return nil
}

func (a *StringArg) Autocomplete(input string) []string {
	return nil
}
