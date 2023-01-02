package plugin

import (
	"github.com/superseriousbusiness/gotosocial/cmd/gotosocial/action"
	"github.com/superseriousbusiness/gotosocial/internal/gtsmodel"

	"context"
	"fmt"
	"github.com/yuin/gopher-lua"
	"layeh.com/gopher-luar"
	"os"
)

type plugin struct {
	Name string
	tbl  lua.LValue
}

type Plugins struct {
	state   lua.LState
	plugins []plugin
}

func initialize(l lua.LState, name string) plugin {
	p := plugin{Name: name}
	l.Push(l.GetGlobal("require"))
	l.Push(luar.New(&l, name))
	l.Call(1, 1)
	p.tbl = l.Get(1)
	l.Pop(1)

	init := l.GetField(p.tbl, "init")
	if lua.LVAsBool(init) {
		l.Push(init)
		l.Call(0, 0)
	}
	return p
}

func Process(ps Plugins, s *gtsmodel.Status) {
	for _, p := range ps.plugins {
		process := ps.state.GetField(p.tbl, "process")
		if lua.LVAsBool(process) {
			ps.state.Push(process)
			ps.state.Push(luar.New(&ps.state, s))
			ps.state.Call(1, 0)
		}
	}
}

var Run action.GTSAction = func(ctx context.Context) error {
	ps := Init()
	var args []string = ctx.Value("args").([]string)
	name := args[0]
	for _, p := range ps.plugins {
		if p.Name == name {
			run := ps.state.GetField(p.tbl, "run")
			// does this plugin have a run function?
			if lua.LVAsBool(run) {
				ps.state.Push(run)
				ps.state.Push(luar.New(&ps.state, args[1:]))
				ps.state.Call(1, 0)
				return nil
			} else {
				return fmt.Errorf("plugin does not have run function: %s", name)
			}
		}
	}
	return fmt.Errorf("plugin not found: %s", name)
}

// eventually we're going to want an environment exposing some callbacks:
// * follow
// * unfollow
// * lookup_user
// * lookup_post
// * post

func Init() Plugins {
	plugins := Plugins{state: *lua.NewState()}
	plugins.state.DoString("_G.io, _G.os = nil") // plugins can't touch host
	plugins.state.DoString("package.path = 'plugins/?.lua;plugins/?/init.lua'")
	files, err := os.ReadDir("plugins/")

	// if we don't have a plugins dir, don't sweat it!
	if err == nil {
		for _, file := range files {
			if file.IsDir() {
				p := initialize(plugins.state, file.Name())
				plugins.plugins = append(plugins.plugins, p)
			}
		}
	}
	return plugins
}
