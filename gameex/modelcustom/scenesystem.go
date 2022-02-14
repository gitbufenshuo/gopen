package modelcustom

type SceneSystem struct {
	scenes map[string]*Scene
}

// 增加 prefab
func (scs *SceneSystem) AddScene(name string, scene *Scene) {
	scs.scenes[name] = scene
}

// 通过名称 获取 prefab
func (scs *SceneSystem) GetScene(name string) *Scene {
	return scs.scenes[name]
}

var SceneSystemIns *SceneSystem

func init() {
	SceneSystemIns = new(SceneSystem)
	SceneSystemIns.scenes = make(map[string]*Scene)
}
