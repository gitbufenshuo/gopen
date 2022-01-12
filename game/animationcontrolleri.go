package game

import "github.com/gitbufenshuo/gopen/game/common"

type AnimationControllerI interface {
	BindBoneNode(name string, transform *common.Transform)
	RecordInitFrame()
}
type AnimationSystemI interface {
	CreateAnimationController(amname string) AnimationControllerI
	Update()
}
