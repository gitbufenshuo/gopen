package logic_changedong

import (
	"fmt"

	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/supports"
	"github.com/gitbufenshuo/gopen/gameex/inputsystem"
	"github.com/go-gl/glfw/v3.1/glfw"
)

type LogicChangedong struct {
	gi *game.GlobalInfo
	*supports.NilLogic
	//
	bootok   bool
	ac       game.AnimationControllerI
	modelist []string
	modeidx  int
}

func NewLogicChangedong(gi *game.GlobalInfo) game.LogicSupportI {
	res := new(LogicChangedong)
	//
	res.NilLogic = supports.NewNilLogic()
	res.gi = gi
	return res
}

func (lc *LogicChangedong) boot(gb game.GameObjectI) {
	// input system
	inputsystem.InitInputSystem(lc.gi)
	inputsystem.GetInputSystem().BeginWatchKey(int(glfw.KeyS))
	lc.gi.SetInputSystem(inputsystem.GetInputSystem())

	lc.ac = lc.gi.AnimationSystem.GetAnimationController(gb.ID_sg())
	if lc.ac != nil {
		lc.modelist = lc.ac.GetModeList()
		fmt.Println("lc.modelist", lc.modelist)
	}
}

func (lc *LogicChangedong) Update(gb game.GameObjectI) {
	if !lc.bootok {
		lc.boot(gb)
		lc.bootok = true
		return
	}
	if lc.gi.InputSystemManager.KeyUp(int(glfw.KeyS)) {
		if lc.ac != nil {
			lc.modeidx++
			lc.modeidx %= len(lc.modelist)
			lc.ac.ChangeMode(lc.modelist[lc.modeidx])
		}
	}
}
