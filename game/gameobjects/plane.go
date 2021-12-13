package gameobjects

import (
	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
)

type PlaneObject struct {
	*BasicObject
	shaderProgram uint32
	shaderCtl     *game.ShaderCtl
	Rotating      bool
	////////////////////////
	cameraCircleRad float64
	cameraVertical  float64
	cameraR         float64
}

func NewPlane(gi *game.GlobalInfo, modelname, texturename string) *PlaneObject {
	innerBasic := NewBasicObject(gi, false)
	innerBasic.ModelAsset_sg(gi.AssetManager.FindByName(modelname))
	innerBasic.ShaderAsset_sg(gi.AssetManager.FindByName("mvp_shader"))
	innerBasic.TextureAsset_sg(gi.AssetManager.FindByName(texturename))
	innerBasic.DrawEnable_sg(true)

	one := new(PlaneObject)
	one.BasicObject = innerBasic
	one.shaderCtl = game.NewShaderCtl(one.ShaderAsset_sg().Resource.(*resource.ShaderProgram).ShaderProgram())
	one.cameraR = 20
	return one
}
func (co *PlaneObject) Start() {
	co.GI().GlobalFrameInfo.Debug = true
}
func (co *PlaneObject) Update() {
	if co.Rotating {
		co.Transform.Rotation.SetIndexValue(0, float32((co.GI().CurFrame)))
		co.Transform.Rotation.SetIndexValue(1, float32((co.GI().CurFrame)))
		// co.Transform.Postion.SetIndexValue(2, float32(math.Sin(float64(co.GI().CurFrame)*0.077))-2)
	}

}
func (co *PlaneObject) OnDraw() {
	co.shaderCtl.M = co.Transform.Model()
	co.shaderCtl.V = co.GI().View()
	co.shaderCtl.P = co.GI().Projection()
	co.shaderCtl.Rotation = co.Transform.RotationMAT4()
	co.shaderCtl.Upload(co)
}
func (co *PlaneObject) OnDrawFinish() {
}
