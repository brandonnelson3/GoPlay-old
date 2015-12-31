package gfx

import (
	"errors"
	"io/ioutil"
	"testing"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

func TestMain(m *testing.M) {
	glfw.Init()
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
}

func TestLoadShader(t *testing.T) {
	testPrepOpenGL(t)
	defer func() { ioutilReadFile = ioutil.ReadFile }()
	ioutilReadFile = func(string) ([]byte, error) {
		return []byte(`#version 330
uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

in vec3 vert;
in vec2 vertTexCoord;

out vec2 fragTexCoord;

void main() {
    fragTexCoord = vertTexCoord;
    gl_Position = projection * camera * model * vec4(vert, 1);
}`), nil
	}

	sMgr := NewShaderManager()
	s1, err := sMgr.Load("filename", gl.VERTEX_SHADER)
	if err != nil {
		t.Fatalf("Got error while compiling valid shader: %v", err)
	}

	s2, err := sMgr.Load("filename", gl.VERTEX_SHADER)
	if err != nil {
		t.Fatalf("Got error while compiling valid shader: %v", err)
	}

	if s1 != s2 {
		t.Errorf("Got different shader value on second load than first, First: %v, Second: %v", s1, s2)
	}
}

func TestLoadShader_ReadFileError(t *testing.T) {
	testPrepOpenGL(t)
	defer func() { ioutilReadFile = ioutil.ReadFile }()
	ioutilReadFile = func(string) ([]byte, error) {
		return nil, errors.New("Test")
	}

	sMgr := NewShaderManager()
	if _, err := sMgr.Load("filename", gl.VERTEX_SHADER); err == nil {
		t.Fatal("Expected error but didn't receive one")
	}
}

func testPrepOpenGL(t *testing.T) {
	glfw.Init()
	window, err := glfw.CreateWindow(1, 1, "Test", nil, nil)
	if err != nil {
		t.Fatalf("Got error while initializing window: %v", err)
	}
	window.MakeContextCurrent()
	if err := gl.Init(); err != nil {
		t.Fatalf("Got error while initializing OpenGL: %v", err)
	}
}
