package main

import (
	"image"
	"image/gif"
	"image/draw"
	"image/color"
	"fmt"
	"net/http"
	"io"
	"bytes"
	"os"
)

// Skipping factor
type MazeImage struct {
	height 		int
	width 		int
	padding 	int     // only want square path so wall & path width should be the same and no factor needed
	wall_color  [3]byte
	path_color  [3]byte
}

func NewMazeImage(height, width, padding int) *MazeImage {
	return &MazeImage{
		height: height,
		width: width,
		padding: padding,
		path_color: [...]byte{255, 255, 255},
	}
}

func (m *MazeImage) SetPathColor(r, g, b byte){
	m.path_color = [3]byte{r,g,b}
}

func (m *MazeImage) SetWallColor(r, g, b byte){
	m.wall_color = [3]byte{r,g,b}
}

func (m MazeImage) String() string {
	return fmt.Sprintf(
		"http://www.hereandabove.com/cgi-bin/maze?%d+%d+%d+%d+0+%d+%d+%d+%d+%d+%d",
		m.width,
		m.height,
		m.padding,
		m.padding,
		m.wall_color[0],
		m.wall_color[1],
		m.wall_color[2],
		m.path_color[0],
		m.path_color[1],
		m.path_color[2],
	)
}

func (m *MazeImage) GetImage() (*ImageFile, error) {
	resp, err := http.Get(m.String())
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	} else {
		if image, err  := NewImageFile(resp.Body, m); err != nil {
			return nil, err
		} else {
			return image, nil
		}
	}
}

type ImageFile struct {
	m [][]*color.RGBA
	i *MazeImage
	x int
	y int
}

func (i ImageFile) Print() {

	if i.i.padding <= 1 {
		for _, fields := range i.m {
			for _, field := range fields {
				if field.R == i.i.path_color[0] && field.G == i.i.path_color[1] && field.B == i.i.path_color[2] {
					fmt.Print(" ")
				} else {
					fmt.Print("#")
				}
			}
			fmt.Print("\n")
		}
	} else {
		for y := 0; y < len(i.m); y += i.i.padding {
			for x := 0; x < len(i.m[y]); x += i.i.padding {
				if i.m[y][x].R == i.i.path_color[0] && i.m[y][x].G == i.i.path_color[1] && i.m[y][x].B == i.i.path_color[2] {
					fmt.Print(" ")
				} else {
					fmt.Print("#")
				}
			}
			fmt.Print("\n")
		}
	}
}

func (i *ImageFile) Up() {
	if i.y - i.i.padding >= 0 {
		i.y -= i.i.padding
	}
}

func (i *ImageFile) Down() {
	if i.y + i.i.padding <= len(i.m) {
		i.y += i.i.padding
	}
}

func (i ImageFile) Save(file *os.File) error {

	rect := image.Rect(0, 0, len(i.m[0]), len(i.m))
	rgba := image.NewRGBA(rect)

	for y, pixels := range i.m {
		for x, c := range pixels {
			rgba.Set(x, y, c)
		}
	}

	return gif.Encode(file, rgba, &gif.Options{NumColors: 256})
}

func NewImageFile(r io.Reader, m *MazeImage) (*ImageFile, error) {
	img, err := gif.Decode(r)

	if err != nil {
		return nil, err
	}

	rect := img.Bounds()
	rgba := image.NewRGBA(rect)
	draw.Draw(rgba, rect, img, rect.Min, draw.Src)
	reader := bytes.NewReader(rgba.Pix)
	ret := make([][]*color.RGBA, rect.Max.Y)

	for i := 0; i < rect.Max.Y; i++ {

		if len(ret[i]) == 0 {
			ret[i] = make([]*color.RGBA, rect.Max.X)
		}

		for c := 0; c < rect.Max.X; c ++ {
			part := make([]byte, 4)
			reader.Read(part)
			ret[i][c] = &color.RGBA{part[0],part[1],part[2],part[3]}
		}
	}

	return &ImageFile{ret, m, 0, 0}, nil;
}

func main(){

	maze := NewMazeImage(10, 20, 5)
	file, _ := maze.GetImage()

	file.Print()

	f, _ := os.Create("foo.gif")
	defer f.Close()
	file.Save(f)
}
