package pk_basic_button

import (
	"fmt"

	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/gitbufenshuo/gopen/gameex/uithing/primes"
	"github.com/gitbufenshuo/gopen/matmath"
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
		buttonModelDefault = resource.NewQuadModel_LeftALign()
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
	UISpec     game.UISpec
	RC         *resource.RenderComponent // 如果不为空则这个覆盖ShaderText 和 TextureR
	Content    string
	ShaderText resource.ShaderText
	TextureR   *resource.Texture
	Bling      bool
	SortZ      float32
	CustomDraw func(shaderOP *game.ShaderOP)
}

var DefaultSpec game.UISpec = game.UISpec{
	Pivot:    matmath.CreateVec4(0, 0, 0, 0),
	LocalPos: matmath.CreateVec4(0, 0, 0, 0),
	Width:    100,
	Height:   30,
}
var DefaultButtonConfig = ButtonConfig{
	UISpec:  DefaultSpec,
	SortZ:   0.001,
	Content: "按钮",
}

type UIButton struct {
	gi              *game.GlobalInfo
	id              int
	renderComponent *resource.RenderComponent
	enabled         bool
	UISpec          game.UISpec
	transform       *game.Transform
	shaderOP        *game.ShaderOP
	bling           bool
	customDraw      func(shaderOP *game.ShaderOP)
	// a_model_loc     int32
	// u_light_loc     int32
	// u_sortz_loc     int32
	sortz      float32
	mouseHover bool
	//
	uitext *primes.UIText
}

func NewDefaultUIButton(gi *game.GlobalInfo) *UIButton {
	return NewCustomButton(gi, DefaultButtonConfig)
}
func NewCustomButton(gi *game.GlobalInfo, buttonconfig ButtonConfig) *UIButton {
	InitDefaultButton()
	uibutton := new(UIButton)
	uibutton.enabled = true
	uibutton.UISpec = buttonconfig.UISpec
	if buttonconfig.SortZ > 0 {
		uibutton.sortz = buttonconfig.SortZ
	}
	uibutton.bling = buttonconfig.Bling
	if buttonconfig.CustomDraw != nil {
		uibutton.customDraw = buttonconfig.CustomDraw
	}
	uibutton.gi = gi
	//
	if buttonconfig.RC != nil {
		uibutton.renderComponent = buttonconfig.RC
		uibutton.shaderOP = game.NewShaderOP()
		uibutton.shaderOP.SetProgram(uibutton.renderComponent.ShaderR.ShaderProgram())
		uibutton.shaderOP.IfUI()
	} else {
		uibutton.renderComponent = new(resource.RenderComponent)
		// model config
		{
			uibutton.renderComponent.ModelR = resource.NewQuadModel_BySpec(uibutton.UISpec.Pivot, uibutton.UISpec.Width, uibutton.UISpec.Height)
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
			uibutton.shaderOP = game.NewShaderOP()
			uibutton.shaderOP.SetProgram(uibutton.renderComponent.ShaderR.ShaderProgram())
			uibutton.shaderOP.IfUI()
		}
	}

	//
	uibutton.transform = game.NewTransform()
	//
	uibutton.uitext = primes.NewUIText(gi)
	gi.AddUIObject(uibutton.uitext)
	uibutton.uitext.SetText(buttonconfig.Content)
	uibutton.uitext.SetParent(uibutton.transform)
	return uibutton
}
func (uibutton *UIButton) ChangeTexture(textureR *resource.Texture) {
	uibutton.renderComponent.TextureR = textureR
	return
}

func (uibutton *UIButton) Disable() {
	uibutton.enabled = false
	if uibutton.uitext != nil {
		uibutton.uitext.Disable()
	}
}

func (uibutton *UIButton) Enable() {
	uibutton.enabled = true
	if uibutton.uitext != nil {
		uibutton.uitext.Enable()
	}
}

// 切换闪烁
func (uibutton *UIButton) SwitchBling() bool {
	uibutton.bling = !uibutton.bling
	return uibutton.bling
}

// 强制闪烁
func (uibutton *UIButton) EnableBling() {
	uibutton.bling = true
}

// 禁用闪烁
func (uibutton *UIButton) DisableBling() {
	uibutton.bling = false
}

func (uibutton *UIButton) AddUniform(name string) {
	// fmt.Println("uibutton,", uibutton.renderComponent.ShaderR)
	uibutton.shaderOP.AddUniform(name)
}
