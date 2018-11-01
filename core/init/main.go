package main

import "github.com/hinshun/kit"

var New kit.Constructor = func() (kit.Command, error) {
	return &command{}, nil
}

func main() {}
