package camera

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brandonnelson3/GoPlay/input"
)

var C FPS

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
}

func init() {
	C = FPS{
		Position:      mgl32.Vec3{3, 3, 3},
		Speed:         10.0,
		FOVDegrees:    45.0,
		NearPlaneDist: 0.1,
		FarPlaneDist:  10000.0,

		direction:       mgl32.Vec3{0, 0, 0},
		horizontalAngle: 0.0,
		verticalAngle:   0.0,
	}

	// Normal wasd movement
	input.M.Register(glfw.KeyW, C.Forward)
	input.M.Register(glfw.KeyS, C.Backward)
	input.M.Register(glfw.KeyA, C.Left)
	input.M.Register(glfw.KeyD, C.Right)

	// Arrow keys adjust view angle
}

func (c *FPS) Forward(bool) {
	c.direction = c.direction.Add(mgl32.Vec3{1, 0, 0})
}

func (c *FPS) Backward(bool) {
	c.direction = c.direction.Add(mgl32.Vec3{-1, 0, 0})
}

func (c *FPS) Left(bool) {
	c.direction = c.direction.Add(mgl32.Vec3{0, 0, 1})
}

func (c *FPS) Right(bool) {
	c.direction = c.direction.Add(mgl32.Vec3{0, 0, -1})
}

// TODO: Move this to the base struct.
func (c *FPS) Update(d float64) {
	if c.direction.X() != 0 || c.direction.Y() != 0 || c.direction.Z() != 0 {
		c.Position = c.Position.Add(c.direction.Normalize().Mul(float32(d) * c.Speed))
		c.direction = mgl32.Vec3{0, 0, 0}
	}
}

// TODO: Move this to the base struct.
func (c *FPS) GetViewMatrix() mgl32.Mat4 {
	// The second parameter here should be calculated based on the horizontal and vertical angle.
	return mgl32.LookAtV(c.Position, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
}
