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
	buttonModelDefault = resource.NewQuadModel()
	buttonModelDefault.Upload()
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
	a_model_loc     int32
	u_light_loc     int32
}

func NewDefaultUIButton(gi *GlobalInfo) *UIButton {
	InitDefaultButton()
	uibutton := new(UIButton)
	uibutton.gi = gi
	/////////////////////////
	uibutton.renderComponent = new(resource.RenderComponent)
	uibutton.renderComponent.ModelR = buttonModelDefault
	uibutton.renderComponent.TextureR = buttonTextureDefault
	uibutton.renderComponent.ShaderR = buttonShaderDefault
	//
	uibutton.transform = common.NewTransform()
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
	if uibutton.a_model_loc == 0 {
		uibutton.a_model_loc = gl.GetUniformLocation(uibutton.renderComponent.ShaderR.ShaderProgram(), gl.Str("model"+"\x00"))
		uibutton.u_light_loc = gl.GetUniformLocation(uibutton.renderComponent.ShaderR.ShaderProgram(), gl.Str("light"+"\x00"))
	}
	//
	modelt := uibutton.transform.Model()
	gl.UniformMatrix4fv(uibutton.a_model_loc, 1, false, modelt.Address())
}

func (uibutton *UIButton) OnDrawFinish() {

}
