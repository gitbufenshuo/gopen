package pk_basic_button

import (
	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/gitbufenshuo/gopen/matmath"
	"github.com/go-gl/gl/v4.1-core/gl"
)

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
func (uibutton *UIButton) GetTransform() *game.Transform {
	return uibutton.transform
}
func (uibutton *UIButton) SetParent(otherTrans *game.Transform) {
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
