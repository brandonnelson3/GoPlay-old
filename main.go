package main

import (
	"fmt"
	_ "image/jpeg"
	"log"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"

	"sync/atomic"

	"github.com/brandonnelson3/GoPlay/gameobjects"
	"github.com/brandonnelson3/GoPlay/input"
)

const windowWidth = 800
const windowHeight = 600

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

var fps uint32

func printFPS() {
	log.Printf("FPS is currently %d/second", atomic.SwapUint32(&fps, 0))
	time.AfterFunc(1*time.Second, printFPS)
}

func main() {
	go printFPS()
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Cube", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	inputmanager := input.NewManager(window)

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	cube, err := gameobjects.NewCube()
	if err != nil {
		panic(err)
	}

	previousTime := glfw.GetTime()

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		inputmanager.RunKeys()

		// Update
		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time

		cube.Update(elapsed)
		cube.Render()

		atomic.AddUint32(&fps, 1)

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
