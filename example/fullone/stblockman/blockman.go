package stblockman

import (
	"fmt"
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
		gi.AssetManager.CreateModelSilent("blockmanwheel.model", customModel)
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
		gi.AssetManager.CreateModelSilent("blockmanbody.model", customModel)
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
		gi.AssetManager.CreateModelSilent("blockmanhead.model", customModel)
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

	blockMan.AnimationCtl.LoadFromFile("blockman.dong")
	fmt.Println(blockMan.AnimationCtl.ModeList)
	return
}
