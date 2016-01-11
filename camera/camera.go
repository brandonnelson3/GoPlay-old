package camera

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"math"
	"sync"

	"fmt"

	"github.com/brandonnelson3/GoPlay/input"
	"github.com/brandonnelson3/GoPlay/window"
)

var C FPS

const Pi2 = math.Pi / 2.0

type FPS struct {
	positionMu sync.RWMutex
	position   mgl32.Vec3

	Speed         float32
	FOVDegrees    float32
	NearPlaneDist float32
	FarPlaneDist  float32

	// TODO: This stuff could probably be abstracted out into a base struct under all camera types.
	direction       mgl32.Vec3
	horizontalAngle float32
	verticalAngle   float32
	sensitivity     float32
}

func init() {
	C = FPS{
		position:      mgl32.Vec3{9, 9, 9},
		Speed:         20.0,
		FOVDegrees:    45.0,
		NearPlaneDist: 0.1,
		FarPlaneDist:  10000.0,

		direction:       mgl32.Vec3{0, 0, 0},
		horizontalAngle: 4.33,
		verticalAngle:   .3,
		sensitivity:     0.1,
	}

	// Normal wasd movement
	input.M.Register(glfw.KeyW, C.Forward)
	input.M.Register(glfw.KeyS, C.Backward)
	input.M.Register(glfw.KeyA, C.StrafeLeft)
	input.M.Register(glfw.KeyD, C.StrafeRight)

	// Arrow keys adjust view angle
	input.M.Register(glfw.KeyUp, C.Up)
	input.M.Register(glfw.KeyDown, C.Down)
	input.M.Register(glfw.KeyLeft, C.Left)
	input.M.Register(glfw.KeyRight, C.Right)

	window.M.W.SetCursorPos(float64(window.M.Height)/2, float64(window.M.Width)/2)
	window.M.W.SetCursorPosCallback(C.cursorPosCallback)
	window.M.W.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
}

func (c *FPS) Forward(bool, float32) {
	c.direction = c.direction.Add(c.GetForward())
}
func (c *FPS) Backward(bool, float32) {
	c.direction = c.direction.Sub(c.GetForward())
}
func (c *FPS) StrafeRight(bool, float32) {
	c.direction = c.direction.Add(c.GetRight())
}
func (c *FPS) StrafeLeft(bool, float32) {
	c.direction = c.direction.Sub(c.GetRight())
}

func (c *FPS) Up(_ bool, d float32) {
	c.verticalAngle += c.sensitivity * d
	if c.verticalAngle > Pi2-0.0001 {
		c.verticalAngle = float32(Pi2 - 0.0001)
	}

	fmt.Printf("Angle: %v", c.verticalAngle)
}
func (c *FPS) Down(_ bool, d float32) {
	c.verticalAngle -= c.sensitivity * d
	if c.verticalAngle < -Pi2 {
		c.verticalAngle = float32(-Pi2 + 0.0001)
	}
}
func (c *FPS) Left(_ bool, d float32) {
	c.horizontalAngle += c.sensitivity * d
	for c.horizontalAngle > 2*math.Pi {
		c.horizontalAngle -= float32(2 * math.Pi)
	}
}
func (c *FPS) Right(_ bool, d float32) {
	c.horizontalAngle -= c.sensitivity * d
	for c.horizontalAngle < 0 {
		c.horizontalAngle += float32(2 * math.Pi)
	}
}

func (c *FPS) cursorPosCallback(w *glfw.Window, x, y float64) {
	x -= float64(window.M.Width) / 2
	y -= float64(window.M.Height) / 2
	c.verticalAngle -= c.sensitivity * float32(y) * .01
	if c.verticalAngle < -Pi2 {
		c.verticalAngle = float32(-Pi2 + 0.0001)
	}
	c.horizontalAngle -= c.sensitivity * float32(x) * .01
	for c.horizontalAngle < 0 {
		c.horizontalAngle += float32(2 * math.Pi)
	}
}

// TODO: Move this to the base struct.
func (c *FPS) Update(d float64) {
	if c.direction.X() != 0 || c.direction.Y() != 0 || c.direction.Z() != 0 {
		delta := c.direction.Normalize().Mul(float32(d) * c.Speed)
		c.positionMu.Lock()
		c.position = c.position.Add(delta)
		c.positionMu.Unlock()
		c.direction = mgl32.Vec3{0, 0, 0}
	}
}

func (c *FPS) GetPosition() mgl32.Vec3 {
	c.positionMu.RLock()
	defer c.positionMu.RUnlock()
	return c.position
}

// TODO: Move this to the base struct. Also not a bad idea to memoize the result, and update only when dirty.
func (c *FPS) GetForward() mgl32.Vec3 {
	return mgl32.Rotate3DY(c.horizontalAngle).Mul3x1(mgl32.Rotate3DZ(c.verticalAngle).Mul3x1((mgl32.Vec3{1, 0, 0})))
}

// TODO: Move this to the base struct. Also not a bad idea to memoize the result, and update only when dirty.
func (c *FPS) GetRight() mgl32.Vec3 {
	return mgl32.Rotate3DY(c.horizontalAngle).Mul3x1(mgl32.Vec3{0, 0, 1})
}

// TODO: Move this to the base struct.
func (c *FPS) GetViewMatrix() mgl32.Mat4 {
	c.positionMu.RLock()
	defer c.positionMu.RUnlock()
	return mgl32.LookAtV(c.position, c.position.Add(c.GetForward()), mgl32.Vec3{0, 1, 0})
}
