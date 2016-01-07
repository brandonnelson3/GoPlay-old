package input

import (
	"log"

	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/brandonnelson3/GoPlay/window"
)

const keyRange = 349

// Global input manager
var M manager

type keyFunction func(bool, float32)
type manager struct {
	down          []bool
	downThisFrame []bool
	functions     [][]keyFunction
}

func init() {
	M = manager{
		down:          make([]bool, keyRange),
		downThisFrame: make([]bool, keyRange),
		functions:     make([][]keyFunction, keyRange),
	}
	window.M.W.SetKeyCallback(M.keyCallBack)
}

func (inputManager *manager) keyCallBack(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		log.Printf("Got key press event: %v", key)
		inputManager.down[key] = true
		inputManager.downThisFrame[key] = true
	}
	if action == glfw.Release {
		log.Printf("Got key release event: %v", key)
		inputManager.down[key] = false
		inputManager.downThisFrame[key] = false
	}
}

// Keycode values are from here: http://www.glfw.org/docs/latest/group__keys.html
func (inputManager *manager) Register(key glfw.Key, f keyFunction) {
	if key < glfw.KeySpace || key > glfw.KeyMenu {
		log.Fatalf("Got invalid keycode in Register: %v", key)
		return
	}
	inputManager.functions[key] = append(inputManager.functions[key], f)
}

func (inputManager *manager) RunKeys(d float32) {
	// glfw.KeySpace is the lowest key.
	for i := glfw.KeySpace; i < keyRange; i++ {
		if inputManager.down[i] {
			for j := 0; j < len(inputManager.functions[i]); j++ {
				// Send true on all function calls except for the very first one after a button press.
				inputManager.functions[i][j](!inputManager.downThisFrame[i], d)
				inputManager.downThisFrame[i] = false
			}
		}
	}
}
