package logic_main

import (
	"fmt"
	"math/rand"

	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/gitbufenshuo/gopen/game/gameobjects"
	"github.com/gitbufenshuo/gopen/game/supports/logicinner"
	"github.com/gitbufenshuo/gopen/matmath"
)

type LogicMain struct {
	gi *game.GlobalInfo
	*gameobjects.NilManageObject
	//
	modelresourcename string
	OjbectsMap        map[int]game.GameObjectI
	Dir               int
}

func NewLogicMain(gi *game.GlobalInfo) *LogicMain {
	res := new(LogicMain)
	//
	res.NilManageObject = gameobjects.NewNilManageObject()
	res.gi = gi
	res.OjbectsMap = make(map[int]game.GameObjectI)
	return res
}
func (lm *LogicMain) Start() {
	model := resource.NewBlockModel_BySpec(
		matmath.CreateVec4FromStr("0,0,0,0"),
		matmath.CreateVec4FromStr("1,1,1,1"),
	)
	lm.modelresourcename = "code.logicmain"
	fmt.Println("modelresourcename", lm.modelresourcename)
	lm.gi.AssetManager.CreateModelSilent(lm.modelresourcename, model)
}
func (lm *LogicMain) Update() {
	if lm.gi.CurFrame%100 == 0 {
		fmt.Println("now object count:", len(lm.OjbectsMap))
	}
	if lm.gi.CurFrame%3 == 0 {
		if lm.Dir > 0 {
			for idx := 0; idx != 10; idx++ {
				gb := gameobjects.NewBasicObject(lm.gi, lm.modelresourcename, "body.png.texture", "mvp_shader")
				lm.gi.AddGameObject(gb)
				gb.AddLogicSupport(logicinner.NewLogicRotate("1,1,1,1"))
				gb.Transform.Postion.SetValue3(10*rand.Float32()-5, 10+0.10*rand.Float32(), 10*rand.Float32()-5)
				gb.Transform.Scale.SetValue3(rand.Float32(), rand.Float32(), rand.Float32())
				lm.OjbectsMap[gb.ID_sg()] = gb
			}
			if len(lm.OjbectsMap) > 100 {
				lm.Dir = -1
				return
			}
		} else {
			for idx := 0; idx != 10; idx++ {
				for k, v := range lm.OjbectsMap {
					lm.gi.DelGameObject(v)
					delete(lm.OjbectsMap, k)
					break
				}
			}
			if len(lm.OjbectsMap) == 0 {
				lm.Dir = 1
				return
			}
		}
	}
}
