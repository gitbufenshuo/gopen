package modelcustom

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/gitbufenshuo/gopen/game/gameobjects"
	"github.com/gitbufenshuo/gopen/game/supports/logicinner"
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
	gi     *game.GlobalInfo
	acgbid int
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
	if cct.acgbid != 0 {
		ac := cct.gi.AnimationSystem.GetAC(cct.acgbid)
		ac.RecordInitFrame()
	}
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
		// nil 对象也是有 position 和 rotation 的
		if posdata, found := attrmap["pos"]; found {
			pos := matmath.CreateVec4FromStr(posdata)
			gbn.GB.GetTransform().Postion.Clone(&pos)
		}
		if rotationdata, found := attrmap["rotation"]; found {
			rotation := matmath.CreateVec4FromStr(rotationdata)
			gbn.GB.GetTransform().Rotation.Clone(&rotation)
		}
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
	if dongid, found := attrmap["dong"]; found {
		if node.Data == "blockroot" { // 根节点可能指定动画id
			cct.gi.AnimationSystem.CreateAC(dongid, gbn.GB.ID_sg()) // 创建 AnimationController
			cct.acgbid = gbn.GB.ID_sg()
		} else {
			cct.gi.AnimationSystem.BindBoneNode(cct.acgbid, dongid, gbn.GB.GetTransform())
		}
	}

	if node.Data == "blockroot" { // 根节点可能指定 内置logic
		if rvdata, found := attrmap["rotate"]; found { // 绑定一个 LogicRotate
			fmt.Println("rvdata", rvdata)
			gbn.GB.AddLogicSupport(logicinner.NewLogicRotate(rvdata))
		}
		if logicdata, found := attrmap["logic"]; found { // 自定义 logic
			fmt.Println("logic", logicdata)
			segs := strings.Split(logicdata, ",")
			for _, onelogic := range segs {
				gbn.GB.AddLogicSupport(cct.gi.LogicSystem.GetLogicByName(cct.gi, fmt.Sprintf("logic_%s", onelogic)))
			}
		}
	}

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
