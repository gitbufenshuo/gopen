package gameobjects

import (
	"github.com/gitbufenshuo/gopen/game"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type BlockObject struct {
	*BasicObject
	shaderProgram uint32
	Rotating      bool
	Color         []float32
}

func NewBlock(gi *game.GlobalInfo, modelname, texturename string) *BlockObject {
	innerBasic := NewBasicObject(gi, false)
	innerBasic.ModelAsset_sg(gi.AssetManager.FindByName(modelname))
	innerBasic.ShaderAsset_sg(gi.AssetManager.FindByName("mvp_shader"))
	innerBasic.TextureAsset_sg(gi.AssetManager.FindByName(texturename))
	innerBasic.DrawEnable_sg(true)

	one := new(BlockObject)
	one.BasicObject = innerBasic
	one.Color = []float32{0, 0, 0}
	return one
}
func (co *BlockObject) Start() {
	co.GI().GlobalFrameInfo.Debug = true

}
func (co *BlockObject) Update() {
	if co.Rotating {
		co.Transform.Rotation.SetIndexValue(1, float32((co.GI().CurFrame)))
	}
}
func (co *BlockObject) OnDraw() {
	sop := co.ShaderOP()
	gl.Uniform3f(sop.UniformLoc("u_Color"), co.Color[0], co.Color[1], co.Color[2])
}
func (co *BlockObject) OnDrawFinish() {
	sop := co.ShaderOP()
	gl.Uniform3f(sop.UniformLoc("u_Color"), 1, 1, 1)

}
