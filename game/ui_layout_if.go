package game

type UICanBeLayout interface {
	GetTransform() *Transform
	GetUISpec() *UISpec
}
