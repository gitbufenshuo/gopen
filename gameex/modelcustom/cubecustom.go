package modelcustom

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/gitbufenshuo/gopen/game/gameobjects"
	"github.com/gitbufenshuo/gopen/matmath"
	"golang.org/x/net/html"
)

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
	gi *game.GlobalInfo
}

func NewCubeCustomTool(gi *game.GlobalInfo) *CubeCustomTool {
	res := new(CubeCustomTool)
	//
	res.gi = gi
	return res
}
func (cct *CubeCustomTool) LoadFromFile(path string) *GameObjectNode {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil
	}
	return cct.LoadFromData(data)
}

func (cct *CubeCustomTool) LoadFromData(data []byte) *GameObjectNode {
	doc, err := html.Parse(bytes.NewReader(data))
	if err != nil {
		return nil
	}
	blockrootnode := cct.FindBlockRoot(doc)
	gbn := new(GameObjectNode)
	cct.ScanNode(blockrootnode, gbn)
	return nil
}
func (cct *CubeCustomTool) FindBlockRoot(node *html.Node) *html.Node {
	if node.Data == "blockroot" {
		return node
	}
	for _node := node.FirstChild; _node != nil; _node = _node.NextSibling {
		if _n := cct.FindBlockRoot(_node); _n != nil {
			return _n
		}
	}
	return nil
}
func (cct *CubeCustomTool) ScanNode(node *html.Node, gbn *GameObjectNode) {
	attrmap := make(map[string]string)
	for _, oneattr := range node.Attr {
		attrmap[oneattr.Key] = oneattr.Val
	}
	//
	gbn.Name = attrmap["name"]
	//
	if attrmap["kind"] == "nil" {
		gbn.GB = gameobjects.NewNilObject(cct.gi)
	} else {
		// 根据 pivot 和 size 生成模型
		fmt.Println(node.Data, "pivot", attrmap["pivot"])
		fmt.Println(node.Data, "size", attrmap["size"])
		fmt.Println(node.Data, "image", attrmap["image"])

		model := resource.NewBlockModel_BySpec(
			matmath.CreateVec4FromStr(attrmap["pivot"]),
			matmath.CreateVec4FromStr(attrmap["size"]),
		)
		modelresourcename := fmt.Sprintf("html.%s.%d%p", node.Data, node.DataAtom, node)
		fmt.Println("modelresourcename", modelresourcename)
		cct.gi.AssetManager.CreateModelSilent(modelresourcename, model)
		gbn.GB = gameobjects.NewBasicObject(cct.gi, modelresourcename, attrmap["image"], "mvp_shader")
		pos := matmath.CreateVec4FromStr(attrmap["pos"])
		gbn.GB.GetTransform().Postion.Clone(&pos)
		rotation := matmath.CreateVec4FromStr(attrmap["rotation"])
		gbn.GB.GetTransform().Rotation.Clone(&rotation)
	}
	cct.gi.AddGameObject(gbn.GB)
	fmt.Println("ScanNode", gbn.GB)
	// 考虑下级 跟链表一样
	child := node.FirstChild
	for ; child != nil; child = child.NextSibling {
		if child.Data == "block" {
			childgbn := new(GameObjectNode)
			gbn.Children = append(gbn.Children, childgbn)
			cct.ScanNode(child, childgbn)
		}
	}
	for idx := range gbn.Children {
		gbn.Children[idx].GB.GetTransform().SetParent(gbn.GB.GetTransform())
	}
}
