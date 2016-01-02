package gfx

import (
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/go-gl/gl/v4.1-core/gl"
)

var ioutilReadFile = ioutil.ReadFile

type shader uint32
type shaderManager struct {
	mu      sync.RWMutex
	shaders map[string]shader
}

func NewShaderManager() *shaderManager {
	return &shaderManager{shaders: make(map[string]shader)}
}

func (m shaderManager) Load(filepath string, shaderType uint32) (shader, error) {
	m.mu.RLock()
	s, ok := m.shaders[filepath]
	m.mu.RUnlock()
	if ok {
		return s, nil
	}
	shaderSrc, err := ioutilReadFile(filepath)
	if err != nil {
		return shader(0), err
	}
	// Shaders require a trailing null terminator.
	shaderSrc = append(shaderSrc, byte(0))
	shaderId := gl.CreateShader(shaderType)
	cSrc := gl.Str(string(shaderSrc))
	gl.ShaderSource(shaderId, 1, &cSrc, nil)
	gl.CompileShader(shaderId)
	var status int32
	gl.GetShaderiv(shaderId, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shaderId, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shaderId, logLength, nil, gl.Str(log))
		return shader(0), fmt.Errorf("failed to compile %v: %v", cSrc, log)
	}
	result := shader(shaderId)
	m.mu.Lock()
	m.shaders[filepath] = result
	m.mu.Unlock()
	return result, nil
}

type program uint32

func Program(shaders ...shader) (program, error) {
	p := gl.CreateProgram()
	for _, s := range shaders {
		gl.AttachShader(p, uint32(s))
	}
	gl.LinkProgram(p)
	var status int32
	gl.GetProgramiv(p, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(p, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(p, logLength, nil, gl.Str(log))

		return program(0), fmt.Errorf("failed to link program: %v", log)
	}
	return program(p), nil
}

func (p program) Activate() {
	gl.UseProgram(uint32(p))
}
