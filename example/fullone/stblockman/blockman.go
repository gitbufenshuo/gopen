package stblockman

import (
	"math/rand"

	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/gitbufenshuo/gopen/game/common"
	"github.com/gitbufenshuo/gopen/game/gameobjects"
)

type BlockManWheel struct {
	*gameobjects.BlockObject
	InnerModel *resource.Model
}

func NewBlockManWheel(gi *game.GlobalInfo) *BlockManWheel {
	{
		customModel := resource.NewBlockModel()

		for idx := 0; idx != 24; idx++ {
			customModel.Vertices[idx*8] *= 0.3
			customModel.Vertices[idx*8+2] *= 0.9
		}
		gi.AssetManager.CreateModel("blockmanwheel.model", customModel)
	}
	block := gameobjects.NewBlock(gi, "blockmanwheel.model", "grid.png.texuture")
	block.Color = []float32{1, 1, 1}
	wheel := new(BlockManWheel)
	///
	wheel.BlockObject = block
	wheel.InnerModel = wheel.ModelAsset_sg().Resource.(*resource.Model)
	wheel.Transform.Postion.SetIndexValue(1, -1)
	return wheel
}

type BlockManHand struct {
	*gameobjects.BlockObject
	InnerModel *resource.Model
}

func NewBlockManHand(gi *game.GlobalInfo) *BlockManHand {
	if res := gi.AssetManager.FindByName("blockmanhand.model"); res == nil {

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
		gi.AssetManager.CreateModel("blockmanhand.model", customModel)
	}
	block := gameobjects.NewBlock(gi, "blockmanhand.model", "hand.png.texuture")
	block.Color = []float32{1, 1, 1}
	hand := new(BlockManHand)
	///
	hand.BlockObject = block
	hand.InnerModel = hand.ModelAsset_sg().Resource.(*resource.Model)
	return hand
}

type BlockManBody struct {
	*gameobjects.BlockObject
	InnerModel *resource.Model
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
		}
		gi.AssetManager.CreateModel("blockmanbody.model", customModel)
	}
	block := gameobjects.NewBlock(gi, "blockmanbody.model", "body.png.texuture")
	block.Color = []float32{1, 1, 1}
	body := new(BlockManBody)
	///
	body.BlockObject = block
	body.InnerModel = body.ModelAsset_sg().Resource.(*resource.Model)
	return body
}

type BlockManHead struct {
	*gameobjects.BlockObject
	InnerModel *resource.Model
}

