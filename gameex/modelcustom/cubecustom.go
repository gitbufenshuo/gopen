package modelcustom

import "github.com/gitbufenshuo/gopen/game"

type GameObjectNode struct {
	Name     string
	GB       game.GameObjectI
	Children []*GameObjectNode
}

type GameObjectNodeSpec struct {
	Name     string
	Kind     string // Nil Basic
	Pivot    []float32
	Size     []float32
	Children []*GameObjectNodeSpec
}

// 读取一个文件，根据文件内容生成复合模型
type CubeCustomTool struct {
}

func (cct *CubeCustomTool) LoadFromData(data []byte) {

}
