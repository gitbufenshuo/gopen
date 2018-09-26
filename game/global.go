package game

import (
	"os"
	"path"

	"github.com/gitbufenshuo/gopen/game/asset_manager"
)

type GlobalInfo struct {
	AssetManager     *asset_manager.AsssetManager
	DefaultAssetList map[string]*asset_manager.Asset
}

func NewGlobalInfo() *GlobalInfo {
	globalInfo := new(GlobalInfo)
	return globalInfo
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
	as := asset_manager.NewAsset("triangle", &data)
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
	as := asset_manager.NewAsset("minimal_shader", &data)
	err := gi.AssetManager.Register(as.Name, as)
	if err != nil {
		panic(err)
	}
	gi.AssetManager.Load(as)
}
