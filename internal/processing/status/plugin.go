package status

import (
	"os"
	"context"
	"github.com/superseriousbusiness/gotosocial/internal/gtsmodel"
	"github.com/superseriousbusiness/gotosocial/internal/config"
	"github.com/superseriousbusiness/gotosocial/internal/log"

	"github.com/yuin/gopher-lua"
	"layeh.com/gopher-luar"
)

var l *lua.LState = nil
var initialized bool = false

func PluginProcess(ctx context.Context, s *gtsmodel.Status) {
	// this is tacky; should probably be its own config but that is a lot more
	// likely to cause git conflicts with upstream, so make it relative to
	// storage base path for now
	pluginPath := (config.GetStorageLocalBasePath() + "/../plugins/init.lua")
	if initialized && l == nil {
		return
	} else if !initialized {
		initialized = true
		_, err := os.Stat(pluginPath)
		if err != nil {
			log.Infof(ctx, "plugin not found: %s", pluginPath)
			return
		} else {
			l := lua.NewState()
			if err := l.DoFile(pluginPath); err != nil {
				log.Errorf(ctx, err.Error())
				return
			} else {
				log.Infof(ctx, "initialized plugin Lua")
			}
		}
	}
	// logs msg="processing status with thread: 0x0"
	// I'm guessing 0x0 means it was not initialized correctly, but shouldn't it
	// be nil if that were the case?
	log.Infof(ctx, "processing status with %s", l)
	// crashes with:
	// ()
	// 	runtime/panic.go:261
	// runtime.sigpanic()
	// 	runtime/signal_unix.go:881
	// ()
	// 	github.com/yuin/gopher-lua@v1.1.1/state.go:902
	// ()
	// 	github.com/yuin/gopher-lua@v1.1.1/state.go:1633
	// gopher-luar.New()
	// 	layeh.com/gopher-luar@v1.0.11/luar.go:112
	// status.PluginProcess()
	// 	github.com/superseriousbusiness/gotosocial/internal/processing/status/plugin.go:39
	// status.(*Processor).Create()
	// 	github.com/superseriousbusiness/gotosocial/internal/processing/status/create.go:151
	// statuses.(*Module).StatusCreatePOSTHandler()
	// 	github.com/superseriousbusiness/gotosocial/internal/api/client/statuses/statuscreate.go:295
	// ()
	// 	github.com/gin-gonic/gin@v1.10.0/context.go:185

	l.SetGlobal("status", luar.New(l, s))
	log.Infof(ctx, "running plugin")
	l.DoString("plugin(status)")
}
