package asset_manager

import (
	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type AssetType string

var (
	AssetTypeTexture AssetType = "texture"
	AssetTypeShader  AssetType = "shaderprogram"
	AssetTypeModel   AssetType = "model"
	AssetTypeCubeMap AssetType = "cubemap"
)

type Resource interface {
	Upload()
	Active()
}

// Asset is where you can get a resource which is the real thing to use.
// When you get a resource, you don't need a asset.
type Asset struct {
	ID       int
	Loaded   bool
	Name     string
	Type     AssetType
	Data     interface{} // type-dependant
	Resource Resource    // type-dependant
}

type TextureDataType struct {
	FilePath    string
	Width       int32
	Height      int32
	RepeatModeU int
	RepeatModeV int
	FlipY       bool
	GenMipMaps  bool
}
type ShaderDataType struct {
	VPath string // vertex shader script file path
	FPath string // fragment shader script file path
}
type ModelDataType struct {
	FilePath string
}
type CubeMapDataType struct {
	PattList []string
}

func NewAsset(name string, Type AssetType, data interface{}) *Asset {
	var as Asset
	as.Name = name
	as.Type = Type
	as.Data = data
	return &as
}
func (as *Asset) Load() error {
	switch as.Type {
	case AssetTypeTexture:
		_data := as.Data.(*TextureDataType)
		_t := resource.NewTexture()
		_t.SetWidth(_data.Width)
		_t.SetHeight(_data.Height)
		_t.SetFormat(gl.RGBA)
		_t.RepeatModeU = _data.RepeatModeU
		_t.RepeatModeV = _data.RepeatModeV
		_t.GenMipMaps = _data.GenMipMaps
		_t.FlipY = _data.FlipY
		_t.ReadFromFile(_data.FilePath)
		as.Resource = _t
	case AssetTypeShader:
		_data := as.Data.(*ShaderDataType)
		_t := resource.NewShaderProgram()
		_t.ReadFromFile(_data.VPath, _data.FPath)
		as.Resource = _t
	case AssetTypeModel:
		if as.Resource == nil {
			_data := as.Data.(*ModelDataType)
			_t := resource.NewModel()
			_t.ReadFromFile(_data.FilePath)
			as.Resource = _t
		}
	case AssetTypeCubeMap:
		if as.Resource == nil {
			_data := as.Data.(*CubeMapDataType)
			_t := resource.NewCubeMap()
			_t.ReadFromFile(_data.PattList)
			as.Resource = _t
		}
	default:
		return ErrTypeNotSupport
	}
	return nil
}
