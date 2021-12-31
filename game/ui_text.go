package game

import (
	"fmt"
	"math/rand"

	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/gitbufenshuo/gopen/game/common"
	"github.com/go-gl/gl/v4.1-core/gl"
)

const pixration float32 = 0.02

var textModelDefault *resource.Model
var InitDefaultTextOK bool

func InitDefaultText() {
	if InitDefaultTextOK {
		return
	}
	InitDefaultTextOK = true
	{
		textModelDefault = resource.NewQuadModel_LeftALign()
		textModelDefault.Upload()
	}
}

type UIText struct {
	gi              *GlobalInfo
	id              int
	renderComponent *resource.RenderComponent
	enabled         bool
	transform       *common.Transform
	shaderOP        *ShaderOP
	sortz           float32
	//
	content string
}

func NewUIText(gi *GlobalInfo) *UIText {
	InitDefaultText()
	uitext := new(UIText)
	uitext.enabled = true
	uitext.gi = gi
	uitext.sortz = 0.0005
	uitext.id = rand.Int()
	/////////////////////////
	uitext.renderComponent = new(resource.RenderComponent)
	uitext.renderComponent.ModelR = textModelDefault
	uitext.renderComponent.TextureR = resource.NewTexture()
	{
		uitext.renderComponent.ShaderR = buttonShaderDefault
		uitext.shaderOP = NewShaderOP()
		uitext.shaderOP.SetProgram(uitext.renderComponent.ShaderR.ShaderProgram())
		uitext.shaderOP.IfUI()
	}
	//
	uitext.transform = common.NewTransform()
	return uitext

}
func (uitext *UIText) ID_sg(_id ...int) int {
	if len(_id) == 0 {
		return uitext.id
	}
	if len(_id) > 1 {
		panic("len(_id)")
	}
	uitext.id = _id[0]
	return uitext.id
}

func (uitext *UIText) GetRenderComponent() *resource.RenderComponent {
	return uitext.renderComponent
}
func (uitext *UIText) Enabled() bool {
	return uitext.enabled
}

func (uitext *UIText) Enable() {
	uitext.enabled = true
}

func (uitext *UIText) Disable() {
	uitext.enabled = false
}

func (uitext *UIText) Start() {
}

func (uitext *UIText) Update() {

}

func (uitext *UIText) OnDraw() {
	// fmt.Println("uibutton,", uibutton.renderComponent.ShaderR)
	uitext.renderComponent.ShaderR.Active()
	uitext.renderComponent.ModelR.Active()
	uitext.renderComponent.TextureR.Active()
	mloc, lightloc, sortzloc := uitext.shaderOP.UniformLoc("model"),
		uitext.shaderOP.UniformLoc("light"),
		uitext.shaderOP.UniformLoc("sortz")

	//
	modelt := uitext.transform.WorldModel()
	gl.UniformMatrix4fv(mloc, 1, false, modelt.Address())
	gl.Uniform1f(lightloc, 1)
	gl.Uniform1f(sortzloc, 0.0005)
	// gl.Uniform1f(whrloc, uitext.gi.GetWHR())

}

func (uitext *UIText) OnDrawFinish() {

}

//
func (uitext *UIText) SetText(content string) {
	if uitext.content == content {
		return
	}
	fmt.Println("SetText", uitext.content, content)
	//
	uitext.content = content
	// re - gen - texture
	tr := uitext.renderComponent.TextureR
	tr.Clear()
	pixWidth := tr.GenFont(content, uitext.gi.FontConfig)
	tr.Upload()
	// re - scale - model
	uitext.transform.Scale.SetIndexValue(0, float32(pixWidth/30)*pixration) // 字图每16像素，是屏幕的1/20
	uitext.transform.Scale.SetIndexValue(1, pixration*1.4)
}

func (uitext *UIText) SortZ() float32 {
	return uitext.sortz
}
