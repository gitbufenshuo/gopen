package animationsystem

import (
	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/gameobjects"
)

type AnimationSystem struct {
	*gameobjects.NilManageObject
	gi *game.GlobalInfo
	//
	AnimationMataMap        map[string]*AnimationMeta
	AnimationControllerList []*AnimationController
}

func NewAnimationSystem(gi *game.GlobalInfo) *AnimationSystem {
	res := new(AnimationSystem)
	res.gi = gi
	res.NilManageObject = gameobjects.NewNilManageObject()
	//
	res.AnimationMataMap = make(map[string]*AnimationMeta)
	return res
}

func (as *AnimationSystem) AddAnimationMeta(amname string, am *AnimationMeta) {
	as.AnimationMataMap[amname] = am
}

func (as *AnimationSystem) CreateAnimationController(amname string) game.AnimationControllerI {
	am := as.AnimationMataMap[amname]
	//
	ac := NewAnimationController()
	ac.UseAimationMeta(am)
	as.AnimationControllerList = append(as.AnimationControllerList, ac)
	var aci game.AnimationControllerI = ac
	return aci
}

func (as *AnimationSystem) Start() {
}
func (as *AnimationSystem) Update() {
	for _, oneac := range as.AnimationControllerList {
		oneac.Update()
	}
}
