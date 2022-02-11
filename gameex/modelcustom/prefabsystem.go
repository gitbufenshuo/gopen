package modelcustom

type PrefabSystem struct {
	prefabs map[string]*Prefab
}

// 增加 prefab
func (pbs *PrefabSystem) AddPrefab(name string, prefab *Prefab) {
	pbs.prefabs[name] = prefab
}

// 通过名称 获取 prefab
func (pbs *PrefabSystem) GetPrefab(name string) *Prefab {
	return pbs.prefabs[name]
}

var PrefabSystemIns *PrefabSystem

func init() {
	PrefabSystemIns = new(PrefabSystem)
	PrefabSystemIns.prefabs = make(map[string]*Prefab)
}
