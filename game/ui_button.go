package game

import (
	"fmt"

	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/gitbufenshuo/gopen/game/common"
	"github.com/go-gl/gl/v4.1-core/gl"
)

var buttonModelDefault *resource.Model
var buttonTextureDefault *resource.Texture
var buttonShaderDefault *resource.ShaderProgram

var InitDefaultButtonOK bool

func InitDefaultButton() {
	if InitDefaultButtonOK {
		return
	}
	InitDefaultButtonOK = true
	{
		buttonModelDefault = resource.NewQuadModel()
		for idx := 0; idx != 4; idx++ {
			// buttonModelDefault.Vertices[idx*5+0] *= 2
			// buttonModelDefault.Vertices[idx*5+1] *= 1
		}
		buttonModelDefault.Upload()
	}
	//
	buttonTextureDefault = resource.NewTexture()
	buttonTextureDefault.GenDefault(1, 1)
	//
	buttonShaderDefault = resource.NewShaderProgram()
	buttonShaderDefault.ReadFromText(resource.ShaderUIButtonText.Vertex, resource.ShaderUIButtonText.Fragment)
	buttonShaderDefault.Upload()
	fmt.Println("InitDefaultButton", buttonShaderDefault.ShaderProgram())
}

type UIButton struct {
	gi              *GlobalInfo
	id              int
	renderComponent *resource.RenderComponent
	enabled         bool
	transform       *common.Transform
	shaderOP        *ShaderOP
	// a_model_loc     int32
	// u_light_loc     int32
	// u_sortz_loc     int32
	sortz float32
	//
	uitext *UIText
}

func NewDefaultUIButton(gi *GlobalInfo) *UIButton {
	InitDefaultButton()
	uibutton := new(UIButton)
	uibutton.sortz = 0.001
	uibutton.gi = gi
	/////////////////////////
	uibutton.renderComponent = new(resource.RenderComponent)
	uibutton.renderComponent.ModelR = buttonModelDefault
	uibutton.renderComponent.TextureR = buttonTextureDefault
	{
		uibutton.renderComponent.ShaderR = buttonShaderDefault
		uibutton.shaderOP = NewShaderOP()
		uibutton.shaderOP.SetProgram(uibutton.renderComponent.ShaderR.ShaderProgram())
		uibutton.shaderOP.IfUI()
	}
	//
	uibutton.transform = common.NewTransform()
	//
	uibutton.uitext = NewUIText(gi)
	gi.AddUIObject(uibutton.uitext)
	uibutton.uitext.SetText("Hello Golang")
	return uibutton
}
func (uibutton *UIButton) ChangeTexture(textureR *resource.Texture) {
	uibutton.renderComponent.TextureR = textureR
	return
}
func (uibutton *UIButton) GetTransform() *common.Transform {
	return uibutton.transform
}
func (uibutton *UIButton) ID_sg(_id ...int) int {
	if len(_id) == 0 {
		return uibutton.id
	}
	if len(_id) > 1 {
		panic("len(_id)")
	}
	uibutton.id = _id[0]
	return uibutton.id
}

func (uibutton *UIButton) GetRenderComponent() *resource.RenderComponent {
	return uibutton.renderComponent
}
func (uibutton *UIButton) Enabled() bool {
	return uibutton.enabled
}

func (uibutton *UIButton) Start() {
}

func (uibutton *UIButton) Update() {
}

func (uibutton *UIButton) OnDraw() {
	// fmt.Println("uibutton,", uibutton.renderComponent.ShaderR)
	uibutton.renderComponent.ShaderR.Active()
	uibutton.renderComponent.ModelR.Active()
	uibutton.renderComponent.TextureR.Active()
	//
	modelMAT := uibutton.transform.Model()
	mloc, lightloc, sortzloc, whrloc := uibutton.shaderOP.UniformLoc("model"),
		uibutton.shaderOP.UniformLoc("light"),
		uibutton.shaderOP.UniformLoc("sortz"),
		uibutton.shaderOP.UniformLoc("whr")
	gl.UniformMatrix4fv(mloc, 1, false, modelMAT.Address())
	gl.Uniform1f(lightloc, 1)
	gl.Uniform1f(sortzloc, uibutton.sortz)
	gl.Uniform1f(whrloc, uibutton.gi.GetWHR())
}

func (uibutton *UIButton) SortZ() float32 {
	return uibutton.sortz
}
func (uibutton *UIButton) OnDrawFinish() {

}
