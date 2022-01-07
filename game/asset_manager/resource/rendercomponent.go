package resource

type RenderComponent struct {
	ModelR   *Model
	TextureR *Texture
	ShaderR  *ShaderProgram
}

func NewRenderComponent(model *Model, texture *Texture, shader *ShaderProgram) *RenderComponent {
	res := new(RenderComponent)
	//
	res.ModelR = model
	res.TextureR = texture
	res.ShaderR = shader
	return res
}
