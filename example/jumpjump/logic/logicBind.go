package logic

import (
	"github.com/gitbufenshuo/gopen/example/jumpjump/logic/logic_changedong"
	"github.com/gitbufenshuo/gopen/example/jumpjump/logic/logic_clone"
	"github.com/gitbufenshuo/gopen/example/jumpjump/logic/logic_jump"
	"github.com/gitbufenshuo/gopen/example/jumpjump/logic/logic_stdin"
	"github.com/gitbufenshuo/gopen/example/jumpjump/logic/logic_walk"
	"github.com/gitbufenshuo/gopen/game"
)

func BindCustomLogic(gi *game.GlobalInfo) {
	gi.LogicSystem.BindLogicByName(gi,
		"jump", logic_jump.NewLogicJump,
	)
	gi.LogicSystem.BindLogicByName(gi,
		"walk", logic_walk.NewLogicWalk,
	)
	gi.LogicSystem.BindLogicByName(gi,
		"changedong", logic_changedong.NewLogicChangedong,
	)
	gi.LogicSystem.BindLogicByName(gi,
		"clone", logic_clone.NewLogicClone,
	)
	gi.LogicSystem.BindLogicByName(gi,
		"stdin", logic_stdin.NewLogicStdin,
	)
}
