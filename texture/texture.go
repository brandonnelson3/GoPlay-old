package texture

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Texture struct {
	id uint32
}

func New(file string) (Texture, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return Texture{0}, err
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return Texture{0}, err
	}
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return Texture{0}, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	var id uint32
	gl.GenTextures(1, &id)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, id)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return Texture{id: id}, nil
}

func (t *Texture) Bind(slot uint32) {
	gl.ActiveTexture(slot)
	gl.BindTexture(gl.TEXTURE_2D, t.id)
}
