package gameobjects

import (
	"github.com/gitbufenshuo/gopen/game"
)

type PlaneObject struct {
	*BasicObject
	shaderProgram uint32
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
}
func (co *PlaneObject) OnDrawFinish() {
}
