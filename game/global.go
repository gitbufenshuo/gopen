package game

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gitbufenshuo/gopen/matmath"

	"github.com/gitbufenshuo/gopen/game/asset_manager"
	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

type Camera struct {
	Pos          *matmath.VECX
	Front        *matmath.VECX
	UP           *matmath.VECX
	Target       *matmath.VECX
	NearDistance float32
	ViewT        *matmath.MATX
	ProjectionT  *matmath.MATX
}
type GlobalFrameInfo struct {
	CurFrame       int
	StartMS        float64 // the time that the globalinfo successfully starts at
	LastFrameMS    float64 // the time that the last frame begins at
	NowMS          float64 // the time that the current frame begins at
	ElapsedMS      float64 // the ms between now and the StartMS
	FrameElapsedMS float64 // the ms between now and the last frame
	FrameRate      float64 // the frame per second
	Debug          bool    // whether print the frame info
}
type GlobalInfo struct {
	AssetManager *asset_manager.AsssetManager
	gameobjects  map[int]GameObjectI
	nowID        int
	width        int
	height       int
	title        string
	CustomInit   func(*GlobalInfo)
	MainCamera   *Camera
	*GlobalFrameInfo
}

func NewGlobalInfo(windowWidth, windowHeight int, title string) *GlobalInfo {
	globalInfo := new(GlobalInfo)
	globalInfo.width = windowWidth
	globalInfo.height = windowHeight
	globalInfo.title = title
	globalInfo.gameobjects = make(map[int]GameObjectI)
	return globalInfo
}
func (gi *GlobalInfo) StartGame(mode string) {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(gi.width, gi.height, gi.title, nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)
	var frame_number int
	gi.Boot()
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(1, 1, 1, 1)
	for !window.ShouldClose() {
		time.Sleep(time.Millisecond * 20)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		///////////////////////////////////////////////////
		// the very update every frame
		gi.update()
		///////////////////////////////////////////////////
		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
		frame_number++
	}

}
func (gi *GlobalInfo) Boot() {
	gi.initAssetManager()
	if gi.CustomInit == nil {
		panic("gi.CustomInit == nil")
	}
	gi.CustomInit(gi)
	if gi.MainCamera == nil {
		panic("MainCamera == nil")
	}
	gi.GlobalFrameInfo = new(GlobalFrameInfo)
	gi.StartMS = float64(time.Now().Unix()*1000 + int64(time.Now().Nanosecond()/1000000))
	for _, gb := range gi.gameobjects {
		gb.Start()
	}
}
func (gi *GlobalInfo) dealWithTime(mode int) {
	if mode == 1 { // new frame begins
		gi.NowMS = float64(time.Now().Unix()*1000 + int64(time.Now().Nanosecond()/1000000))
		gi.ElapsedMS = gi.NowMS - gi.StartMS
		gi.FrameElapsedMS = gi.NowMS - gi.LastFrameMS
		gi.FrameRate = 1000 / gi.FrameElapsedMS
	}
	if mode == 2 { // frame ends
		gi.LastFrameMS = gi.NowMS
	}
	if mode == 0 { // only for print
		if gi.GlobalFrameInfo.Debug {
			info, _ := json.Marshal(gi.GlobalFrameInfo)
			fmt.Println(string(info))
		}
	}
}
func (gi *GlobalInfo) update() {
	gi.CurFrame++
	gi.dealWithTime(1)
	gi.dealWithTime(0)
	for _, gb := range gi.gameobjects {
		gi.draw(gb)
	}
	gi.dealWithTime(2)
}
func (gi *GlobalInfo) draw(gb GameObjectI) {
	if gb.NotDrawable() {
		return
	}
	if !gb.DrawEnable_sg() {
		return
	}
	if !gb.ReadyForDraw_sg() {
		// set something
		gb.ShaderAsset_sg().Resource.Upload()
		gb.ModelAsset_sg().Resource.Upload()
		if _asset := gb.TextureAsset_sg(); _asset != nil {
			_asset.Resource.Upload()
		}

		gb.ReadyForDraw_sg(true)
	}
	gb.Update() // call the gameobjects' Update function
	gb.OnDraw() // call the gameobjects' OnDraw function
	// change context
	gb.ShaderAsset_sg().Resource.Active() // shader
	gb.ModelAsset_sg().Resource.Active()  // model
	if _asset := gb.TextureAsset_sg(); _asset != nil {
		_asset.Resource.Active()
	}
	// draw
	modelResource := gb.ModelAsset_sg().Resource.(*resource.Model)
	vertexNum := len(modelResource.Indices)
	gl.DrawElements(gl.TRIANGLES, int32(vertexNum), gl.UNSIGNED_INT, gl.PtrOffset(0))
	// refresh something that vary every frame
	matmath.DontNeedMATXAnyMore(gi.MainCamera.ViewT)
	matmath.DontNeedMATXAnyMore(gi.MainCamera.ProjectionT)
	gi.MainCamera.ViewT = nil
	gi.MainCamera.ProjectionT = nil
}
func (gi *GlobalInfo) AddGameObject(gb GameObjectI) {
	gb.ID_sg(gi.nowID + 1)

	gi.nowID++
	gi.gameobjects[gb.ID_sg()] = gb
}

