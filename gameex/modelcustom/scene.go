package modelcustom

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/gameobjects"
	"github.com/gitbufenshuo/gopen/matmath"
	"golang.org/x/net/html"
)

type Scene struct {
	From      PrefabFrom
	RootNode  *SceneNode
	runtimeGB map[string]game.GameObjectI
}

// 通过这个 scene 生成一整个 场景
func (pf *Scene) Instantiate(gi *game.GlobalInfo) game.GameObjectI {
	//
	return pf.RootNode.instantiate(gi, pf)
}

func LoadSceneFromFile(pname, filepath string) *Scene {
	newScene := new(Scene)
	newScene.runtimeGB = make(map[string]game.GameObjectI)
	//
	newScene.From.From = FromFile
	newScene.From.Meta = filepath
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil
	}
	newScene.From.Content = data
	newScene.RootNode = loadSceneFromContent(newScene.From.Content)
	SceneSystemIns.AddScene(pname, newScene)
	return newScene
}

func LoadSceneFromContent(pname string, content []byte) *Scene {
	newScene := new(Scene)
	newScene.runtimeGB = make(map[string]game.GameObjectI)
	//
	newScene.From.From = FromContent
	newScene.From.Content = content
	//
	newScene.RootNode = loadSceneFromContent(newScene.From.Content)
	SceneSystemIns.AddScene(pname, newScene)
	return newScene
}

func loadSceneFromContent(content []byte) *SceneNode {
	doc, err := html.Parse(bytes.NewReader(content))
	if err != nil {
		return nil
	}
	blockrootnode := FindHTMLRoot(doc, "root")
	return readOneHTMLNode2SceneNode(blockrootnode)
}

func readOneHTMLNode2SceneNode(htmlnode *html.Node) *SceneNode {
	thenode := new(SceneNode)
	thenode.ReadDataFromHTMLNode(htmlnode) // 将本node数据搞定
	// 然后搞多个子节点
	child := htmlnode.FirstChild
	for ; child != nil; child = child.NextSibling {
		if child.Data == "gameobject" {
			childNode := readOneHTMLNode2SceneNode(child)
			thenode.Children = append(thenode.Children, childNode)
		}
	}
	return thenode
}

func FindHTMLRoot(node *html.Node, rootname string) *html.Node {
	if node.Data == rootname {
		return node
	}
	for _node := node.FirstChild; _node != nil; _node = _node.NextSibling {
		if _n := FindHTMLRoot(_node, rootname); _n != nil {
			return _n
		}
	}
	return nil
}

type SceneNode struct {
	Name   string // node name
	Kind   string // nil basic
	Dong   string //
	Logic  []string
	Model  string
	Image  string
	Prefab string
	//
	Pos      matmath.Vec4 //
	Rotation matmath.Vec4 //
	Pivot    matmath.Vec4
	Size     matmath.Vec4
	///////////////
	Children []*SceneNode
}

func (pn *SceneNode) instantiate(gi *game.GlobalInfo, scene *Scene) game.GameObjectI {

	var res game.GameObjectI
	if pn.Prefab == "" {
		newGB := gameobjects.NewNilObject(gi)
		// nil 对象也是有 position 和 rotation 的
		newGB.GetTransform().Postion.Clone(&pn.Pos)
		newGB.GetTransform().Rotation.Clone(&pn.Rotation)
		res = newGB
		for _, onelogic := range pn.Logic {
			res.AddLogicSupport(gi.LogicSystem.GetLogicByName(gi, fmt.Sprintf("logic_%s", onelogic)))
		}
	} else { // 这是一个prefab
		prefab := PrefabSystemIns.GetPrefab(pn.Prefab)
		newgb := prefab.Instantiate(gi)
		newgb.GetTransform().Postion.Clone(&pn.Pos)
		res = newgb
	}
	gi.AddGameObject(res)
	//
	for _, onechild := range pn.Children {
		newChildGB := onechild.instantiate(gi, scene)
		newChildGB.GetTransform().SetParent(res.GetTransform())
	}
	scene.runtimeGB[pn.Name] = res
	return res
}
func (pn *SceneNode) ReadDataFromHTMLNode(htmlnode *html.Node) {
	attrmap := make(map[string]string)
	for _, oneattr := range htmlnode.Attr {
		attrmap[oneattr.Key] = oneattr.Val
	}
	// name
	if v, found := attrmap["name"]; found {
		pn.Name = v
	}
	// prefab
	if v, found := attrmap["prefab"]; found {
		pn.Prefab = v
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
