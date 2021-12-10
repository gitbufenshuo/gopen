package gameobjects

import (
	"math/rand"

	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
)

type BlockObject struct {
	*BasicObject
	shaderProgram uint32
	shaderCtl     *game.ShaderCtl
	Rotating      bool
}

func NewBlock(gi *game.GlobalInfo, modelname, texturename string) *BlockObject {
	innerBasic := NewBasicObject(gi, false)
	innerBasic.ModelAsset_sg(gi.AssetManager.FindByName(modelname))
	innerBasic.ShaderAsset_sg(gi.AssetManager.FindByName("mvp_shader"))
	innerBasic.TextureAsset_sg(gi.AssetManager.FindByName(texturename))
	innerBasic.DrawEnable_sg(true)

	one := new(BlockObject)
	one.BasicObject = innerBasic
	one.shaderCtl = game.NewShaderCtl(one.ShaderAsset_sg().Resource.(*resource.ShaderProgram).ShaderProgram())
	return one
}
func (co *BlockObject) Start() {
	co.GI().GlobalFrameInfo.Debug = true
}
func (co *BlockObject) Update() {
	if co.Rotating {
		co.Transform.Rotation.SetIndexValue(1, float32((co.GI().CurFrame)))
		// co.Transform.Postion.SetIndexValue(2, float32(math.Sin(float64(co.GI().CurFrame)*0.077))-2)
	}

}
func (co *BlockObject) OnDraw() {
	co.shaderCtl.M = co.Transform.Model()
	co.shaderCtl.Rotation = co.Transform.RotationMAT4()
	if co.Transform.Parent != nil { // not root
		parentM := co.Transform.Parent.Model()
		co.shaderCtl.M.RightMul_InPlace(&parentM)
		parentR := co.Transform.Parent.RotationMAT4()
		co.shaderCtl.Rotation.RightMul_InPlace(&parentR)
	}
	co.shaderCtl.V = co.GI().View()
	co.shaderCtl.P = co.GI().Projection()
	co.shaderCtl.Upload(co)
	//
	co.shaderCtl.UniformU_Colur(rand.Float32(), rand.Float32(), rand.Float32())
}
func (co *BlockObject) OnDrawFinish() {
	co.shaderCtl.UniformU_Colur(1, 1, 1)
}
