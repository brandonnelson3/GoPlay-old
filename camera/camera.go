package camera

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brandonnelson3/GoPlay/input"
)

var C FPS

type FPS struct {
	Position  mgl32.Vec3
	Speed     float32
	direction mgl32.Vec3
}

func init() {
	C = FPS{Position: mgl32.Vec3{3, 3, 3}, Speed: 10, direction: mgl32.Vec3{0, 0, 0}}
	input.M.Register(glfw.KeyW, C.Forward)
	input.M.Register(glfw.KeyS, C.Backward)
}

func (c *FPS) Forward(w bool) {
	c.direction = c.direction.Add(mgl32.Vec3{1, 0, 0})
}

func (c *FPS) Backward(w bool) {
	c.direction = c.direction.Add(mgl32.Vec3{-1, 0, 0})
}

func (c *FPS) Update(d float64) {
	if c.direction.X() != 0 || c.direction.Y() != 0 || c.direction.Z() != 0 {
		c.Position = c.Position.Add(c.direction.Normalize().Mul(float32(d) * c.Speed))
		c.direction = mgl32.Vec3{0, 0, 0}
	}
}
