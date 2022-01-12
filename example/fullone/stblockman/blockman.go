package stblockman

import (
	"fmt"
	"math/rand"

	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/gitbufenshuo/gopen/game/common"
	"github.com/gitbufenshuo/gopen/game/gameobjects"
)

type BlockManLeg struct {
	*gameobjects.BasicObject
}

func NewBlockManLeg(gi *game.GlobalInfo) *BlockManLeg {
	{
		customModel := resource.NewBlockModel()
		for idx := 0; idx != 24; idx++ {
			if customModel.Vertices[idx*8] < 0 {
				customModel.Vertices[idx*8] = -0.1
			} else {
				customModel.Vertices[idx*8] = 0.1
			}
			customModel.Vertices[idx*8+1] -= 0.5
			customModel.Vertices[idx*8+2] *= 0.2
		}
		gi.AssetManager.CreateModelSilent("blockmanleg.model", customModel)
	}
	blockobject := gameobjects.NewBlockObject(gi, "blockmanleg.model", "grid.png.texture")

	wheel := new(BlockManLeg)
	///
	wheel.BasicObject = blockobject
	return wheel
}

type BlockManHand struct {
	*gameobjects.BasicObject
}

func NewBlockManHand(gi *game.GlobalInfo) *BlockManHand {
	{
		customModel := resource.NewBlockModel()
		for idx := 0; idx != 24; idx++ {
			if customModel.Vertices[idx*8] < 0 {
				customModel.Vertices[idx*8] = -0.1
			} else {
				customModel.Vertices[idx*8] = 0.1
			}
			customModel.Vertices[idx*8+1] -= 0.5
			customModel.Vertices[idx*8+2] *= 0.2
		}
		gi.AssetManager.CreateModelSilent("blockmanhand.model", customModel)
	}
	blockobject := gameobjects.NewBlockObject(gi, "blockmanhand.model", "hand.png.texture")
	hand := new(BlockManHand)
	///
	hand.BasicObject = blockobject
	return hand
}

type BlockManBody struct {
	*gameobjects.BasicObject
}

func NewBlockManBody(gi *game.GlobalInfo) *BlockManBody {
	{
		customModel := resource.NewBlockModel()

		for idx := 0; idx != 24; idx++ {
			if customModel.Vertices[idx*8+1] < 0 {
				customModel.Vertices[idx*8+1] = -1
			} else {
				customModel.Vertices[idx*8+1] = 1
			}
			customModel.Vertices[idx*8+2] *= 0.3
		}
		gi.AssetManager.CreateModelSilent("blockmanbody.model", customModel)
	}
	blockobject := gameobjects.NewBlockObject(gi, "blockmanbody.model", "body.png.texture")
	body := new(BlockManBody)
	body.BasicObject = blockobject
	logic := NewLogicColorControl()
	body.AddLogicSupport(logic)
	///
	return body
}

type BlockManHead struct {
	*gameobjects.BasicObject
}

func NewBlockManHead(gi *game.GlobalInfo) *BlockManHead {
	{
		customModel := resource.NewBlockModel()

		for idx := 0; idx != 24; idx++ {
			customModel.Vertices[idx*8] *= 0.8
			customModel.Vertices[idx*8+1] *= 0.8
			customModel.Vertices[idx*8+2] *= 0.8
			if idx >= 4 {
				customModel.Vertices[idx*8+3] = 0.8
				customModel.Vertices[idx*8+4] = 0.8
			}
		}
		gi.AssetManager.CreateModelSilent("blockmanhead.model", customModel)
	}
	blockobject := gameobjects.NewBlockObject(gi, "blockmanhead.model", "head.png.texture")
	head := new(BlockManHead)
	///
	head.BasicObject = blockobject
	head.Transform.Postion.SetIndexValue(1, 1.5)
	return head
}

type BlockManCore struct {
	*gameobjects.NilObject
}

func NewBlockManCore(gi *game.GlobalInfo) *BlockManCore {
	core := new(BlockManCore)
	///
	core.NilObject = gameobjects.NewNilObject(gi)
	return core
}

