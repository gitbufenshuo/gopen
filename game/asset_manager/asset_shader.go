package asset_manager

import "github.com/gitbufenshuo/gopen/game/asset_manager/resource"

func (am *AsssetManager) LoadShaderFromFile(vetexPath, fragPath, assetname string) {
	var data ShaderDataType
	data.VPath = vetexPath
	data.FPath = fragPath
	as := NewAsset(assetname, AssetTypeShader, &data)
	err := am.Register(as.Name, as)
	if err != nil {
		panic(err)
	}
	am.Load(as)
	as.Resource.Upload()
}
func (am *AsssetManager) LoadShaderFromText(vetexText, fragText, assetname string) {
	shaderR := resource.NewShaderProgram()
	shaderR.ReadFromText(vetexText, fragText)
	shaderR.Upload()
	//
	as := NewAsset(assetname, AssetTypeShader, nil)
	err := am.Register(as.Name, as)
	if err != nil {
		panic(err)
	}
	as.Resource = shaderR
	am.Load(as)
}
