package sceneloader

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/gameex/modelcustom"
)

// 资源加载
// 模型加载
// 等等

type SceneLoader struct {
	gi          *game.GlobalInfo
	cct         *modelcustom.CubeCustomTool
	SpecPath    string // 路径
	TextureList []string
}

func NewSceneLoader(gi *game.GlobalInfo, specpath string) *SceneLoader {
	res := new(SceneLoader)
	res.gi = gi
	res.SpecPath = specpath
	res.cct = modelcustom.NewCubeCustomTool(gi)
	return res
}

func (sl *SceneLoader) LoadTextureList() {
	filename := path.Join(sl.SpecPath, "pick", "texture.csv")
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("LoadTextureList", err)
		return
	}
	defer file.Close()
	//
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		segs := strings.Split(text, " ")
		textureid, texturepath := segs[0], segs[1]
		path := path.Join(sl.SpecPath, "asset", texturepath)
		sl.gi.AssetManager.LoadTextureFromFile(path, textureid)
		sl.TextureList = append(sl.TextureList, textureid)
	}
}

func (sl *SceneLoader) LoadCubeModelList() {
	filename := path.Join(sl.SpecPath, "pick", "cubemodel.csv")
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("LoadCubeModelList", err)
		return
	}
	defer file.Close()
	//
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		segs := strings.Split(text, " ")
		_, cubepath := segs[0], segs[1]
		path := path.Join(sl.SpecPath, "asset", cubepath)
		sl.cct.LoadFromFile(path)
	}
}
