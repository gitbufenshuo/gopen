package game

import (
	"fmt"
	"math"

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
			if buttonModelDefault.Vertices[idx*5+0] < 0 {
				buttonModelDefault.Vertices[idx*5+0] = 0
			} else {
				buttonModelDefault.Vertices[idx*5+0] *= 2
			}
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

type ButtonConfig struct {
	Width      float32
	Height     float32
	Content    string
	ShaderText resource.ShaderText
	TextureR   *resource.Texture
	Bling      bool
	CustomDraw func(shaderOP *ShaderOP)
}

var DefaultButtonConfig = ButtonConfig{
	Width:   1,
	Height:  1,
	Content: "按钮",
}

type UIButton struct {
	gi              *GlobalInfo
	id              int
	renderComponent *resource.RenderComponent
	enabled         bool
	transform       *common.Transform
	shaderOP        *ShaderOP
	bling           bool
	customDraw      func(shaderOP *ShaderOP)

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
	{
		uibutton.renderComponent.ModelR = buttonModelDefault
	}
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
	uibutton.uitext.SetText(DefaultButtonConfig.Content)
	uibutton.uitext.transform.SetParent(uibutton.transform)
	return uibutton
}
func NewCustomButton(gi *GlobalInfo, buttonconfig ButtonConfig) *UIButton {
	InitDefaultButton()
	uibutton := new(UIButton)
	uibutton.sortz = 0.001
	uibutton.bling = buttonconfig.Bling
	if buttonconfig.CustomDraw != nil {
		uibutton.customDraw = buttonconfig.CustomDraw
	}
	uibutton.gi = gi
	//
	uibutton.renderComponent = new(resource.RenderComponent)
	// model config
	{
		uibutton.renderComponent.ModelR = resource.NewQuadModel()
		for idx := 0; idx != 4; idx++ {
			uibutton.renderComponent.ModelR.Vertices[idx*5+0] *= buttonconfig.Width
			uibutton.renderComponent.ModelR.Vertices[idx*5+1] *= buttonconfig.Height
			if uibutton.renderComponent.ModelR.Vertices[idx*5+0] < 0 {
				uibutton.renderComponent.ModelR.Vertices[idx*5+0] = 0
			} else {
				uibutton.renderComponent.ModelR.Vertices[idx*5+0] *= 2
			}
		}
		uibutton.renderComponent.ModelR.Upload()
	}
	// texture config
	if buttonconfig.TextureR == nil {
		uibutton.renderComponent.TextureR = buttonTextureDefault
	} else {
		uibutton.renderComponent.TextureR = buttonconfig.TextureR
	}
	// shader config
	{
		if buttonconfig.ShaderText.Vertex == "" {
			uibutton.renderComponent.ShaderR = buttonShaderDefault
		} else {
			newShaderR := resource.NewShaderProgram()
			newShaderR.ReadFromText(buttonconfig.ShaderText.Vertex, buttonconfig.ShaderText.Fragment)
			newShaderR.Upload()
			uibutton.renderComponent.ShaderR = newShaderR
		}
		uibutton.shaderOP = NewShaderOP()
		uibutton.shaderOP.SetProgram(uibutton.renderComponent.ShaderR.ShaderProgram())
		uibutton.shaderOP.IfUI()
	}
	//
	uibutton.transform = common.NewTransform()
	//
	uibutton.uitext = NewUIText(gi)
	gi.AddUIObject(uibutton.uitext)
	uibutton.uitext.SetText(buttonconfig.Content)
	uibutton.uitext.transform.SetParent(uibutton.transform)
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
func (uibutton *UIButton) AddUniform(name string) {
	// fmt.Println("uibutton,", uibutton.renderComponent.ShaderR)
	uibutton.shaderOP.AddUniform(name)
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
	if uibutton.bling {
		v := math.Sin(float64(uibutton.gi.CurFrame)/0.001) + 1
		v /= 20
		v += 0.8
		gl.Uniform1f(lightloc, float32(v))
	}
	gl.Uniform1f(sortzloc, uibutton.sortz)
	gl.Uniform1f(whrloc, uibutton.gi.GetWHR())
	return
	//
	if uibutton.customDraw != nil {
		uibutton.customDraw(uibutton.shaderOP)
	}
}

func (uibutton *UIButton) SortZ() float32 {
	return uibutton.sortz
}
func (uibutton *UIButton) OnDrawFinish() {

}
