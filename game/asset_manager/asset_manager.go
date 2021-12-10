// Package asset_manager provides asset definition and asset management.

package asset_manager

import (
	"errors"
	"fmt"

	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
)

var (
	ErrNameDup        = errors.New("ErrNameDup")
	ErrTypeNotSupport = errors.New("ErrTypeNotSupport")
)

type AsssetManager struct {
	nowID          int
	assets_by_name map[string]*Asset
	assets_by_id   map[int]*Asset
}

func NewAsssetManager() *AsssetManager {
	var am AsssetManager
	am.assets_by_name = make(map[string]*Asset)
	am.assets_by_id = make(map[int]*Asset)
	return &am
}
func (am *AsssetManager) PrintAllAsset() {
	for name, as := range am.assets_by_name {
		fmt.Println(name, "--", as.Type, as.ID, as.Resource)
	}
}
func (am *AsssetManager) FindByName(name string) *Asset {
	if as, found := am.assets_by_name[name]; found {
		return as
	}
	return nil
}
func (am *AsssetManager) Register(name string, as *Asset) error {
	if _as := am.FindByName(name); _as != nil {
		return ErrNameDup
	}
	am.assets_by_name[name] = as
	return nil
}

func (am *AsssetManager) LoadModelFromFile(filepath, assetname string) {
	var data ModelDataType
	data.FilePath = filepath
	as := NewAsset(assetname, AssetTypeModel, &data)
	err := am.Register(as.Name, as)
	if err != nil {
		panic(err)
	}
	am.Load(as)
	as.Resource.Upload()
}

func (am *AsssetManager) CreateModel(assetname string, modelResource *resource.Model) {
	as := NewAsset(assetname, AssetTypeModel, nil)
	err := am.Register(as.Name, as)
	if err != nil {
		panic(err)
	}
	as.Resource = modelResource
	am.Load(as)
	as.Resource.Upload()
}

func (am *AsssetManager) LoadTextureFromFile(filepath, assetname string) {
	var data TextureDataType
	data.FilePath = filepath
	data.FlipY = true
	as := NewAsset(assetname, AssetTypeTexture, &data)
	err := am.Register(as.Name, as)
	if err != nil {
		panic(err)
	}
	am.Load(as)
	as.Resource.Upload()
}
func (am *AsssetManager) LoadShaderFromFile(vetexPath, fragPath, assetname string) {
	var data ShaderDataType
	data.VPath = vetexPath
	data.FPath = fragPath
	as := NewAsset("mvp_shader", AssetTypeShader, &data)
	err := am.Register(as.Name, as)
	if err != nil {
		panic(err)
	}
	am.Load(as)
	as.Resource.Upload()
}

// will assign id field
func (am *AsssetManager) Load(as *Asset) {
	err := as.Load()
	if err == nil {
		as.ID = am.nowID + 1
		am.assets_by_id[as.ID] = as
		am.nowID++
	}
	return
}
