package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/gitbufenshuo/gopen/help"
	"github.com/gitbufenshuo/gopen/matmath"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"

	"github.com/gitbufenshuo/gopen/game/asset_manager"
	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

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
	AssetManager          *asset_manager.AsssetManager
	gameobjects           map[int]GameObjectI
	manageobjects         map[int]ManageObjectI
	uiobjects             map[int]UIObjectI
	nowID                 int
	nowMD                 int
	nowUD                 int
	width                 int
	height                int
	title                 string
	FontConfig            *help.FontConfig
	CustomInit            func(*GlobalInfo)
	MainCamera            *Camera
	ParticalSystem        *Particle
	window                *glfw.Window
	InputSystemKeyPress   []bool
	InputSystemKeyRelease []bool
	keyCallback           func(win *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey)
	InputMouseCtl         *InputMouse
	CursorMode            int
	InputSystemManager    InputSystemI
	//
	MouseXDiff, MouseYDiff float64
	*GlobalFrameInfo
}

func NewGlobalInfo(windowWidth, windowHeight int, title string) *GlobalInfo {
	globalInfo := new(GlobalInfo)
	globalInfo.width = windowWidth
	globalInfo.height = windowHeight
	globalInfo.title = title
	globalInfo.gameobjects = make(map[int]GameObjectI)
	globalInfo.manageobjects = make(map[int]ManageObjectI)
	globalInfo.uiobjects = make(map[int]UIObjectI)
	globalInfo.CursorMode = glfw.CursorNormal
	return globalInfo
}

func (gi *GlobalInfo) Window() *glfw.Window {
	return gi.window
}

func (gi *GlobalInfo) GetWHR() float32 {
	return float32(gi.width) / float32(gi.height)
}

func (gi *GlobalInfo) LoadFont(fontpath string) {
	fontBytes, err := ioutil.ReadFile(fontpath)
	if err != nil {
		log.Println(err)
		return
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}

	// Truetype stuff
	opts := truetype.Options{}
	opts.Size = 30
	face := truetype.NewFace(font, &opts)
	gi.FontConfig = help.NewFontConfig(font, face)
}
func (gi *GlobalInfo) StartGame(mode string) {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()
	rand.Seed(help.GetNowMS())
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(gi.width, gi.height, gi.title, nil, nil)
	if err != nil {
		panic(err)
	}
	gi.window = window
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
	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.BLEND)
	gl.DepthFunc(gl.LESS)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.ClearColor(1, 1, 1, 1)
	{
		// interact
		if gi.keyCallback != nil {
			window.SetKeyCallback(gi.keyCallback)
		}
		gi.InputMouseCtl = NewInputMouse()
		window.SetCursorPosCallback(gi.InputMouseCtl.CursorCallback)
		window.SetMouseButtonCallback(gi.InputMouseCtl.MouseButtonCallback)
	}
	gi.startlogic()
	// start hook
	for !window.ShouldClose() {
		// time.Sleep(time.Millisecond * 10)
		gl.ClearColor(0.5, 0.5, 0.5, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		// window.SwapBuffers()
		time.Sleep(time.Millisecond * 20)
		///////////////////////////////////////////////////
		// the very update every frame
		gi.update()
		gi.draw()
		gi.OnFrameEnd()
		///////////////////////////////////////////////////
		// Maintenance
		window.SwapBuffers()
		// time.Sleep(time.Second)
		glfw.PollEvents()
		frame_number++
	}

}

func (gi *GlobalInfo) startlogic() {
	if gi.InputSystemManager != nil {
		gi.InputSystemManager.Start()
	}
	for _, mb := range gi.manageobjects {
		mb.Start()
	}
	for _, gb := range gi.gameobjects {
		gb.Start()
	}
	for _, ub := range gi.uiobjects {
		ub.Start()
	}
}

func (gi *GlobalInfo) SetInputSystem(is InputSystemI) {
	gi.InputSystemManager = is
}

func (gi *GlobalInfo) OnFrameEnd() {
	gi.MouseXDiff = 0
	gi.MouseYDiff = 0
}

func (gi *GlobalInfo) SetCursorMode(mode int) {
	gi.CursorMode = mode
	gi.window.SetInputMode(glfw.CursorMode, gi.CursorMode)
	fmt.Println(":asdfasdfasdf")
}

func (gi *GlobalInfo) SetKeyCallback(callback func(win *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey)) {
	gi.keyCallback = callback
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
	// gi.dealWithTime(0)
	if gi.InputSystemManager != nil {
		gi.InputSystemManager.Update()
	}
	for _, gb := range gi.gameobjects {
		gb.Update() // call the gameobjects' Update function
	}
	for _, mb := range gi.manageobjects {
		mb.Update() // call the manageobjects' Update function
	}
	for _, ub := range gi.uiobjects {
		ub.Update()
	}
	gi.dealWithTime(2)
}
func (gi *GlobalInfo) draw() {
	// gi.dealWithTime(0)
	for _, gb := range gi.gameobjects {
		gi.drawGameobject(gb)
	}
	gl.DepthFunc(gl.LEQUAL)
	// draw the skybox
	if gi.MainCamera.CubeMapObject != nil {
		gi.drawSkyBox()
	}
	gl.DepthFunc(gl.LESS)
	// particle system
	if gi.ParticalSystem != nil {
		gi.ParticalSystem.Update()
		gi.ParticalSystem.Draw()
	}
	// ui system
	{
		var aluiob []UIObjectI = make([]UIObjectI, len(gi.uiobjects))
		var appendidx int
		for _, ub := range gi.uiobjects {
			aluiob[appendidx] = ub
			appendidx++
		}
		sort.Slice(aluiob, func(i, j int) bool {
			return aluiob[i].SortZ() > aluiob[j].SortZ()
		})
		for _, ub := range aluiob {
			if !ub.Enabled() {
				continue
			}
			ub.OnDraw()
			rc := ub.GetRenderComponent()
			vertexNum := len(rc.ModelR.Indices)
			gl.DrawElements(gl.TRIANGLES, int32(vertexNum), gl.UNSIGNED_INT, gl.PtrOffset(0))
		}
	}
}

