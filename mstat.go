package mstat

import (
	"encoding/json"
	"math"
	"net/http"
	"strings"

	"github.com/cloudfoundry/gosigar"
)

type Machine struct {
	Unit string
	Next http.Handler
}

func New() *Machine {
	return new(Machine)
}

type Uptime struct {
	Time float64 `json:"time"`
}

type Swap struct {
	Total uint64 `json:"total"`
	Free  uint64 `json:"free"`
	Used  uint64 `json:"used"`
}

type Memory struct {
	Total uint64 `json:"total"`
	Free  uint64 `json:"free"`
	Used  uint64 `json:"used"`
}

type FileSystem struct {
	Total     uint64  `json:"total"`
	Free      uint64  `json:"free"`
	Used      uint64  `json:"used"`
	Available uint64  `json:"available"`
	Files     uint64  `json:"files"`
	FreeFiles uint64  `json:"freeFiles"`
	Percent   float64 `json:"percent"`
}

func (m *Machine) Uptime() Uptime {
	uptime := sigar.Uptime{}
	if err := uptime.Get(); err != nil {
		panic(err)
	}
	return Uptime{uptime.Length}
}

func (m *Machine) Swap() Swap {
	swap := sigar.Swap{}
	if err := swap.Get(); err != nil {
		panic(err)
	}
	return Swap{
		m.format(swap.Total),
		m.format(swap.Free),
		m.format(swap.Used),
	}
}

func (m *Machine) Memory() Memory {
	mem := sigar.Mem{}
	if err := mem.Get(); err != nil {
		panic(err)
	}
	return Memory{
		m.format(mem.Total),
		m.format(mem.Free),
		m.format(mem.Used),
	}
}

func (m *Machine) FileSystem(path string) FileSystem {
	fs := sigar.FileSystemUsage{}
	if err := fs.Get(path); err != nil {
		panic(err)
	}
	return FileSystem{
		m.format(fs.Total),
		m.format(fs.Free),
		m.format(fs.Used),
		m.format(fs.Avail),
		fs.Files,
		fs.FreeFiles,
		fs.UsePercent(),
	}
}

func (m *Machine) unitFormat() (i float64) {
	switch strings.ToLower(m.Unit) {
	case "gb":
		i = 3
	case "mb":
		i = 2
	case "kb":
		i = 1
	default:
		i = 0
	}
	return i
}

func (m *Machine) format(size uint64) uint64 {
	return size / uint64(math.Pow(1024, m.unitFormat()))
}

func (m *Machine) writeJSON(w http.ResponseWriter, i interface{}) {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	w.Write(b)
}

func (m *Machine) setUnit(params map[string][]string) {
	if len(params) == 0 {
		return
	}

	for _, u := range []string{"mb", "MB", "gb", "GB", "kb", "KB"} {
		if _, ok := params[u]; ok {
			m.Unit = u
			break
		}
	}
}

func (m *Machine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		params := r.URL.Query()
		m.setUnit(params)

		switch r.URL.Path {
		case "/swap/", "/swap":
			m.writeJSON(w, m.Swap())
			return
		case "/uptime/", "/uptime":
			m.writeJSON(w, m.Uptime())
			return
		case "/memory/", "/memory", "/mem/", "/mem":
			m.writeJSON(w, m.Memory())
			return
		case "/filesystem/", "/filesystem", "/fs/", "/fs":
			path := "/"
			if p, ok := params["path"]; len(p) > 0 && ok {
				path = p[0]
			}
			m.writeJSON(w, m.FileSystem(path))
			return
		}
	}

	if m.Next != nil {
		m.Next.ServeHTTP(w, r)
	}
}
