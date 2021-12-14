package stchunk

import (
	"fmt"
	"math/rand"

	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/gitbufenshuo/gopen/game/gameobjects"
)

type Chunk struct {
	*gameobjects.BlockObject
	InnerModel *resource.Model
}

func NewChunk(gi *game.GlobalInfo) *Chunk {
	block := gameobjects.NewBlock(gi, "block.model", "grid.png.texuture")
	block.Color = []float32{1, 1, 1}
	chunk := new(Chunk)
	///
	chunk.BlockObject = block
	chunk.InnerModel = chunk.ModelAsset_sg().Resource.(*resource.Model)
	chunk.InnerModel.CopyFrom(resource.NewBlockModel())
	return chunk
}
func (chunk *Chunk) Update() {
	if chunk.GI().CurFrame%10 == 0 {
		fmt.Println("chunk update")
		//
		for idx := 0; idx != 24; idx++ {
			chunk.InnerModel.Vertices[idx*8+3] = rand.Float32()
			chunk.InnerModel.Vertices[idx*8+4] = rand.Float32()
		}
		chunk.InnerModel.Clear()
		chunk.InnerModel.Upload()
	}
}
