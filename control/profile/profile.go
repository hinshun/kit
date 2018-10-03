package profile

import (
	"io"
	"os"
	"runtime/pprof"
)

var (
	CpuProfile = "cpu.prof"
)

type Profile struct {
	w io.WriteCloser
}

func NewProfile() (*Profile, error) {
	w, err := os.Create(CpuProfile)
	if err != nil {
		return nil, err
	}

	pprof.StartCPUProfile(w)
	return &Profile{
		w: w,
	}, nil
}

func (p *Profile) Close() {
	defer p.w.Close()
	pprof.StopCPUProfile()
}
