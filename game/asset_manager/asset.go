package asset_manager

import (
	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/go-gl/gl/v2.1/gl"
)

type AssetType string

var (
	AssetTypeTexture AssetType = "texture"
)

// Asset is where you can get a resource which is the real thing to use.
// When you get a resource, you don't need a asset.
type Asset struct {
	ID       int
	Loaded   bool
	Name     string
	Type     AssetType
	Data     interface{} // type-dependant
	Resource interface{} // type-dependant
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

func NewAsset(name string, data interface{}) *Asset {
	var as Asset
	as.Name = name
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
	default:
		return ErrTypeNotSupport
	}
	return nil
}
