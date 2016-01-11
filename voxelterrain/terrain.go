package voxelterrain

import (
	"time"

	"github.com/aquilax/go-perlin"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"sync"

	"fmt"

	"github.com/brandonnelson3/GoPlay/camera"
	"github.com/brandonnelson3/GoPlay/shaders"
	"github.com/brandonnelson3/GoPlay/texture"
	"github.com/brandonnelson3/GoPlay/window"
)

const (
	cellsize     = 32
	cellsizep1   = cellsize + 1
	cellsizep1_2 = cellsizep1 * cellsizep1
	cellsizep1_3 = cellsizep1 * cellsizep1 * cellsizep1

	worldSize   = 6
	worldSizem1 = worldSize - 1
	worldSizet2 = worldSize * 2
	worldTotal  = (2 * worldSize) * (2 * worldSize) * (2 * worldSize)
)

var (
	halfCell = mgl32.Vec3{cellsize / 2, cellsize / 2, cellsize / 2}
)

type cell struct {
	verts []shaders.DefaultShader_Vertex
	vbo   *shaders.DefaultShader_VertexBuffer

	data [cellsizep1_3]byte
	id   cellid
}

type terrain struct {
	shader  *shaders.DefaultShader
	texture texture.Texture

	mu    sync.Mutex
	world map[cellid]*cell
}

func idx(x, y, z int32) int32 {
	return (x*cellsizep1_2 + y*cellsizep1 + z)
}

func (c *cell) generate() {
	noise := perlin.NewPerlin(2, 2, 3, int64(0))
	for x := int32(0); x < cellsizep1; x++ {
		for z := int32(0); z < cellsizep1; z++ {
			h := ((noise.Noise2D(float64(c.id.x*cellsize+x)/10.0, float64(c.id.z*cellsize+z)/10.0) + 1) / 2.0) * 20
			for y := int32(0); y < cellsizep1; y++ {
				index := idx(x, y, z)

				if float64(c.id.y*cellsize+y) < h {
					//if y == 0 {
					c.data[index] = 1
				} else {
					c.data[index] = 0
				}
			}
		}
	}
}

func (c *cell) polygonize(s *shaders.DefaultShader) {
	verts := []shaders.DefaultShader_Vertex{}
	for x := int32(0); x < cellsize; x++ {
		for y := int32(0); y < cellsize; y++ {
			for z := int32(0); z < cellsize; z++ {
				fx := float32(x)
				fy := float32(y)
				fz := float32(z)

				index := idx(x, y, z)

				index1x := idx(x+1, y, z)
				index1y := idx(x, y+1, z)
				index1z := idx(x, y, z+1)

				if c.data[index] == 0 && c.data[index1x] != 0 || c.data[index] != 0 && c.data[index1x] == 0 {
					verts = append(verts,
						shaders.DefaultShader_Vertex{mgl32.Vec3{fx, fy, fz}, mgl32.Vec2{0.0, 0.0}},
						shaders.DefaultShader_Vertex{mgl32.Vec3{fx, fy, fz - 1}, mgl32.Vec2{0.0, 1.0}},
						shaders.DefaultShader_Vertex{mgl32.Vec3{fx, fy - 1, fz}, mgl32.Vec2{1.0, 0.0}},
						shaders.DefaultShader_Vertex{mgl32.Vec3{fx, fy - 1, fz}, mgl32.Vec2{1.0, 0.0}},
						shaders.DefaultShader_Vertex{mgl32.Vec3{fx, fy, fz - 1}, mgl32.Vec2{0.0, 1.0}},
						shaders.DefaultShader_Vertex{mgl32.Vec3{fx, fy - 1, fz - 1}, mgl32.Vec2{1.0, 1.0}})
				}

				if c.data[index] == 0 && c.data[index1y] != 0 || c.data[index] != 0 && c.data[index1y] == 0 {
					verts = append(verts,
						shaders.DefaultShader_Vertex{mgl32.Vec3{fx, fy, fz}, mgl32.Vec2{0.0, 0.0}},
						shaders.DefaultShader_Vertex{mgl32.Vec3{fx, fy, fz - 1}, mgl32.Vec2{0.0, 1.0}},
						shaders.DefaultShader_Vertex{mgl32.Vec3{fx - 1, fy, fz}, mgl32.Vec2{1.0, 0.0}},
						shaders.DefaultShader_Vertex{mgl32.Vec3{fx - 1, fy, fz}, mgl32.Vec2{1.0, 0.0}},
						shaders.DefaultShader_Vertex{mgl32.Vec3{fx, fy, fz - 1}, mgl32.Vec2{0.0, 1.0}},
						shaders.DefaultShader_Vertex{mgl32.Vec3{fx - 1, fy, fz - 1}, mgl32.Vec2{1.0, 1.0}})
				}

				if c.data[index] == 0 && c.data[index1z] != 0 || c.data[index] != 0 && c.data[index1z] == 0 {
					verts = append(verts,
						shaders.DefaultShader_Vertex{mgl32.Vec3{fx, fy, fz}, mgl32.Vec2{0.0, 0.0}},
						shaders.DefaultShader_Vertex{mgl32.Vec3{fx, fy - 1, fz}, mgl32.Vec2{0.0, 1.0}},
						shaders.DefaultShader_Vertex{mgl32.Vec3{fx - 1, fy, fz}, mgl32.Vec2{1.0, 0.0}},
						shaders.DefaultShader_Vertex{mgl32.Vec3{fx - 1, fy, fz}, mgl32.Vec2{1.0, 0.0}},
						shaders.DefaultShader_Vertex{mgl32.Vec3{fx, fy - 1, fz}, mgl32.Vec2{0.0, 1.0}},
						shaders.DefaultShader_Vertex{mgl32.Vec3{fx - 1, fy - 1, fz}, mgl32.Vec2{1.0, 1.0}})
				}
			}
		}
	}
	if len(verts) == 0 {
		return
	}
	c.verts = verts
}

