package game

import (
	"github.com/go-gl/glfw/v3.1/glfw"
)

type InputMouse struct {
	gi                     *GlobalInfo
	MouseXPos, MouseYPos   float64
	MouseXR, MouseYR       float32
	MouseXDiff, MouseYDiff float64
}

func NewInputMouse(gi *GlobalInfo) *InputMouse {
	res := new(InputMouse)
	res.gi = gi
	return res
}

func (im *InputMouse) Update() {
	im.MouseXPos, im.MouseYPos = im.gi.window.GetCursorPos()
	rwidthhalf, rheighthalf := im.gi.GetWindowWidth()/2, im.gi.GetWindowHeight()/2
	xr := float32(im.MouseXPos) - rwidthhalf
	xr /= rwidthhalf

	yr := float32(im.MouseYPos) - rheighthalf
	yr /= rheighthalf

	im.MouseXR, im.MouseYR = xr, -yr
}

func (im *InputMouse) CursorCallback(win *glfw.Window, xpos float64, ypos float64) {
	// im.MouseXDiff = xpos - im.MouseXPos
	// im.MouseYDiff = ypos - im.MouseYPos
	// im.MouseXPos, im.MouseYPos = xpos, ypos

}

func (im *InputMouse) MouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {

}
