package game

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/gitbufenshuo/gopen/game/asset_manager"
	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

type GlobalInfo struct {
	AssetManager *asset_manager.AsssetManager
	gameobjects  map[int]*GameObject
	nowID        int
	width        int
	height       int
	title        string
}

func NewGlobalInfo(windowWidth, windowHeight int, title string) *GlobalInfo {
	globalInfo := new(GlobalInfo)
	globalInfo.width = windowWidth
	globalInfo.height = windowHeight
	globalInfo.title = title
	globalInfo.gameobjects = make(map[int]*GameObject)
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
	gi.initAssetManager()
	if mode == "test_init" {
		gi.AssetManager.PrintAllAsset()
		return
	}
	one := NewGameObject()
	one.ModelAsset = gi.AssetManager.FindByName("triangle")
	one.ShaderAsset = gi.AssetManager.FindByName("minimal_shader")
	gi.AddGameObject(one)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(1, 1, 1, 1)
	for !window.ShouldClose() {
		time.Sleep(time.Millisecond * 30)
		// r = float32(math.Sin(math.Pi*float64((frame_number*2)%1000)/500))/2 + 0.5
		// fmt.Println(r)
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
func (gi *GlobalInfo) update() {
	for _, gb := range gi.gameobjects {
		fmt.Println("update", time.Now().Unix())
		gi.draw(gb)
	}
}
func (gi *GlobalInfo) draw(gb *GameObject) {
	if !gb.readyForDraw {
		// set something
		fmt.Println("set something")
		gb.ShaderAsset.Resource.Upload()
		gb.ModelAsset.Resource.Upload()
		gb.readyForDraw = true
	}
	// change context
	gb.ShaderAsset.Resource.Active()
	gb.ModelAsset.Resource.Active()
	// draw
	modelResource := gb.ModelAsset.Resource.(*resource.Model)
	// vertexNum := len(modelResource.Indices)
	fmt.Println(modelResource.Indices)
	fmt.Println(modelResource.Vertices)
	gl.DrawElements(gl.TRIANGLES, 3, gl.UNSIGNED_INT, gl.PtrOffset(0))
	// gl.DrawArrays(gl.TRIANGLES, 0, 3)
}
func (gi *GlobalInfo) AddGameObject(gb *GameObject) {
	gb.ID = gi.nowID + 1
	gi.nowID++
	gi.gameobjects[gb.ID] = gb
}

// init assetmanager and some default assets
func (gi *GlobalInfo) initAssetManager() {
	gi.AssetManager = asset_manager.NewAsssetManager()
	// default model
	gi.initDefaultModel_Triangle()
	// default shader program
	gi.initDefaultShaderprogram_minimal()
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
