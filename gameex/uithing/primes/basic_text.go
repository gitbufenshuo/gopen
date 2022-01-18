package primes

import (
	"fmt"
	"math/rand"

	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/gitbufenshuo/gopen/matmath"
	"github.com/go-gl/gl/v4.1-core/gl"
)

const pixration float32 = 0.02

var shaderDefault *resource.ShaderProgram

func STInit() {
	if shaderDefault == nil {
		shaderDefault = resource.NewShaderProgram()
		shaderDefault.ReadFromText(resource.ShaderUIButtonText.Vertex, resource.ShaderUIButtonText.Fragment)
		shaderDefault.Upload()
	}
}

type UIText struct {
	gi              *game.GlobalInfo
	id              int
	renderComponent *resource.RenderComponent
	enabled         bool
	transform       *game.Transform
	shaderOP        *game.ShaderOP
	sortz           float32
	mouseHover      bool
	//
	content string
}

func NewUIText(gi *game.GlobalInfo) *UIText {
	STInit()
	uitext := new(UIText)
	uitext.enabled = true
	uitext.gi = gi
	uitext.sortz = 0.0005
	uitext.id = rand.Int()
	/////////////////////////
	uitext.renderComponent = new(resource.RenderComponent)

	uitext.renderComponent.TextureR = resource.NewTexture()
	{
		uitext.renderComponent.ShaderR = shaderDefault
		uitext.shaderOP = game.NewShaderOP()
		uitext.shaderOP.SetProgram(uitext.renderComponent.ShaderR.ShaderProgram())
		uitext.shaderOP.IfUI()
	}
	//
	uitext.transform = game.NewTransform()
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
func (uitext *UIText) HoverCheck() bool {
	return false
}
func (uitext *UIText) HoverSet(ing bool) {
	uitext.mouseHover = ing
}
func (uitext *UIText) Bounds() []*matmath.Vec2 {
	return nil
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
	uitext.renderComponent.ShaderR.Active()
	uitext.renderComponent.ModelR.Active()
	uitext.renderComponent.TextureR.Active()
	modelMAT := uitext.transform.WorldModel()

	proloc, mloc, lightloc, sortzloc := uitext.shaderOP.UniformLoc("projection"),
		uitext.shaderOP.UniformLoc("model"),
		uitext.shaderOP.UniformLoc("light"),
		uitext.shaderOP.UniformLoc("sortz")
	orProjection := uitext.gi.UICanvas.Orthographic()
	gl.UniformMatrix4fv(proloc, 1, false, orProjection.Address())
	gl.UniformMatrix4fv(mloc, 1, false, modelMAT.Address())
	gl.Uniform1f(lightloc, 1)
	gl.Uniform1f(sortzloc, uitext.sortz)
}

func (uitext *UIText) OnDrawFinish() {

}
func (uitext *UIText) SetParent(parent *game.Transform) {
	uitext.transform.SetParent(parent)
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
	width, height, outterWidth := tr.GenFont(content, uitext.gi.FontConfig)
	tr.Upload()
	if uitext.renderComponent.ModelR != nil {
		uitext.renderComponent.ModelR.Clear()
	}
	modelr := resource.NewUITextModel_BySpec(
		matmath.CreateVec4(-1, 0, 0, 0),
		width, height, outterWidth,
	)
	modelr.Upload()
	uitext.renderComponent.ModelR = modelr
	// re - scale - model
	uitext.transform.Scale.SetIndexValue(0, 0.5)
	uitext.transform.Scale.SetIndexValue(1, 0.5)
}

func (uitext *UIText) SortZ() float32 {
	return uitext.sortz
}
