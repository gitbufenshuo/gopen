package game

import (
	"fmt"

	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/gitbufenshuo/gopen/game/common"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type UIText struct {
	gi              *GlobalInfo
	id              int
	renderComponent *resource.RenderComponent
	enabled         bool
	transform       *common.Transform
	a_model_loc     int32
	u_light_loc     int32
	u_sortz_loc     int32
	sortz           float32
	//
	content string
}

func NewUIText(gi *GlobalInfo) *UIText {
	InitDefaultButton()
	uitext := new(UIText)
	uitext.gi = gi
	uitext.sortz = 0.0005
	/////////////////////////
	uitext.renderComponent = new(resource.RenderComponent)
	uitext.renderComponent.ModelR = buttonModelDefault
	uitext.renderComponent.TextureR = resource.NewTexture()
	uitext.renderComponent.ShaderR = buttonShaderDefault
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

func (uitext *UIText) Start() {
}

func (uitext *UIText) Update() {
	if uitext.gi.CurFrame%100 == 0 {
		fmt.Println("func (uitext *UIText) Update() {")
		if len(uitext.content) <= 30 {
			uitext.SetText(uitext.content + "0 ")
		} else {
			uitext.SetText("0")
		}
	}
}

func (uitext *UIText) OnDraw() {
	// fmt.Println("uibutton,", uibutton.renderComponent.ShaderR)
	uitext.renderComponent.ShaderR.Active()
	uitext.renderComponent.ModelR.Active()
	uitext.renderComponent.TextureR.Active()
	if uitext.a_model_loc == 0 {
		uitext.a_model_loc = gl.GetUniformLocation(uitext.renderComponent.ShaderR.ShaderProgram(), gl.Str("model"+"\x00"))
		uitext.u_light_loc = gl.GetUniformLocation(uitext.renderComponent.ShaderR.ShaderProgram(), gl.Str("light"+"\x00"))
		uitext.u_sortz_loc = gl.GetUniformLocation(uitext.renderComponent.ShaderR.ShaderProgram(), gl.Str("sortz"+"\x00"))
	}
	//
	modelt := uitext.transform.WorldModel()
	gl.UniformMatrix4fv(uitext.a_model_loc, 1, false, modelt.Address())
	gl.Uniform1f(uitext.u_light_loc, 1)
	gl.Uniform1f(uitext.u_sortz_loc, 0.0005)
}

func (uitext *UIText) OnDrawFinish() {

}

//
func (uitext *UIText) SetText(content string) {
	if uitext.content == content {
		return
	}
	//
	uitext.content = content
	// re - gen - texture
	tr := uitext.renderComponent.TextureR
	tr.Clear()
	pixWidth := tr.GenFont(content, uitext.gi.FontConfig)
	tr.Upload()
	// re - scale - model
	uitext.transform.Scale.SetIndexValue(0, float32(pixWidth/16)/3)
	uitext.transform.Scale.SetIndexValue(1, 0.3)
}

func (uitext *UIText) SortZ() float32 {
	return uitext.sortz
}
