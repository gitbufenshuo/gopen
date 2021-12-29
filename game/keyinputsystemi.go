package game

type InputSystemI interface {
	Start()
	Update()
	BeginWatchKey(key int)
	StopWatchKey(key int)
	KeyListInWatching() []int
	KeyDown(key int) bool
	KeyUp(key int) bool
	KeyDoubleClick(key int) bool
	KeyHoldRelease(key int) float64
}