type cellid struct {
	x, y, z int32
}

func (lhs *cellid) Equal(rhs cellid) bool {
	return lhs.x == rhs.x && lhs.y == rhs.y && lhs.z == rhs.z
}

func NewCell(s *shaders.DefaultShader, id cellid) *cell {
	cell := &cell{id: id}
	cell.generate()
	cell.polygonize(s)
	return cell
}

func isCellInWorld(cell, centroidCell cellid) bool {
	if cell.x < centroidCell.x-worldSizem1 {
		return false
	}
	if cell.y < centroidCell.y-worldSizem1 {
		return false
	}
	if cell.z < centroidCell.z-worldSizem1 {
		return false
	}

	if cell.x > centroidCell.x+worldSize {
		return false
	}
	if cell.y > centroidCell.y+worldSize {
		return false
	}
	if cell.z > centroidCell.z+worldSize {
		return false
	}
	return true
}

func (t *terrain) generate(x, y, z int32) {
	lastCell := cellid{-1, -1, -1}
	for {
		// No point in checking more often then every 100ms.
		<-time.After(100 * time.Millisecond)

		// Positions are shifted by half a cell from cell positions since cell positions are in the lower left corner.
		pos := camera.C.GetPosition().Sub(halfCell)

		// If this is the same cell as last iteration bail.
		thisCell := cellid{int32(pos.X())/cellsize + x, int32(pos.Y())/cellsize + y, int32(pos.Z())/cellsize + z}
		if lastCell.Equal(thisCell) {
			continue
		}
		lastCell = thisCell

		// If this cell is already present in the world bail.
		t.mu.Lock()
		_, newOk := t.world[thisCell]
		t.mu.Unlock()
		if newOk {
			continue
		}

		// This is a new cell not currently present in the world. Generate then insert.
		c := NewCell(t.shader, thisCell)
		t.mu.Lock()
		t.world[thisCell] = c
		t.mu.Unlock()
	}
}

func NewTerrain() (*terrain, error) {
	shader, err := shaders.NewDefaultShader()
	if err != nil {
		return nil, err
	}
	shader.Activate()
	// Load the texture
	texture, err := texture.New("assets/crate.jpg")
	if err != nil {
		return nil, err
	}
	t := &terrain{shader: shader, texture: texture, world: make(map[cellid]*cell)}

	for x := int32(1 - worldSize); x <= worldSize; x++ {
		for y := int32(1 - worldSize); y <= worldSize; y++ {
			for z := int32(1 - worldSize); z <= worldSize; z++ {
				go t.generate(x, y, z)
			}
		}
	}
	return t, nil
}

func (t *terrain) Render() {
	t.shader.Activate()
	t.shader.SetProjection(mgl32.Perspective(mgl32.DegToRad(camera.C.FOVDegrees), float32(window.M.Width)/float32(window.M.Height), camera.C.NearPlaneDist, camera.C.FarPlaneDist))
	t.shader.SetView(camera.C.GetViewMatrix())
	t.texture.Bind(gl.TEXTURE0)
	pos := camera.C.GetPosition()
	pos = pos.Sub(halfCell)
	centroidCell := cellid{int32(pos.X()) / cellsize, int32(pos.Y()) / cellsize, int32(pos.Z()) / cellsize}
	t.mu.Lock()
	defer t.mu.Unlock()
	if size := len(t.world); size > worldTotal {
		fmt.Printf("WARNING: Scene has %d cells should be %v\n", len(t.world), worldTotal)
	}
	for id, c := range t.world {
		if !isCellInWorld(c.id, centroidCell) {
			delete(t.world, c.id)
			continue
		}
		if c.verts == nil {
			continue
		}
		if c.vbo == nil {
			c.vbo = shaders.NewDefaultShader_VertexBuffer(t.shader, c.verts)
		}
		c.vbo.Activate()
		t.shader.SetModel(mgl32.Translate3D(float32(id.x*cellsize), float32(id.y*cellsize), float32(id.z*cellsize)))
		gl.DrawArrays(gl.TRIANGLES, 0, c.vbo.Size)
	}
}
