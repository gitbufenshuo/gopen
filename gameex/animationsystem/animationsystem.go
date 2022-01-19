package animationsystem

import (
	"fmt"

	"github.com/gitbufenshuo/gopen/game"
)

type AnimationSystem struct {
	gi *game.GlobalInfo
	//
	AnimationMataMap map[string]*AnimationMeta
	ACRuntimeList    map[int]*AnimationController // 这个是会执行的AC
	ACStoreList      map[int]*AnimationController // 这个相当于AC库
}

func NewAnimationSystem(gi *game.GlobalInfo) *AnimationSystem {
	res := new(AnimationSystem)
	res.gi = gi
	//
	res.AnimationMataMap = make(map[string]*AnimationMeta)
	res.ACRuntimeList = make(map[int]*AnimationController)
	res.ACStoreList = make(map[int]*AnimationController)
	return res
}

func (as *AnimationSystem) AddAnimationMeta(amname string, am *AnimationMeta) {
	as.AnimationMataMap[amname] = am
}

func (as *AnimationSystem) GetAC(gbid int) game.AnimationControllerI {
	if v, found := as.ACStoreList[gbid]; found {
		return v
	}
	return nil
}

func (as *AnimationSystem) CreateAC(amname string, gbid int) game.AnimationControllerI {
	am := as.AnimationMataMap[amname]
	//
	ac := NewAnimationController()
	ac.UseAimationMeta(am)
	as.ACRuntimeList[gbid] = ac
	as.ACStoreList[gbid] = ac
	fmt.Println(">>>>>>CreateAnimationController", amname, gbid)
	var aci game.AnimationControllerI = ac
	return aci
}

// 一个 gameobject 被删除了
// 应该将runtimeac 和storeac 中的相应元素删掉
func (as *AnimationSystem) GameobjectDel(gbid int) {
	if _, found := as.ACRuntimeList[gbid]; found {
		delete(as.ACRuntimeList, gbid)
	}
	if _, found := as.ACStoreList[gbid]; found {
		delete(as.ACStoreList, gbid)
	}
}

// 一个 gameobject 被剔除 system
// 应该将 runtimeac 中的相应元素删掉
// storeac 保留
func (as *AnimationSystem) GameobjectDetach(gbid int) {
	if _, found := as.ACRuntimeList[gbid]; found {
		delete(as.ACRuntimeList, gbid)
	}
}

func (as *AnimationSystem) Start() {
}

func (as *AnimationSystem) Update() {
	for _, oneac := range as.ACRuntimeList {
		oneac.Update()
	}
}
