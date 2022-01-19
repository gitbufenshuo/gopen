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
	AnimationMataMap map[string]*AnimationMeta
	ACRuntimeList    map[int]*AnimationController // 这个是会执行的AC
	ACStoreList      map[int]*AnimationController // 这个相当于AC库
	MovingList       map[int][]*game.AniMoving
}

func NewAnimationSystem(gi *game.GlobalInfo) *AnimationSystem {
	res := new(AnimationSystem)
	res.gi = gi
	//
	res.AnimationMataMap = make(map[string]*AnimationMeta)
	res.ACRuntimeList = make(map[int]*AnimationController)
	res.ACStoreList = make(map[int]*AnimationController)
	//
	res.MovingList = make(map[int][]*game.AniMoving)
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
func (as *AnimationSystem) GetMoving(gbid int) []*game.AniMoving {
	if v, found := as.MovingList[gbid]; found {
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
	return ac
}

func (as *AnimationSystem) CloneAC(oldgbid, newgbid int) game.AnimationControllerI {
	ac := as.ACStoreList[oldgbid]
	newac := ac.Clone() // NodeList 没有复制
	as.ACRuntimeList[newgbid] = newac
	as.ACStoreList[newgbid] = newac
	return newac
}

// gbid: 主gameobject id
func (as *AnimationSystem) BindBoneNode(gbid int, bonename string, transform *game.Transform) {
	ac := as.GetAC(gbid)
	ac.BindBoneNode(bonename, transform)
	//
	movingGBID := transform.GB.ID_sg()
	newanimov := NewAniMoving(gbid, bonename)
	as.addMovingList(movingGBID, newanimov)
}

func (as *AnimationSystem) addMovingList(gbid int, mov *game.AniMoving) {
	if v, found := as.MovingList[gbid]; found {
		v = append(v, mov)
	} else {
		as.MovingList[gbid] = []*game.AniMoving{mov}
	}
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
	fmt.Println("AnimationSystem MovingList")
	for gbid, list := range as.MovingList {
		fmt.Println("          ", gbid, list[0].BoneName, list[0].GBID)
	}
	for _, oneac := range as.ACRuntimeList {
		oneac.Update()
	}
}
