package input

import (
	"log"

	"github.com/go-gl/glfw/v3.1/glfw"
)

const keyRange = 349

type keyFunction func(bool)
type manager struct {
	down          []bool
	downThisFrame []bool
	functions     [][]keyFunction
}

func NewManager(w *glfw.Window) manager {
	result := manager{
		down:          make([]bool, keyRange),
		downThisFrame: make([]bool, keyRange),
		functions:     make([][]keyFunction, keyRange),
	}
	w.SetKeyCallback(result.keyCallBack)
	return result
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
func (inputManager *manager) Register(keycode uint, f keyFunction) {
	if keycode < uint(glfw.KeySpace) || keycode > uint(glfw.KeyMenu) {
		log.Fatalf("Got invalid keycode in Register: %v", keycode)
		return
	}
	inputManager.functions[keycode] = append(inputManager.functions[keycode], f)
}

func (inputManager *manager) RunKeys() {
	// glfw.KeySpace is the lowest key.
	for i := glfw.KeySpace; i < keyRange; i++ {
		if inputManager.down[i] {
			for j := 0; j < len(inputManager.functions[i]); j++ {
				// Send true on all function calls except for the very first one after a button press.
				inputManager.functions[i][j](!inputManager.downThisFrame[i])
				inputManager.downThisFrame[i] = false
			}
		}
	}
}