func NewBlockManHead(gi *game.GlobalInfo) *BlockManHead {
	{
		customModel := resource.NewBlockModel()

		for idx := 0; idx != 24; idx++ {
			customModel.Vertices[idx*8] *= 0.8
			customModel.Vertices[idx*8+1] *= 0.8
			customModel.Vertices[idx*8+2] *= 0.8
		}
		gi.AssetManager.CreateModel("blockmanhead.model", customModel)
	}
	block := gameobjects.NewBlock(gi, "blockmanhead.model", "head.png.texuture")
	block.Color = []float32{1, 1, 1}
	head := new(BlockManHead)
	///
	head.BlockObject = block
	head.InnerModel = head.ModelAsset_sg().Resource.(*resource.Model)
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

func (bmc *BlockManCore) Update() {
	v := float32(bmc.GI().CurFrame) * 1.2
	bmc.Transform.Rotation.SetIndexValue(1, v)
}

type BlockMan struct {
	ID        int
	Core      *BlockManCore
	Body      *BlockManBody
	Head      *BlockManHead
	HandLeft  *BlockManHand
	HandRight *BlockManHand
	Wheel     *BlockManWheel
	//
	AnimationCtl *common.AnimationController
}

func (bm *BlockMan) Start() {

}

func (bm *BlockMan) Update() {
	bm.AnimationRun()
	if bm.Body.GI().CurFrame%111 == 0 {
		rint := rand.Int()
		rint %= len(bm.AnimationCtl.ModeList)
		bm.AnimationCtl.ChangeMode(bm.AnimationCtl.ModeList[rint])
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
	//
	blockMan.Core = NewBlockManCore(gi)
	blockMan.Body = NewBlockManBody(gi)
	blockMan.Head = NewBlockManHead(gi)
	blockMan.HandLeft = NewBlockManHand(gi)
	blockMan.HandLeft.Transform.Postion.SetValue2(-0.7, 1)
	blockMan.HandRight = NewBlockManHand(gi)
	blockMan.HandRight.Transform.Postion.SetValue2(0.7, 1)
	blockMan.Wheel = NewBlockManWheel(gi)

	gi.AddGameObject(blockMan.Core)
	gi.AddGameObject(blockMan.Body)
	gi.AddGameObject(blockMan.Head)
	gi.AddGameObject(blockMan.HandLeft)
	gi.AddGameObject(blockMan.HandRight)
	gi.AddGameObject(blockMan.Wheel)
	//
	blockMan.Body.Transform.SetParent(blockMan.Core.Transform)
	blockMan.Head.Transform.SetParent(blockMan.Body.Transform)
	blockMan.HandLeft.Transform.SetParent(blockMan.Body.Transform)
	blockMan.HandRight.Transform.SetParent(blockMan.Body.Transform)
	blockMan.Wheel.Transform.SetParent(blockMan.Body.Transform)
	//
	blockMan.CreateAnimation()
	return blockMan
}

func (blockMan *BlockMan) CreateAnimation() {
	blockMan.AnimationCtl = common.NewAnimationController()
	blockMan.AnimationCtl.BindBoneNode(
		blockMan.Head.Transform,
		blockMan.Body.Transform,
		blockMan.HandLeft.Transform,
		blockMan.HandRight.Transform,
		blockMan.Wheel.Transform)
	{
		STOPMODE := make([]*common.AnimationFrame, 80)
		for idx := 0; idx != 80; idx++ {
			STOPMODE[idx] = &common.AnimationFrame{
				HeadStatus:      common.NewBoneSatus(0, 0, 0, 0, float32(idx)*2, 0),
				BodyStatus:      common.NewBoneSatus(0, 0, 0, 0, 0, 0),
				HandLeftStatus:  common.NewBoneSatus(0, 0, 0, 0, 0, 0),
				HandRightStatus: common.NewBoneSatus(0, 0, 0, 0, 0, 0),
				WheelStatus:     common.NewBoneSatus(0, 0, 0, 0, 0, 0),
			}
			if idx > 60 {
				STOPMODE[idx] = &common.AnimationFrame{
					HeadStatus:      common.NewBoneSatus(0, 0, 0, 0, float32(60)*0.5, 0),
					BodyStatus:      common.NewBoneSatus(0, 0, 0, 0, 0, 0),
					HandLeftStatus:  common.NewBoneSatus(0, 0, 0, 0, 0, 0),
					HandRightStatus: common.NewBoneSatus(0, 0, 0, 0, 0, 0),
					WheelStatus:     common.NewBoneSatus(0, 0, 0, 0, 0, 0),
				}
			}
		}

		blockMan.AnimationCtl.AddMode("__init", STOPMODE)
		blockMan.AnimationCtl.CurMode = "__init"
		blockMan.AnimationCtl.ModeList = append(blockMan.AnimationCtl.ModeList, "__init")
	}
	{
		MOVINGMODE := make([]*common.AnimationFrame, 60)
		for idx := 0; idx != 15; idx++ {
			MOVINGMODE[idx] = &common.AnimationFrame{
				HeadStatus:      common.NewBoneSatus(0, 0, 0, 0, 0, 0),
				BodyStatus:      common.NewBoneSatus(0, 0, 0, 0, 0, 0),
				HandLeftStatus:  common.NewBoneSatus(0, 0, 0, float32(idx)*4, 0, 0),
				HandRightStatus: common.NewBoneSatus(0, 0, 0, -float32(idx)*4, 0, 0),
				WheelStatus:     common.NewBoneSatus(0, 0, 0, 0, 0, 0),
			}
		}
		for idx := 15; idx != 30; idx++ {
			MOVINGMODE[idx] = MOVINGMODE[30-idx-1]
		}
		for idx := 30; idx != 45; idx++ {
			MOVINGMODE[idx] = &common.AnimationFrame{
				HeadStatus:      common.NewBoneSatus(0, 0, 0, 0, 0, 0),
				BodyStatus:      common.NewBoneSatus(0, 0, 0, 0, 0, 0),
				HandLeftStatus:  common.NewBoneSatus(0, 0, 0, -float32(idx-30)*4, 0, 0),
				HandRightStatus: common.NewBoneSatus(0, 0, 0, float32(idx-30)*4, 0, 0),
				WheelStatus:     common.NewBoneSatus(0, 0, 0, 0, 0, 0),
			}
		}
		for idx := 45; idx != 60; idx++ {
			MOVINGMODE[idx] = MOVINGMODE[60-idx-1+30]
		}
		blockMan.AnimationCtl.AddMode("MOVING", MOVINGMODE)
		blockMan.AnimationCtl.CurMode = "MOVING"
		blockMan.AnimationCtl.ModeList = append(blockMan.AnimationCtl.ModeList, "MOVING")

	}
	{
		MOVINGMODE := make([]*common.AnimationFrame, 15)
		for idx := 0; idx != 15; idx++ {
			MOVINGMODE[idx] = &common.AnimationFrame{
				HeadStatus:      common.NewBoneSatus(0, 0, 0, 0, 0, 0),
				BodyStatus:      common.NewBoneSatus(0, float32(idx)*0.05, 0, 0, 0, 0),
				HandLeftStatus:  common.NewBoneSatus(0, 0, 0, float32(idx)*4, 0, 0),
				HandRightStatus: common.NewBoneSatus(0, 0, 0, float32(idx)*4, 0, 0),
				WheelStatus:     common.NewBoneSatus(0, 0, 0, 0, 0, 0),
			}
		}
		blockMan.AnimationCtl.AddMode("JUMPING", MOVINGMODE)
		blockMan.AnimationCtl.CurMode = "JUMPING"
		blockMan.AnimationCtl.ModeList = append(blockMan.AnimationCtl.ModeList, "JUMPING")

	}
	{
		FIREMODE := make([]*common.AnimationFrame, 60)
		for idx := 0; idx != 20; idx++ {
			FIREMODE[idx] = &common.AnimationFrame{
				HeadStatus:      common.NewBoneSatus(0, 0, 0, 0, 6.5, 0),
				BodyStatus:      common.NewBoneSatus(rand.Float32()/10-0.05, rand.Float32()/10-0.05, rand.Float32()/10-0.05, 0, float32(idx)*5.5, 0),
				HandLeftStatus:  common.NewBoneSatus(0, 0, 0, -float32(idx)*4.5, 0, 0),
				HandRightStatus: common.NewBoneSatus(0, 0, 0, 0, 0, 0),
				WheelStatus:     common.NewBoneSatus(0, 0, 0, 0, 0, 0),
			}
		}
		for idx := 20; idx != 40; idx++ {
			FIREMODE[idx] = &common.AnimationFrame{
				HeadStatus:      common.NewBoneSatus(0, 0, 0, 0, 0, 0),
				BodyStatus:      common.NewBoneSatus(rand.Float32()/10-0.05, rand.Float32()/10-0.05, rand.Float32()/10-0.05, 0, float32(20)*5.5, 0),
				HandLeftStatus:  common.NewBoneSatus(0, 0, 0, -float32(20)*4.5, 0, 0),
				HandRightStatus: common.NewBoneSatus(0, 0, 0, 0, 0, 0),
				WheelStatus:     common.NewBoneSatus(0, 0, 0, 0, 0, 0),
			}
		}
		for idx := 40; idx != 60; idx++ {
			FIREMODE[idx] = &common.AnimationFrame{
				HeadStatus:      common.NewBoneSatus(0, 0, 0, 0, 0, 0),
				BodyStatus:      common.NewBoneSatus(rand.Float32()/10-0.05, rand.Float32()/10-0.05, rand.Float32()/10-0.05, 0, float32(20)*5.5, 0),
				HandLeftStatus:  common.NewBoneSatus(rand.Float32()/2-0.25, 0, float32(idx-40), -float32(20)*4.5, float32(idx-40)*20, 0),
				HandRightStatus: common.NewBoneSatus(0, 0, 0, 0, 0, 0),
				WheelStatus:     common.NewBoneSatus(0, 0, 0, 0, 0, 0),
			}
		}
		blockMan.AnimationCtl.AddMode("FIREMODE", FIREMODE)
		blockMan.AnimationCtl.CurMode = "FIREMODE"
		blockMan.AnimationCtl.ModeList = append(blockMan.AnimationCtl.ModeList, "FIREMODE")

	}
}
