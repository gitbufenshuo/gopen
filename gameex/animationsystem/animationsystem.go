package animationsystem

import (
	"fmt"

	"github.com/gitbufenshuo/gopen/game"
)

func NewAniMoving(gbid int, bonename string) *game.AniMoving {
	res := new(game.AniMoving)
	//
	res.GBID = gbid
	res.BoneName = bonename
	return res
}

type AnimationSystem struct {
	gi *game.GlobalInfo
	//
	innerid          int64
	AnimationMataMap map[string]*AnimationMeta      // 动画资源
	ACRuntimeList    map[int64]*AnimationController // 这个是会执行的AC
	// ACStoreList      map[int]*AnimationController // 这个相当于AC库
	// MovingList       map[int][]*game.AniMoving    // gb 被AC管控Move。key 是 被管控gb的id
}

func NewAnimationSystem(gi *game.GlobalInfo) *AnimationSystem {
	res := new(AnimationSystem)
	res.gi = gi
	//
	res.AnimationMataMap = make(map[string]*AnimationMeta)
	res.ACRuntimeList = make(map[int64]*AnimationController)
	return res
}

func (as *AnimationSystem) AddAnimationMeta(amname string, am *AnimationMeta) {
	as.AnimationMataMap[amname] = am
}

func (as *AnimationSystem) allocID() int64 {
	as.innerid++
	return as.innerid
}

/*
	指定动画资源名称，创建一个AC
*/
func (as *AnimationSystem) CreateAC(amname string) game.AnimationControllerI {
	am := as.AnimationMataMap[amname]
	//
	acid := as.allocID()
	ac := NewAnimationController(acid)
	ac.UseAimationMeta(am)
	as.ACRuntimeList[acid] = ac
	fmt.Println(">>>>>>CreateAnimationController", amname)
	return ac
}

// gbid: 主gameobject id
// func (as *AnimationSystem) BindBoneNode(gbid int, bonename string, transform *game.Transform) {
// 	ac := as.GetAC(gbid)
// 	ac.BindBoneNode(bonename, transform)
// 	//
// 	movingGBID := transform.GB.ID_sg()
// 	newanimov := NewAniMoving(gbid, bonename)
// 	as.addMovingList(movingGBID, newanimov)
// }

// func (as *AnimationSystem) addMovingList(gbid int, mov *game.AniMoving) {
// 	if v, found := as.MovingList[gbid]; found {
// 		v = append(v, mov)
// 	} else {
// 		as.MovingList[gbid] = []*game.AniMoving{mov}
// 	}
// }

// 一个 gameobject 被删除了
// 应该将runtimeac 和storeac 中的相应元素删掉
// func (as *AnimationSystem) GameobjectDel(gbid int) {
// 	if _, found := as.ACRuntimeList[gbid]; found {
// 		delete(as.ACRuntimeList, gbid)
// 	}
// 	if _, found := as.ACStoreList[gbid]; found {
// 		delete(as.ACStoreList, gbid)
// 	}
// }

// 一个 gameobject 被剔除 system
// 应该将 runtimeac 中的相应元素删掉
// storeac 保留
// func (as *AnimationSystem) GameobjectDetach(gbid int) {
// 	if _, found := as.ACRuntimeList[gbid]; found {
// 		delete(as.ACRuntimeList, gbid)
// 	}
// }

func (as *AnimationSystem) Start() {
}
