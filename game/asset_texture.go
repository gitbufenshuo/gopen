package game

type SourceType int

const (
	SourceFromFile SourceType = 1 // from local file system
)

type AssetTexture struct {
	Source     SourceType // which place this asset is loaded from
	SourceInfo string     // the source detail such as "$assetpath/my_texture.png"

}
