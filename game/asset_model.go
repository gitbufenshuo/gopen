package game

import (
	"sync"

	"github.com/gitbufenshuo/gopen/matmath"
)

var all_model_pool *sync.Pool

var all_model_asset map[int]*AssetModel

func init() {
	var NewAssetModel = func() interface{} {
		var assetModel AssetModel
		return &assetModel
	}
	all_model_pool = &sync.Pool{
		New: NewAssetModel,
	}
	all_model_asset = make(map[int]*AssetModel)
}

type Vertex struct {
	Position *matmath.VECX // 3
	UV       *matmath.VECX // 2
	Color    *matmath.VECX // 4
}

// Model includes :
// 1. vertices
// 2. indices
type AssetModel struct {
	AssetID  int // gameobject should reference this id
	Vertices []*Vertex
	Indices  []int32
}
