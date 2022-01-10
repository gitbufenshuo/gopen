package game

import (
	"fmt"

	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/gitbufenshuo/gopen/game/common"
	"github.com/gitbufenshuo/gopen/matmath"
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
	UISpec     UISpec
	RC         *resource.RenderComponent // 如果不为空则这个覆盖ShaderText 和 TextureR
	Content    string
	ShaderText resource.ShaderText
	TextureR   *resource.Texture
	Bling      bool
	SortZ      float32
	CustomDraw func(shaderOP *ShaderOP)
}

var DefaultSpec UISpec = UISpec{
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
	gi              *GlobalInfo
	id              int
	renderComponent *resource.RenderComponent
	enabled         bool
	UISpec          UISpec
	transform       *common.Transform
	shaderOP        *ShaderOP
	bling           bool
	customDraw      func(shaderOP *ShaderOP)
	// a_model_loc     int32
	// u_light_loc     int32
	// u_sortz_loc     int32
	sortz      float32
	mouseHover bool
	//
	uitext *UIText
}

func NewDefaultUIButton(gi *GlobalInfo) *UIButton {
	return NewCustomButton(gi, DefaultButtonConfig)
}
func NewCustomButton(gi *GlobalInfo, buttonconfig ButtonConfig) *UIButton {
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
		uibutton.shaderOP = NewShaderOP()
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
			uibutton.shaderOP = NewShaderOP()
			uibutton.shaderOP.SetProgram(uibutton.renderComponent.ShaderR.ShaderProgram())
			uibutton.shaderOP.IfUI()
		}
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
func (uibutton *UIButton) SetParent(otherTrans *common.Transform) {
	uibutton.transform.SetParent(otherTrans)
}
func (uibutton *UIButton) HoverCheck() bool {
	return true
}
func (uibutton *UIButton) HoverSet(ing bool) {
	uibutton.mouseHover = ing
}

func (uibutton *UIButton) Bounds() []*matmath.Vec2 {
	modelMAT := uibutton.transform.WorldModel()
	projectionMAT := uibutton.gi.UICanvas.Orthographic()
	modelMAT.RightMul_InPlace(&projectionMAT)
	vertices := uibutton.renderComponent.ModelR.Vertices
	bound1 := matmath.CreateVec2FromVec4(
		matmath.CreateVec4(vertices[0], vertices[1], vertices[2], 1).LeftMulMAT(modelMAT),
	)
	bound2 := matmath.CreateVec2FromVec4(
		matmath.CreateVec4(vertices[5], vertices[6], vertices[7], 1).LeftMulMAT(modelMAT),
	)

	bound3 := matmath.CreateVec2FromVec4(
		matmath.CreateVec4(vertices[10], vertices[11], vertices[12], 1).LeftMulMAT(modelMAT),
	)

	bound4 := matmath.CreateVec2FromVec4(
		matmath.CreateVec4(vertices[15], vertices[16], vertices[17], 1).LeftMulMAT(modelMAT),
	)
	// bound1.PrettyShow("1")
	// bound2.PrettyShow(" 2")
	// bound3.PrettyShow("  3")
	// bound4.PrettyShow("   4")
	return []*matmath.Vec2{&bound1, &bound2, &bound3, &bound4}
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

func (uibutton *UIButton) CheckPoint(x, y float32) bool {
	// bounds := uibutton.Bounds()
	//
	return false
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
	uibutton.bling = uibutton.mouseHover
}
func (uibutton *UIButton) AddUniform(name string) {
	// fmt.Println("uibutton,", uibutton.renderComponent.ShaderR)
	uibutton.shaderOP.AddUniform(name)
}
func (uibutton *UIButton) OnDraw() {
	uibutton.renderComponent.ShaderR.Active()
	uibutton.renderComponent.ModelR.Active()
	uibutton.renderComponent.TextureR.Active()
	//
	{
		// 根据 UISpec 得到真正要渲染的参数
		widthDeform := uibutton.gi.GetWindowWidth() / uibutton.gi.UICanvas.DesignWidth
		heightDeform := uibutton.gi.GetWindowHeight() / uibutton.gi.UICanvas.DesignHeight
		// 1. pos
		{
			posx, posy := uibutton.UISpec.LocalPos.GetValue2()
			posrx, posry := uibutton.UISpec.PosRelativity.GetValue2()
			// 根据真实分辨率，计算新的位置
			posxNew := posx * widthDeform
			posyNew := posy * heightDeform
			posxNew = (1-posrx)*posx + posrx*posxNew
			posyNew = (1-posry)*posy + posry*posyNew
			uibutton.transform.Postion.SetIndexValue(0, posxNew)
			uibutton.transform.Postion.SetIndexValue(1, posyNew)
		}
		// 2. scale
		{
			scalex := 1 + uibutton.UISpec.SizeRelativity.GetIndexValue(0)*(widthDeform-1)
			scaley := 1 + uibutton.UISpec.SizeRelativity.GetIndexValue(1)*(heightDeform-1)
			uibutton.transform.Scale.SetValue2(
				scalex,
				scaley,
			)
		}
		// 3. rotate
		{
			// uibutton.transform.Rotation.SetZ(
			// 	float32(uibutton.gi.CurFrame) / 2,
			// )
		}
	}
	modelMAT := uibutton.transform.WorldModel()
	proloc, mloc, lightloc, sortzloc := uibutton.shaderOP.UniformLoc("projection"),
		uibutton.shaderOP.UniformLoc("model"),
		uibutton.shaderOP.UniformLoc("light"),
		uibutton.shaderOP.UniformLoc("sortz")
	orProjection := uibutton.gi.UICanvas.Orthographic()
	gl.UniformMatrix4fv(proloc, 1, false, orProjection.Address())
	gl.UniformMatrix4fv(mloc, 1, false, modelMAT.Address())
	if uibutton.bling {
		v := float32(uibutton.gi.CurFrame % 20)
		v /= 20
		v += 1
		v /= 2
		gl.Uniform1f(lightloc, float32(v))
	} else {
		gl.Uniform1f(lightloc, 1)
	}
	gl.Uniform1f(sortzloc, uibutton.sortz)

	if uibutton.customDraw != nil {
		uibutton.customDraw(uibutton.shaderOP)
	}
}

func (uibutton *UIButton) SortZ() float32 {
	return uibutton.sortz
}
func (uibutton *UIButton) OnDrawFinish() {

}
