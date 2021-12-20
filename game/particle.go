package game

import (
	"fmt"
	"math/rand"

	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/gitbufenshuo/gopen/game/common"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type ParticleCore struct {
	modelResource *resource.Model
	Transform     *common.Transform
}

func NewParticleCore() *ParticleCore {
	pc := new(ParticleCore)
	{
		// model : just a quad
		quadModel := resource.NewQuadModel()
		quadModel.Upload()
		pc.modelResource = quadModel
	}
	pc.Transform = common.NewTransform()
	return pc
}

func (pc *ParticleCore) UploadUniforms(ml int32) {
	m := pc.Transform.Model()
	gl.UniformMatrix4fv(ml, 1, false, m.Address())
}

type Particle struct {
	gi                              *GlobalInfo
	ID                              int
	CoreList                        []*ParticleCore
	ShaderResource                  *resource.ShaderProgram
	TextureResource                 *resource.Texture
	MLocation, VLocation, PLocation int32
}

func (parti *Particle) Start() {

}

func (parti *Particle) Update() {
	randint := rand.Int()
	//
	randint %= len(parti.CoreList)
	//
	gi := parti.gi
	parti.CoreList[randint].Transform.Postion.SetIndexValue(1, 3)
	parti.CoreList[randint].Transform.Rotation.SetIndexValue(1, float32(gi.CurFrame))

}

func (parti *Particle) Draw() {

	parti.ShaderResource.Active() // shader
	parti.TextureResource.Active()
	parti.UploadUniforms(parti.VLocation, parti.PLocation)
	for _, onecore := range parti.CoreList {
		// prepare the uniforms
		onecore.UploadUniforms(parti.MLocation)
		// change context
		onecore.modelResource.Active() // model
		// draw
		vertexNum := len(onecore.modelResource.Indices)
		gl.DrawElements(gl.TRIANGLES, int32(vertexNum), gl.UNSIGNED_INT, gl.PtrOffset(0))
	}
}

func (parti *Particle) ID_sg(_id ...int) int {
	if len(_id) == 0 {
		return parti.ID
	}
	if len(_id) > 1 {
		panic("len(_id)")
	}
	parti.ID = _id[0]
	return parti.ID
}
func (parti *Particle) UploadUniforms(vl, pl int32) {
	v := parti.gi.View()
	p := parti.gi.Projection()
	//////////////////////////////////////////////
	gl.UniformMatrix4fv(vl, 1, false, v.Address())
	gl.UniformMatrix4fv(pl, 1, false, p.Address())
}

func NewParticle(gi *GlobalInfo, texture *resource.Texture) *Particle {
	parti := new(Particle)
	parti.gi = gi
	// return blockMan
	{
		// texture
		texture.Upload()
		parti.TextureResource = texture
	}
	{
		// shader
		quadShader := resource.NewShaderProgram()
		quadShader.ReadFromText(resource.ShaderQuadText.Vertex, resource.ShaderQuadText.Fragment)
		quadShader.Upload()
		parti.ShaderResource = quadShader
		//
		fmt.Println("quadShader.ShaderProgram()", quadShader.ShaderProgram())
		parti.MLocation = gl.GetUniformLocation(quadShader.ShaderProgram(), gl.Str("model"+"\x00"))
		parti.VLocation = gl.GetUniformLocation(quadShader.ShaderProgram(), gl.Str("view"+"\x00"))
		parti.PLocation = gl.GetUniformLocation(quadShader.ShaderProgram(), gl.Str("projection"+"\x00"))
	}
	{
		// model list
		for idx := 0; idx != 10; idx++ {
			parti.CoreList = append(parti.CoreList, NewParticleCore())
		}
	}
	return parti
}
