package shaders

import (
	"fmt"
	"strings"

	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const vsSrc = `#version 330
uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

in vec3 vert;
in vec2 vertTexCoord;

out vec2 fragTexCoord;

void main() {
    fragTexCoord = vertTexCoord;
    gl_Position = projection * camera * model * vec4(vert, 1);
}` + "\x00"

const fsSrc = `#version 330
uniform sampler2D tex;

in vec2 fragTexCoord;

out vec4 outputColor;

void main() {
    outputColor = texture(tex, fragTexCoord);
}` + "\x00"

type DefaultShader struct {
	id uint32

	uniformLocation_projection int32
	uniformLocation_camera     int32
	uniformLocation_model      int32
	uniformLocation_tex        int32
}

func NewDefaultShader() (*DefaultShader, error) {
	pID := gl.CreateProgram()

	// Vertex Shader
	vsID := gl.CreateShader(gl.VERTEX_SHADER)
	cSrc := gl.Str(string(vsSrc))
	gl.ShaderSource(vsID, 1, &cSrc, nil)
	gl.CompileShader(vsID)
	var status int32
	gl.GetShaderiv(vsID, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(vsID, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(vsID, logLength, nil, gl.Str(log))
		return nil, fmt.Errorf("failed to compile vertex shader: %v", log)
	}
	gl.AttachShader(pID, vsID)

	// Fragment Shader
	fsID := gl.CreateShader(gl.FRAGMENT_SHADER)
	cSrc = gl.Str(string(fsSrc))
	gl.ShaderSource(fsID, 1, &cSrc, nil)
	gl.CompileShader(fsID)
	gl.GetShaderiv(fsID, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(fsID, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(fsID, logLength, nil, gl.Str(log))
		return nil, fmt.Errorf("failed to compile fragment shader: %v", log)
	}
	gl.AttachShader(pID, fsID)

	// Program
	gl.LinkProgram(pID)
	gl.GetProgramiv(pID, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(pID, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(pID, logLength, nil, gl.Str(log))

		return nil, fmt.Errorf("failed to link program: %v", log)
	}

	// Fragment Outputs
	gl.BindFragDataLocation(pID, 0, gl.Str("outputColor\x00"))

	result := &DefaultShader{
		id: pID,
		uniformLocation_projection: gl.GetUniformLocation(pID, gl.Str("projection\x00")),
		uniformLocation_camera:     gl.GetUniformLocation(pID, gl.Str("camera\x00")),
		uniformLocation_model:      gl.GetUniformLocation(pID, gl.Str("model\x00")),
		uniformLocation_tex:        gl.GetUniformLocation(pID, gl.Str("tex\x00")),
	}

	// Set Texture to slot 0
	gl.Uniform1i(result.uniformLocation_tex, 0)

	return result, nil
}

func (s *DefaultShader) Activate() {
	gl.UseProgram(s.id)
}

func (s *DefaultShader) SetProjection(d mgl32.Mat4) {
	gl.UniformMatrix4fv(s.uniformLocation_projection, 1, false, &d[0])
}

func (s *DefaultShader) SetCamera(d mgl32.Mat4) {
	gl.UniformMatrix4fv(s.uniformLocation_camera, 1, false, &d[0])
}

func (s *DefaultShader) SetModel(d mgl32.Mat4) {
	gl.UniformMatrix4fv(s.uniformLocation_model, 1, false, &d[0])
}

type DefaultShader_Vertex struct {
	Vert         mgl32.Vec3
	VertTexCoord mgl32.Vec2
}

type DefaultShader_VertexBuffer struct {
	id   uint32
	Size int32
}

func NewDefaultShader_VertexBuffer(s *DefaultShader, d []DefaultShader_Vertex) *DefaultShader_VertexBuffer {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(d)*20, gl.Ptr(d), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(s.id, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, int32(unsafe.Sizeof(DefaultShader_Vertex{})), gl.PtrOffset(0))

	texCoordAttrib := uint32(gl.GetAttribLocation(s.id, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, int32(unsafe.Sizeof(DefaultShader_Vertex{})), gl.PtrOffset(3*4))

	return &DefaultShader_VertexBuffer{id: vao, Size: int32(len(d))}
}

func (vbo *DefaultShader_VertexBuffer) Activate() {
	gl.BindVertexArray(vbo.id)
}
