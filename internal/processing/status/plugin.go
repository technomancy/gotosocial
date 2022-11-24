package status

import (
	"os"
	"github.com/superseriousbusiness/gotosocial/internal/gtsmodel"

	"github.com/yuin/gopher-lua"
	"layeh.com/gopher-luar"
)

var l *lua.LState = nil
var initialized bool = false

func PluginProcess(s *gtsmodel.Status) {
	if initialized && l == nil {
		return
	} else if !initialized {
		initialized = true
		_, err := os.Stat("plugins/init.lua")
		if err != nil {
			return
		} else {
			l := lua.NewState()
			l.DoString("dofile('plugins/init.lua')")
		}
	}
	l.SetGlobal("status", luar.New(l, s))
	l.DoString("plugin(status)")
}
