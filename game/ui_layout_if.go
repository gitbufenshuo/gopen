package game

import "github.com/gitbufenshuo/gopen/game/common"

type UICanBeLayout interface {
	GetTransform() *common.Transform
	GetUISpec() *UISpec
}