type BlockMan struct {
	gi        *game.GlobalInfo
	ID        int
	Core      *BlockManCore
	Body      *BlockManBody
	Head      *BlockManHead
	HandLeft  *BlockManHand
	HandRight *BlockManHand
	LegLeft   *BlockManLeg
	LegRight  *BlockManLeg
	//
	AnimationCtl *common.AnimationController
}

func (bm *BlockMan) Start() {

}

func (bm *BlockMan) Update() {
	gi := bm.gi
	// gi.MainCamera.Transform.Rotation.SetValue3(0, float32(gi.CurFrame), 0)

	bm.AnimationRun()
	if gi.CurFrame%111 == 0 {
		rint := rand.Int()
		rint %= len(bm.AnimationCtl.AM.ModeList)
		bm.AnimationCtl.ChangeMode(bm.AnimationCtl.AM.ModeList[rint])
	}
}
func (bm *BlockMan) AnimationRun() {
	//
	if bm.AnimationCtl == nil {
		return
	}
	//
	bm.AnimationCtl.Update()
}

func (bm *BlockMan) ID_sg(_id ...int) int {
	if len(_id) == 0 {
		return bm.ID
	}
	if len(_id) > 1 {
		panic("len(_id)")
	}
	bm.ID = _id[0]
	return bm.ID
}

func NewBlockMan(gi *game.GlobalInfo) *BlockMan {
	blockMan := new(BlockMan)
	blockMan.gi = gi
	// return blockMan
	//
	blockMan.Core = NewBlockManCore(gi)
	blockMan.Body = NewBlockManBody(gi)
	blockMan.Head = NewBlockManHead(gi)
	blockMan.HandLeft = NewBlockManHand(gi)
	blockMan.HandLeft.Transform.Postion.SetValue2(-0.7, 1)
	blockMan.HandRight = NewBlockManHand(gi)
	blockMan.HandRight.Transform.Postion.SetValue2(0.7, 1)
	blockMan.LegLeft = NewBlockManLeg(gi)
	blockMan.LegLeft.Transform.Postion.SetValue2(0.3, -1)
	blockMan.LegRight = NewBlockManLeg(gi)
	blockMan.LegRight.Transform.Postion.SetValue2(-0.3, -1)
	gi.AddGameObject(blockMan.Core)
	gi.AddGameObject(blockMan.Body)
	gi.AddGameObject(blockMan.Head)
	gi.AddGameObject(blockMan.HandLeft)
	gi.AddGameObject(blockMan.HandRight)
	gi.AddGameObject(blockMan.LegLeft)
	gi.AddGameObject(blockMan.LegRight)
	//
	blockMan.Body.Transform.SetParent(blockMan.Core.Transform)
	blockMan.Head.Transform.SetParent(blockMan.Body.Transform)
	blockMan.HandLeft.Transform.SetParent(blockMan.Body.Transform)
	blockMan.HandRight.Transform.SetParent(blockMan.Body.Transform)
	blockMan.LegLeft.Transform.SetParent(blockMan.Body.Transform)
	blockMan.LegRight.Transform.SetParent(blockMan.Body.Transform)
	//
	blockMan.CreateAnimation()
	return blockMan
}

func (blockMan *BlockMan) CreateAnimation() {
	blockMan.AnimationCtl = common.NewAnimationController()
	am := common.LoadAnimationMetaFromFile("blockman.dong")
	blockMan.AnimationCtl.UseAimationMeta(am)
	blockMan.AnimationCtl.BindBoneNodeList(
		[]*common.Transform{
			blockMan.Head.Transform,
			blockMan.Body.Transform,
			blockMan.HandLeft.Transform,
			blockMan.HandRight.Transform,
			blockMan.LegLeft.Transform,
			blockMan.LegRight.Transform,
		},
	)
	fmt.Println(blockMan.AnimationCtl.AM.ModeList)
	return
}