func (gi *GlobalInfo) prepareMVP(co GameObjectI) {
	co.SetUniform()
	// co.ShaderCtl().M = co.GetTransform().Model()
	// co.ShaderCtl().Rotation = co.GetTransform().RotationMAT4()
	// var curTransform *common.Transform
	// curTransform = co.GetTransform()
	// for {
	// 	if curTransform.Parent != nil { // not root
	// 		parentM := curTransform.Parent.Model()
	// 		co.ShaderCtl().M.RightMul_InPlace(&parentM)
	// 		parentR := curTransform.Parent.RotationMAT4()
	// 		co.ShaderCtl().Rotation.RightMul_InPlace(&parentR)
	// 	} else {
	// 		break
	// 	}
	// 	curTransform = curTransform.Parent
	// }
	// co.ShaderCtl().V = gi.View()
	// co.ShaderCtl().P = gi.Projection()

	// co.ShaderCtl().Upload(co)
}

func (gi *GlobalInfo) drawGameobject(gb GameObjectI) {
	if gb.NotDrawable() {
		return
	}
	if !gb.DrawEnable_sg() {
		return
	}
	gb.ShaderAsset_sg().Resource.Active() // shader
	gi.prepareMVP(gb)
	gb.OnDraw() // call the gameobjects' OnDraw function
	// change context
	gb.ModelAsset_sg().Resource.Active() // model
	if _asset := gb.TextureAsset_sg(); _asset != nil {
		_asset.Resource.Active()
	}
	// draw
	modelResource := gb.ModelAsset_sg().Resource.(*resource.Model)
	vertexNum := len(modelResource.Indices)
	gl.DrawElements(gl.TRIANGLES, int32(vertexNum), gl.UNSIGNED_INT, gl.PtrOffset(0))
	gb.OnDrawFinish()
}
func (gi *GlobalInfo) drawSkyBox() {
	var rotation = gi.MainCamera.Transform.RotationMAT4()
	gi.MainCamera.CubeMapObject.shaderResource.Active()
	gl.UniformMatrix4fv(gi.MainCamera.CubeMapObject.RotationLocation, 1, false, rotation.Address())
	// change context
	gi.MainCamera.CubeMapObject.modelResource.Active()
	gi.MainCamera.CubeMapObject.cubemapResource.Active()
	// draw
	vertexNum := len(gi.MainCamera.CubeMapObject.modelResource.Indices)
	gl.DrawElements(gl.TRIANGLES, int32(vertexNum), gl.UNSIGNED_INT, gl.PtrOffset(0))
}

func (gi *GlobalInfo) AddGameObject(gb GameObjectI) {
	gb.ID_sg(gi.nowID + 1)

	gi.nowID++
	gi.gameobjects[gb.ID_sg()] = gb
}
func (gi *GlobalInfo) AddManageObject(mb ManageObjectI) {
	mb.ID_sg(gi.nowMD + 1)

	gi.nowMD++
	gi.manageobjects[mb.ID_sg()] = mb
}
func (gi *GlobalInfo) AddUIObject(ub UIObjectI) {
	ub.ID_sg(gi.nowUD + 1)

	gi.nowUD++
	gi.uiobjects[ub.ID_sg()] = ub
}

// init assetmanager and some default assets
func (gi *GlobalInfo) initAssetManager() {
	gi.AssetManager = asset_manager.NewAsssetManager()
	return
	// default model
	gi.initDefaultModel_Triangle()
	// default shader program
	gi.initDefaultShaderprogram_minimal()
	// default texture
	gi.initDefaultTexture_logo()
}
func (gi *GlobalInfo) initDefaultModel_Triangle() {
	gi.AssetManager.LoadModelFromFile(path.Join(os.Getenv("HOME"), ".gopen", "assets", "models", "triangle.json"), "triangle")
}
func (gi *GlobalInfo) initDefaultShaderprogram_minimal() {
	gi.AssetManager.LoadShaderFromFile(path.Join(os.Getenv("HOME"), ".gopen", "assets", "shaderprograms", "minimal_vertex.glsl"),
		path.Join(os.Getenv("HOME"), ".gopen", "assets", "shaderprograms", "minimal_fragment.glsl"), "minimal_shader",
	)
}
func (gi *GlobalInfo) initDefaultTexture_logo() {
	gi.AssetManager.LoadTextureFromFile(path.Join(os.Getenv("HOME"), ".gopen", "assets", "textures", "logo.png"), "logo_texture")
}

func (gi *GlobalInfo) View() matmath.MATX {
	return gi.MainCamera.ViewMat()
	// viewT := matmath.LookAtFrom4(&gi.MainCamera.Pos, &gi.MainCamera.Target, &gi.MainCamera.UP)
	// gi.MainCamera.ViewT = viewT
	// return viewT
}

func (gi *GlobalInfo) Projection() matmath.MATX {
	projectionT := matmath.Perspective(gi.MainCamera.NearDistance, gi.MainCamera.FarDistance, gi.MainCamera.FOV)
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
