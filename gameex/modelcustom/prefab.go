package modelcustom

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"

	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/gitbufenshuo/gopen/game/gameobjects"
	"github.com/gitbufenshuo/gopen/matmath"
	"golang.org/x/net/html"
)

type FromKind string

const (
	FromFile    FromKind = "FromFile"
	FromContent FromKind = "FromContent"
)

type PrefabFrom struct {
	From    FromKind //
	Meta    string   // filepath
	Content []byte   // content
}

type Prefab struct {
	From     PrefabFrom
	RootNode *PrefabNode
}

// 通过这个prefab 生成一个 gameobject
func (pf *Prefab) Instantiate(gi *game.GlobalInfo) game.GameObjectI {
	//
	gb := pf.RootNode.instantiate(gi)
	if ac := gb.GetACSupport(); ac != nil {
		ac.RecordInitFrame()
	}
	return gb
}

func LoadPrefabFromFile(pname, filepath string) *Prefab {
	newPrefab := new(Prefab)
	//
	newPrefab.From.From = FromFile
	newPrefab.From.Meta = filepath
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil
	}
	newPrefab.From.Content = data
	newPrefab.RootNode = loadPrefabFromContent(newPrefab.From.Content)
	PrefabSystemIns.AddPrefab(pname, newPrefab)
	return newPrefab
}

func LoadPrefabFromContent(pname string, content []byte) *Prefab {
	newPrefab := new(Prefab)
	//
	newPrefab.From.From = FromContent
	newPrefab.From.Content = content
	//
	newPrefab.RootNode = loadPrefabFromContent(newPrefab.From.Content)
	PrefabSystemIns.AddPrefab(pname, newPrefab)
	return newPrefab
}

func loadPrefabFromContent(content []byte) *PrefabNode {
	doc, err := html.Parse(bytes.NewReader(content))
	if err != nil {
		return nil
	}
	blockrootnode := FindBlockRoot(doc)
	return readOneHTMLNode(blockrootnode)
	// gbn := new(GameObjectNode)
	// cct.ScanNode(blockrootnode, gbn)
	// if cct.acgbid != 0 {
	// 	ac := cct.gi.AnimationSystem.GetAC(cct.acgbid)
	// 	ac.RecordInitFrame()
	// }
	// return gbn

}

func readOneHTMLNode(htmlnode *html.Node) *PrefabNode {
	thenode := new(PrefabNode)
	thenode.ReadDataFromHTMLNode(htmlnode) // 将本node数据搞定
	// 然后搞多个子节点
	child := htmlnode.FirstChild
	for ; child != nil; child = child.NextSibling {
		if child.Data == "block" {
			childNode := readOneHTMLNode(child)
			thenode.Children = append(thenode.Children, childNode)
		}
	}
	return thenode
}

func FindBlockRoot(node *html.Node) *html.Node {
	if node.Data == "blockroot" {
		return node
	}
	for _node := node.FirstChild; _node != nil; _node = _node.NextSibling {
		if _n := FindBlockRoot(_node); _n != nil {
			return _n
		}
	}
	return nil
}

type PrefabNode struct {
	Data  string
	Name  string // node name
	Kind  string // nil basic
	Dong  string //
	Logic []string
	Model string
	Image string
	//
	Pos      matmath.Vec4 //
	Rotation matmath.Vec4 //
	Pivot    matmath.Vec4
	Size     matmath.Vec4
	///////////////
	Children []*PrefabNode
	ac       game.AnimationControllerI
}

func (pn *PrefabNode) instantiate(gi *game.GlobalInfo) game.GameObjectI {
	var res game.GameObjectI
	if pn.Kind == "nil" {
		newGB := gameobjects.NewNilObject(gi)
		// nil 对象也是有 position 和 rotation 的
		newGB.GetTransform().Postion.Clone(&pn.Pos)
		newGB.GetTransform().Rotation.Clone(&pn.Rotation)
		res = newGB
	} else {
		var modelkind = "block"
		if pn.Model != "" {
			modelkind = pn.Model
		}
		var model *resource.Model
		if modelkind == "block" {
			model = resource.NewBlockModel_BySpec(pn.Pivot, pn.Size)
		} else {
			model = resource.NewSphereModel_BySpec(pn.Pivot, pn.Size)
		}
		modelresourcename := fmt.Sprintf("prefabhtml.%s.%d.%d", pn.Name, rand.Int()%1000, rand.Int()%1000000)
		// fmt.Println("modelresourcename", modelresourcename)
		gi.AssetManager.CreateModelSilent(modelresourcename, model)
		newGB := gameobjects.NewBasicObject(gi, modelresourcename, pn.Image, "mvp_shader")
		newGB.GetTransform().Postion.Clone(&pn.Pos)
		newGB.GetTransform().Rotation.Clone(&pn.Rotation)
		res = newGB
	}
	for _, onelogic := range pn.Logic {
		res.AddLogicSupport(gi.LogicSystem.GetLogicByName(gi, fmt.Sprintf("logic_%s", onelogic)))
	}
	gi.AddGameObject(res)
	//
	if pn.Dong != "" {
		// 可能有动画
		fmt.Println(pn.Data, pn.Dong)
		if pn.Data == "blockroot" { // 根节点
			pn.ac = gi.AnimationSystem.CreateAC(pn.Dong)
			res.SetACSupport(pn.ac)
		} else {
			pn.ac.BindBoneNode(pn.Dong, res.GetTransform())
		}
	}
	for _, onechild := range pn.Children {
		onechild.ac = pn.ac
		newChildGB := onechild.instantiate(gi)
		newChildGB.GetTransform().SetParent(res.GetTransform())
	}
	return res
}
func (pn *PrefabNode) ReadDataFromHTMLNode(htmlnode *html.Node) {
	attrmap := make(map[string]string)
	for _, oneattr := range htmlnode.Attr {
		attrmap[oneattr.Key] = oneattr.Val
	}
	// data

	pn.Data = htmlnode.Data

	// name
	if v, found := attrmap["name"]; found {
		pn.Name = v
	}
	// kind
	if v, found := attrmap["kind"]; found {
		pn.Kind = v
	}
	// Dong
	if v, found := attrmap["dong"]; found {
		pn.Dong = v
	}
	// Model
	if v, found := attrmap["model"]; found {
		pn.Model = v
	}
	// Image
	if v, found := attrmap["image"]; found {
		pn.Image = v
	}
	// Logic list
	if v, found := attrmap["logic"]; found {
		segs := strings.Split(v, ",")
		pn.Logic = segs
	}
	// 数值型
	if v, found := attrmap["pos"]; found {
		pn.Pos = matmath.CreateVec4FromStr(v)
	}
	if v, found := attrmap["rotation"]; found {
		pn.Rotation = matmath.CreateVec4FromStr(v)
	}
	if v, found := attrmap["pivot"]; found {
		pn.Pivot = matmath.CreateVec4FromStr(v)
	}
	if v, found := attrmap["size"]; found {
		pn.Size = matmath.CreateVec4FromStr(v)
	}
}
