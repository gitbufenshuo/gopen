package sceneloader

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/gameex/animationsystem"
	"github.com/gitbufenshuo/gopen/gameex/modelcustom"
)

// 资源加载
// 模型加载
// 动画加载
// 等等

type SceneLoader struct {
	gi          *game.GlobalInfo
	SpecPath    string // 路径
	TextureList []string
	//
	GameMap map[string]game.GameObjectI
}

var SceneLoaderMap map[string]*SceneLoader

func FindGameobjectByName(scpath, gbname string) game.GameObjectI {
	scene := SceneLoaderMap[scpath]
	return scene.GameMap[gbname]
}

func NewSceneLoader(gi *game.GlobalInfo, specpath string) *SceneLoader {
	res := new(SceneLoader)
	res.gi = gi
	res.SpecPath = specpath
	res.GameMap = make(map[string]game.GameObjectI)
	if SceneLoaderMap == nil {
		SceneLoaderMap = make(map[string]*SceneLoader)
	}
	SceneLoaderMap[specpath] = res
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
		path := path.Join(sl.SpecPath, "asset", "png", texturepath)
		sl.gi.AssetManager.LoadTextureFromFile(path, textureid)
		sl.TextureList = append(sl.TextureList, textureid)
	}
}

func (sl *SceneLoader) LoadDongList() {
	filename := path.Join(sl.SpecPath, "pick", "dong.csv")
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("LoadDongList", err)
		return
	}
	defer file.Close()
	//
	as := animationsystem.NewAnimationSystem(sl.gi) // 动画统一管理器
	sl.gi.AnimationSystem = as
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "//") {
			continue
		}
		segs := strings.Split(text, " ")
		dongname, dongpath := segs[0], segs[1]
		path := path.Join(sl.SpecPath, "asset", "dong", dongpath)
		as.AddAnimationMeta(dongname, animationsystem.LoadAnimationMetaFromFile(path))
	}
}

func (sl *SceneLoader) LoadPrefabList() {
	filename := path.Join(sl.SpecPath, "pick", "prefab.csv")
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("LoadPrefabList", err)
		return
	}
	defer file.Close()
	//
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "//") {
			continue
		}
		segs := strings.Split(text, " ")
		_name, _path := segs[0], segs[1]
		fullpath := path.Join(sl.SpecPath, "asset", "prefab", _path)
		modelcustom.LoadPrefabFromFile(_name, fullpath)
	}
}

func (sl *SceneLoader) LoadSceneList() {
	filename := path.Join(sl.SpecPath, "pick", "scene.csv")
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("LoadSceneList", err)
		return
	}
	defer file.Close()
	//
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "//") {
			continue
		}
		segs := strings.Split(text, " ")
		_name, _path := segs[0], segs[1]
		fullpath := path.Join(sl.SpecPath, "asset", "scene", _path)
		modelcustom.LoadSceneFromFile(_name, fullpath)
	}
}
