package game

import (
	"github.com/go-gl/glfw/v3.1/glfw"
)

type InputMouse struct {
	mouseXPos, mouseYPos   float64
	MouseXDiff, MouseYDiff float64
}

func NewInputMouse() *InputMouse {
	return new(InputMouse)
}

func (im *InputMouse) CursorCallback(win *glfw.Window, xpos float64, ypos float64) {
	im.MouseXDiff = xpos - im.mouseXPos
	im.MouseYDiff = ypos - im.mouseYPos
	im.mouseXPos, im.mouseYPos = xpos, ypos

}

func (im *InputMouse) MouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {

}