// init assetmanager and some default assets
func (gi *GlobalInfo) initAssetManager() {
	gi.AssetManager = asset_manager.NewAsssetManager()
	// default model
	gi.initDefaultModel_Triangle()
	// default shader program
	gi.initDefaultShaderprogram_minimal()
	// default texture
	gi.initDefaultTexture_logo()
}
func (gi *GlobalInfo) initDefaultModel_Triangle() {
	var data asset_manager.ModelDataType
	data.FilePath = path.Join(os.Getenv("HOME"), ".gopen", "assets", "models", "triangle.json")
	as := asset_manager.NewAsset("triangle", asset_manager.AssetTypeModel, &data)
	err := gi.AssetManager.Register(as.Name, as)
	if err != nil {
		panic(err)
	}
	gi.AssetManager.Load(as)
}
func (gi *GlobalInfo) initDefaultShaderprogram_minimal() {
	var data asset_manager.ShaderDataType
	data.VPath = path.Join(os.Getenv("HOME"), ".gopen", "assets", "shaderprograms", "minimal_vertex.glsl")
	data.FPath = path.Join(os.Getenv("HOME"), ".gopen", "assets", "shaderprograms", "minimal_fragment.glsl")
	as := asset_manager.NewAsset("minimal_shader", asset_manager.AssetTypeShader, &data)
	err := gi.AssetManager.Register(as.Name, as)
	if err != nil {
		panic(err)
	}
	gi.AssetManager.Load(as)
}
func (gi *GlobalInfo) initDefaultTexture_logo() {
	var data asset_manager.TextureDataType
	data.FilePath = path.Join(os.Getenv("HOME"), ".gopen", "assets", "textures", "logo.png")
	data.FlipY = true
	as := asset_manager.NewAsset("logo_texture", asset_manager.AssetTypeTexture, &data)
	err := gi.AssetManager.Register(as.Name, as)
	if err != nil {
		panic(err)
	}
	gi.AssetManager.Load(as)
}

func (gi *GlobalInfo) View() *matmath.MATX {
	if gi.MainCamera.ViewT != nil {
		return gi.MainCamera.ViewT
	}
	gi.MainCamera.Target = gi.MainCamera.Pos.Add(gi.MainCamera.Front)
	viewT := matmath.LookAtFrom4(gi.MainCamera.Pos, gi.MainCamera.Target, gi.MainCamera.UP)
	gi.MainCamera.ViewT = viewT
	return viewT
}
func (gi *GlobalInfo) Projection() *matmath.MATX {
	if gi.MainCamera.ProjectionT != nil {
		return gi.MainCamera.ProjectionT
	}
	projectionT := matmath.Homoz4(gi.MainCamera.NearDistance)
	gi.MainCamera.ProjectionT = projectionT
	return projectionT
}

//// fortestonly
var vao uint32
var vbo uint32

// var program uint32

func fortestonly(mode string) {
	if mode == "init" {
		// Configure the vertex and fragment shaders
		// program, _ = newProgram(vertexShader, fragmentShader)

		gl.GenVertexArrays(1, &vao)
		gl.BindVertexArray(vao)

		gl.GenBuffers(1, &vbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, gl.Ptr(cubeVertices), gl.STATIC_DRAW)

		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	}
	if mode == "update" {
		// gl.UseProgram(program)

		gl.BindVertexArray(vao)

		gl.DrawArrays(gl.TRIANGLES, 0, 3)

	}
	// Configure the vertex data
}

func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

var vertexShader = `
#version 330

in vec3 vert;

void main() {
    gl_Position = vec4(vert.xyz, 1);
}
` + "\x00"

var fragmentShader = `
#version 330

out vec4 outputColor;

void main() {
    outputColor = vec4(1,0,0,1);
}
` + "\x00"

var cubeVertices = []float32{
	//  X, Y, Z, U, V
	// Bottom
	-1.0, -1.0, -1.0, 0.0, 0.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	0, 1.0, 1.0, 0.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, -1.0, 1.0, 1.0, 1.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,

	// Top
	-1.0, 1.0, -1.0, 0.0, 0.0,
	-1.0, 1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, -1.0, 1.0, 0.0,
	1.0, 1.0, -1.0, 1.0, 0.0,
	-1.0, 1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, 1.0, 1.0, 1.0,

	// Front
	-1.0, -1.0, 1.0, 1.0, 0.0,
	1.0, -1.0, 1.0, 0.0, 0.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,
	1.0, -1.0, 1.0, 0.0, 0.0,
	1.0, 1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,

	// Back
	-1.0, -1.0, -1.0, 0.0, 0.0,
	-1.0, 1.0, -1.0, 0.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	-1.0, 1.0, -1.0, 0.0, 1.0,
	1.0, 1.0, -1.0, 1.0, 1.0,

	// Left
	-1.0, -1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, -1.0, 1.0, 0.0,
	-1.0, -1.0, -1.0, 0.0, 0.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,
	-1.0, 1.0, -1.0, 1.0, 0.0,

	// Right
	1.0, -1.0, 1.0, 1.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, 1.0, -1.0, 0.0, 0.0,
	1.0, -1.0, 1.0, 1.0, 1.0,
	1.0, 1.0, -1.0, 0.0, 0.0,
	1.0, 1.0, 1.0, 0.0, 1.0,
}
