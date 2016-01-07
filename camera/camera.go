package camera

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"math"

	"github.com/brandonnelson3/GoPlay/input"
)

var C FPS

const Pi2 = math.Pi / 2.0

type FPS struct {
	Position      mgl32.Vec3
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
		Position:      mgl32.Vec3{-6, 0, 0},
		Speed:         10.0,
		FOVDegrees:    45.0,
		NearPlaneDist: 0.1,
		FarPlaneDist:  10000.0,

		direction:       mgl32.Vec3{0, 0, 0},
		horizontalAngle: 0.0,
		verticalAngle:   0.0,
		sensitivity:     0.5,
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

// TODO: Move this to the base struct.
func (c *FPS) Update(d float64) {
	if c.direction.X() != 0 || c.direction.Y() != 0 || c.direction.Z() != 0 {
		c.Position = c.Position.Add(c.direction.Normalize().Mul(float32(d) * c.Speed))
		c.direction = mgl32.Vec3{0, 0, 0}
	}
}

// TODO: Move this to the base struct. Also not a bad idea to memoize the result, and update only when dirty.
func (c *FPS) GetForward() mgl32.Vec3 {
	return mgl32.Rotate3DY(c.horizontalAngle).Mul3x1(mgl32.Rotate3DZ(c.verticalAngle).Mul3x1((mgl32.Vec3{1, 0, 0})))
}

// TODO: Move this to the base
// struct. Also not a bad idea to memoize the result, and update only when dirty.
func (c *FPS) GetRight() mgl32.Vec3 {
	return mgl32.Rotate3DY(c.horizontalAngle).Mul3x1(mgl32.Vec3{0, 0, 1})
}

// TODO: Move this to the base struct.
func (c *FPS) GetViewMatrix() mgl32.Mat4 {
	return mgl32.LookAtV(c.Position, c.Position.Add(c.GetForward()), mgl32.Vec3{0, 1, 0})
}
