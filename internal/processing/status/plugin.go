package status

import (
	"github.com/superseriousbusiness/gotosocial/internal/gtsmodel"

	"github.com/yuin/gopher-lua"
	"layeh.com/gopher-luar"
)

func InitPlugin(name string) *lua.LState {
	l := lua.NewState()
	l.DoFile("plugins/" + name + ".lua")
	return l
}

func PluginProcess(l *lua.LState, s *gtsmodel.Status) {
	l.SetGlobal("status", luar.New(l, s))
	l.DoString("plugin(status)")
}