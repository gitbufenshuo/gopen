package manage_main

import (
	"fmt"

	"github.com/gitbufenshuo/gopen/example/jumpjump/logic/logic_jump"
	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/gameobjects"
	"github.com/gitbufenshuo/gopen/gameex/inputsystem"
	"github.com/gitbufenshuo/gopen/gameex/sceneloader"
	"github.com/go-gl/glfw/v3.1/glfw"
)

type ManageMain struct {
	gi *game.GlobalInfo
	*gameobjects.NilManageObject
	//
	which          int
	MainPlayer     game.GameObjectI
	MainPlayerJump *logic_jump.LogicJump
	SubPlayer      game.GameObjectI
	SubPlayerJump  *logic_jump.LogicJump
}

func NewManageMain(gi *game.GlobalInfo) *ManageMain {
	res := new(ManageMain)
	//
	res.NilManageObject = gameobjects.NewNilManageObject()
	res.gi = gi
	return res
}

func (lm *ManageMain) Start() {
	inputsystem.InitInputSystem(lm.gi)
	inputsystem.GetInputSystem().BeginWatchKey(int(glfw.KeyP))

	lm.MainPlayer = sceneloader.FindGameobjectByName("scenespec", "MainPlayer")
	logiclist := lm.MainPlayer.GetLogicSupport()
	for idx := range logiclist {
		if v, ok := logiclist[idx].(*logic_jump.LogicJump); ok {
			lm.MainPlayerJump = v
			lm.MainPlayerJump.Chosen = true
		}
	}
}

func (lm *ManageMain) clonePlayer() {
	if lm.SubPlayerJump != nil {
		return
	}
	logicx := lm.MainPlayerJump.GetLogicPosX()
	if lm.gi.CurFrame%30 == 0 {
		fmt.Println("MainPlayerLogicX", logicx)
	}
	if logicx < -5 {
		lm.SubPlayer = lm.gi.InstantiateGameObject(lm.MainPlayer)
		logiclist := lm.SubPlayer.GetLogicSupport()
		for idx := range logiclist {
			if v, ok := logiclist[idx].(*logic_jump.LogicJump); ok {
				lm.SubPlayerJump = v
			}
		}
	}
}

func (lm *ManageMain) Update() {
	lm.clonePlayer()
	if lm.SubPlayerJump == nil {
		return
	}
	if inputsystem.GetInputSystem().KeyDown(int(glfw.KeyP)) {
		lm.which++
	}
	//
	if lm.which%2 == 0 {
		lm.MainPlayerJump.Chosen = true
		lm.SubPlayerJump.Chosen = false
	} else {
		lm.MainPlayerJump.Chosen = false
		lm.SubPlayerJump.Chosen = true
	}
}
