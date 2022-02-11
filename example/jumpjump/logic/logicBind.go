package logic

import (
	"github.com/gitbufenshuo/gopen/example/htmlblockcustom/logic/logic_follow"
	"github.com/gitbufenshuo/gopen/example/jumpjump/logic/logic_changedong"
	"github.com/gitbufenshuo/gopen/example/jumpjump/logic/logic_jump"
	"github.com/gitbufenshuo/gopen/example/jumpjump/logic/logic_stdin"
	"github.com/gitbufenshuo/gopen/game"
)

func BindCustomLogic(gi *game.GlobalInfo) {
	gi.LogicSystem.BindLogicByName(gi,
		"stdin", logic_stdin.NewLogicStdin,
	)
	gi.LogicSystem.BindLogicByName(gi,
		"jump", logic_jump.NewLogicJump,
	)
	gi.LogicSystem.BindLogicByName(gi,
		"changedong", logic_changedong.NewLogicChangedong,
	)
	gi.LogicSystem.BindLogicByName(gi,
		"follow", logic_follow.NewLogicFollow,
	)

}
