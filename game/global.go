package game

import (
	"fmt"
	"log"
	"math"
	"os"
	"path"
	"time"

	"github.com/gitbufenshuo/gopen/game/asset_manager"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

type GlobalInfo struct {
	AssetManager *asset_manager.AsssetManager
	width        int
	height       int
	title        string
}

func NewGlobalInfo(windowWidth, windowHeight int, title string) *GlobalInfo {
	globalInfo := new(GlobalInfo)
	globalInfo.width = windowWidth
	globalInfo.height = windowHeight
	globalInfo.title = title

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
	var r float32
	var frame_number int
	gi.initAssetManager()
	gi.AssetManager.PrintAllAsset()
	if mode == "test" {
		return
	}
	for !window.ShouldClose() {
		time.Sleep(time.Millisecond * 30)
		r = float32(math.Sin(math.Pi*float64((frame_number*2)%1000)/500))/2 + 0.5
		fmt.Println(r)
		gl.ClearColor(r, 0, 0, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
		frame_number++
	}

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
