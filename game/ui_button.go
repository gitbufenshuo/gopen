package game

import "github.com/gitbufenshuo/gopen/game/asset_manager/resource"

var buttonModelDefault *resource.Model
var buttonTextureDefault *resource.Texture
var buttonShaderDefault *resource.ShaderProgram

func init() {
	buttonModelDefault = resource.NewQuadModel()

}

type UIButton struct {
	ModelR   *resource.Model
	TextureR *resource.Texture
	ShaderR  *resource.ShaderProgram
}

func NewDefaultUIButton() *UIButton {
	uibutton := new(UIButton)
	/////////////////////////
	// modelr := resource.NewQuadModel()

	return uibutton
}
