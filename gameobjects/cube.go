package gameobjects

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brandonnelson3/GoPlay/camera"
	"github.com/brandonnelson3/GoPlay/shaders"
	"github.com/brandonnelson3/GoPlay/texture"
)

var cubeVertices = []shaders.DefaultShader_Vertex{
	// Bottom
	{mgl32.Vec3{-1.0, -1.0, -1.0}, mgl32.Vec2{0.0, 0.0}},
	{mgl32.Vec3{1.0, -1.0, -1.0}, mgl32.Vec2{1.0, 0.0}},
	{mgl32.Vec3{-1.0, -1.0, 1.0}, mgl32.Vec2{0.0, 1.0}},
	{mgl32.Vec3{1.0, -1.0, -1.0}, mgl32.Vec2{1.0, 0.0}},
	{mgl32.Vec3{1.0, -1.0, 1.0}, mgl32.Vec2{1.0, 1.0}},
	{mgl32.Vec3{-1.0, -1.0, 1.0}, mgl32.Vec2{0.0, 1.0}},

	// Top
	{mgl32.Vec3{-1.0, 1.0, -1.0}, mgl32.Vec2{0.0, 0.0}},
	{mgl32.Vec3{-1.0, 1.0, 1.0}, mgl32.Vec2{0.0, 1.0}},
	{mgl32.Vec3{1.0, 1.0, -1.0}, mgl32.Vec2{1.0, 0.0}},
	{mgl32.Vec3{1.0, 1.0, -1.0}, mgl32.Vec2{1.0, 0.0}},
	{mgl32.Vec3{-1.0, 1.0, 1.0}, mgl32.Vec2{0.0, 1.0}},
	{mgl32.Vec3{1.0, 1.0, 1.0}, mgl32.Vec2{1.0, 1.0}},

	// Front
	{mgl32.Vec3{-1.0, -1.0, 1.0}, mgl32.Vec2{1.0, 0.0}},
	{mgl32.Vec3{1.0, -1.0, 1.0}, mgl32.Vec2{0.0, 0.0}},
	{mgl32.Vec3{-1.0, 1.0, 1.0}, mgl32.Vec2{1.0, 1.0}},
	{mgl32.Vec3{1.0, -1.0, 1.0}, mgl32.Vec2{0.0, 0.0}},
	{mgl32.Vec3{1.0, 1.0, 1.0}, mgl32.Vec2{0.0, 1.0}},
	{mgl32.Vec3{-1.0, 1.0, 1.0}, mgl32.Vec2{1.0, 1.0}},

	// Back
	{mgl32.Vec3{-1.0, -1.0, -1.0}, mgl32.Vec2{0.0, 0.0}},
	{mgl32.Vec3{-1.0, 1.0, -1.0}, mgl32.Vec2{0.0, 1.0}},
	{mgl32.Vec3{1.0, -1.0, -1.0}, mgl32.Vec2{1.0, 0.0}},
	{mgl32.Vec3{1.0, -1.0, -1.0}, mgl32.Vec2{1.0, 0.0}},
	{mgl32.Vec3{-1.0, 1.0, -1.0}, mgl32.Vec2{0.0, 1.0}},
	{mgl32.Vec3{1.0, 1.0, -1.0}, mgl32.Vec2{1.0, 1.0}},

	// Left
	{mgl32.Vec3{-1.0, -1.0, 1.0}, mgl32.Vec2{0.0, 1.0}},
	{mgl32.Vec3{-1.0, 1.0, -1.0}, mgl32.Vec2{1.0, 0.0}},
	{mgl32.Vec3{-1.0, -1.0, -1.0}, mgl32.Vec2{0.0, 0.0}},
	{mgl32.Vec3{-1.0, -1.0, 1.0}, mgl32.Vec2{0.0, 1.0}},
	{mgl32.Vec3{-1.0, 1.0, 1.0}, mgl32.Vec2{1.0, 1.0}},
	{mgl32.Vec3{-1.0, 1.0, -1.0}, mgl32.Vec2{1.0, 0.0}},

	// Right
	{mgl32.Vec3{1.0, -1.0, 1.0}, mgl32.Vec2{1.0, 1.0}},
	{mgl32.Vec3{1.0, -1.0, -1.0}, mgl32.Vec2{1.0, 0.0}},
	{mgl32.Vec3{1.0, 1.0, -1.0}, mgl32.Vec2{0.0, 0.0}},
	{mgl32.Vec3{1.0, -1.0, 1.0}, mgl32.Vec2{1.0, 1.0}},
	{mgl32.Vec3{1.0, 1.0, -1.0}, mgl32.Vec2{0.0, 0.0}},
	{mgl32.Vec3{1.0, 1.0, 1.0}, mgl32.Vec2{0.0, 1.0}},
}

type cube struct {
	vbo     *shaders.DefaultShader_VertexBuffer
	shader  *shaders.DefaultShader
	texture texture.Texture

	angle float64
}

func NewCube() (*cube, error) {
	shader, err := shaders.NewDefaultShader()
	if err != nil {
		return nil, err
	}
	shader.Activate()
	shader.SetProjection(mgl32.Perspective(mgl32.DegToRad(45.0), float32(800)/600, 0.1, 10.0))
	shader.SetModel(mgl32.Ident4())

	vbo := shaders.NewDefaultShader_VertexBuffer(shader, cubeVertices)

	// Load the texture
	texture, err := texture.New("assets/crate.jpg")
	if err != nil {
		return nil, err
	}

	return &cube{vbo: vbo, shader: shader, texture: texture}, nil
}

func (c *cube) Update(t float64) {
	c.angle += t
	c.shader.SetModel(mgl32.HomogRotate3D(float32(c.angle), mgl32.Vec3{0, 1, 0}))
}

func (c *cube) Render() {
	c.shader.Activate()
	c.shader.SetCamera(mgl32.LookAtV(camera.C.Position, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0}))
	c.texture.Bind(gl.TEXTURE0)
	gl.DrawArrays(gl.TRIANGLES, 0, c.vbo.Size)
}
