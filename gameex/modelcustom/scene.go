package modelcustom

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"path"
	"strings"

	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/gitbufenshuo/gopen/game/gameobjects"
	"github.com/gitbufenshuo/gopen/help"
	"github.com/gitbufenshuo/gopen/matmath"
	"golang.org/x/net/html"
)

type Scene struct {
	From       PrefabFrom
	RootNode   *SceneNode
	CameraNode *SceneNode
	LightNode  *SceneNode
	runtimeGB  map[string]game.GameObjectI
}

// 通过这个 scene 生成一整个 场景
func (pf *Scene) Instantiate(gi *game.GlobalInfo) game.GameObjectI {
	//
	pf.CreateCamera(gi)
	pf.CreateLight(gi)
	return pf.RootNode.instantiate(gi, pf)
}
func (pf *Scene) CreateLight(gi *game.GlobalInfo) {
	gi.MainLight = game.NewLight()
	gi.MainLight.SetLightColor(pf.LightNode.Color.GetValue3())
	gi.MainLight.SetLightDirection(pf.LightNode.Direction.GetValue3())
}
func (pf *Scene) CreateCamera(gi *game.GlobalInfo) {
	//
	gi.MainCamera = game.NewDefaultCamera()
	gi.MainCamera.Transform.Postion.SetValue3(pf.CameraNode.Pos.GetValue3())
	gi.MainCamera.SetForward(pf.CameraNode.Forward.GetValue3())
	gi.MainCamera.NearDistance = pf.CameraNode.Near
	{
		pngs := []string{}
		for idx := range pf.CameraNode.SkyBox {
			pngs = append(pngs,
				path.Join("scenespec/asset/skybox/", pf.CameraNode.SkyBox[idx]),
			)
		}
		// skybox
		// 1. load the cubemap
		cubemap := resource.NewCubeMap()
		// C17F4DB8CC0274DA12D60A6944762679.png
		// []string{
		// 	// "scenespec/asset/skybox/128.png",
		// 	// "scenespec/asset/skybox/128.png",
		// 	// "scenespec/asset/skybox/128.png",
		// 	// "scenespec/asset/skybox/128.png",
		// 	// "scenespec/asset/skybox/128.png",
		// 	// "scenespec/asset/skybox/right.png",
		// 	// "scenespec/asset/skybox/right.png",
		// 	// "scenespec/asset/skybox/left.png",
		// 	// "scenespec/asset/skybox/top.png",
		// 	// "scenespec/asset/skybox/bottom.png",
		// 	// "scenespec/asset/skybox/back.png",
		// 	// "scenespec/asset/skybox/front.png",
		// }
		cubemap.ReadFromFile(pngs)
		cubemap.Upload()
		gi.AssetManager.CreateCubemapSilent("skybox.cubemap", cubemap)
		model := resource.NewCubemapModel_BySpec(matmath.CreateVec4FromStr("0,0,0,1"), matmath.CreateVec4FromStr("2,2,2,1"))
		modelresourcename := fmt.Sprintf("cameraskybox.%d.%d", rand.Int()%1000, rand.Int()%1000000)
		// fmt.Println("modelresourcename", modelresourcename)
		gi.AssetManager.CreateModelSilent(modelresourcename, model)
		newGB := gameobjects.NewCubemapObject(gi, modelresourcename, "skybox.cubemap", "skybox_shader")
		newGB.AddLogicSupport(gi.LogicSystem.GetLogicByName(gi, "logic_zhuan"))
		gi.MainCamera.AddSkyBox(newGB)
	}

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
	newScene.RootNode = newScene.loadSceneFromContent(newScene.From.Content)
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
	newScene.RootNode = newScene.loadSceneFromContent(newScene.From.Content)
	SceneSystemIns.AddScene(pname, newScene)
	return newScene
}

func (sc *Scene) loadSceneFromContent(content []byte) *SceneNode {
	doc, err := html.Parse(bytes.NewReader(content))
	if err != nil {
		return nil
	}
	blockrootnode := FindHTMLRoot(doc, "root")
	sceneNode := readOneHTMLNode2SceneNode(blockrootnode)
	{
		// camera node
		cameraHTMLNode := FindHTMLRoot(doc, "camera")
		cameraNode := new(SceneNode)
		cameraNode.ReadCameraFromHTMLNode(cameraHTMLNode)
		sc.CameraNode = cameraNode
	}
	{
		// light node
		lightHTMLNode := FindHTMLRoot(doc, "light")
		lightNode := new(SceneNode)
		lightNode.ReadLightFromHTMLNode(lightHTMLNode)
		sc.LightNode = lightNode
	}
	return sceneNode
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
	Name      string // node name
	Kind      string // nil basic
	Dong      string //
	Logic     []string
	SkyBox    []string
	Model     string
	Image     string
	Prefab    string
	Color     matmath.Vec4
	Direction matmath.Vec4
	//
	Pos      matmath.Vec4 //
	Forward  matmath.Vec4
	Rotation matmath.Vec4 //
	Near     float32
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
		if prefab == nil {
			fmt.Println("SceneNode instantiate 错误", scene.From.From, scene.From.Meta)
			fmt.Printf("    --->找不到 prefab:%s\n", pn.Prefab)
			panic("")
		}
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

func (pn *SceneNode) ReadCameraFromHTMLNode(htmlnode *html.Node) {
	attrmap := make(map[string]string)
	for _, oneattr := range htmlnode.Attr {
		attrmap[oneattr.Key] = oneattr.Val
	}
	if v, found := attrmap["name"]; found {
		pn.Name = v
	}
	if v, found := attrmap["pos"]; found {
		pn.Pos = matmath.CreateVec4FromStr(v)
	}
	if v, found := attrmap["forward"]; found {
		pn.Forward = matmath.CreateVec3FromStr(v)
	}
	pn.Near = 0.5
	if v, found := attrmap["near"]; found {
		pn.Near = help.Str2Float32(v)
	}
	if v, found := attrmap["skybox"]; found {
		segs := strings.Split(v, ",")
		pn.SkyBox = segs
	}
}

func (pn *SceneNode) ReadLightFromHTMLNode(htmlnode *html.Node) {
	attrmap := make(map[string]string)
	for _, oneattr := range htmlnode.Attr {
		attrmap[oneattr.Key] = oneattr.Val
	}
	if v, found := attrmap["name"]; found {
		pn.Name = v
	}
	if v, found := attrmap["color"]; found {
		pn.Color = matmath.CreateVec3FromStr(v)
	}
	if v, found := attrmap["direction"]; found {
		pn.Direction = matmath.CreateVec3FromStr(v)
	}
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
