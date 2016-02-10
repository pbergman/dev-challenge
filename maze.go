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
	"sync"
)

// MazeImage is basic configuration for fetching a maze image
type MazeImage struct {
	height 		int
	width 		int
	// ration i replacement for path and wall width
	// because we only want square paths/walls this
	// also drops need for factor length. This will
	// also be ignored when fetching image and only
	// used with re-rendering of image
	ratio 		uint
	wall_color  *color.RGBA
	path_color  *color.RGBA
}

func NewMazeImage(height, width int) *MazeImage {
	return &MazeImage{
		height: 	height,
		width: 		width,
		ratio: 		1,
		path_color: &color.RGBA{255,255,255,255},
		wall_color: &color.RGBA{0,0,0,255},
	}
}

// SetPathColor will set path color for maze, default 255,255,255
func (m *MazeImage) SetPathColor(r, g, b byte){
	m.path_color.R = r
	m.path_color.G = g
	m.path_color.B = b
}

// SetWallColor will set path color for maze, default 0,0,0
func (m *MazeImage) SetWallColor(r, g, b byte){
	m.wall_color.R = r
	m.wall_color.G = g
	m.wall_color.B = b
}

// Will set ration, so ef set to 2 every "pixel block" is 2x2 pixel
func (m *MazeImage) SetRatio(r uint) {
	m.ratio = r
}

// String will return url for fetching the maze image
func (m MazeImage) String() string {
	return fmt.Sprintf(
		"http://www.hereandabove.com/cgi-bin/maze?%d+%d+%d+%d+0+%d+%d+%d+%d+%d+%d",
		m.width,
		m.height,
		1,
		1,
		m.wall_color.R,
		m.wall_color.G,
		m.wall_color.B,
		m.path_color.R,
		m.path_color.G,
		m.path_color.B,
	)
}

// GetMatrix will return and create matrix based on fetched image
func (m *MazeImage) GetMatrix() (*MazeMatrix, error) {
	resp, err := http.Get(m.String())
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	} else {
		if image, err  := NewMazeMatrix(resp.Body, m); err != nil {
			return nil, err
		} else {
			return image, nil
		}
	}
}

type MatrixToken uint64
type Direction uint64

func (d Direction) String() string {
	switch(d) {
	case LEFT:
		return "LEFT"
	case UP:
		return "UP"
	case RIGHT:
		return "RIGHT"
	case DOWN:
		return "DOWN"
	default:
		return "UNKNOWN"
	}

}

const (
	WALL MatrixToken = 1 << iota
	PATH
	BORDER
	START
	END
	VISITED
	ROUTE

	LEFT Direction = 1 <<iota
	UP
	RIGHT
	DOWN
)

//
//		|
// Y	|
//		|
// 		-----------
//
//			X
//

type MazeMatrix struct {
	m [][]MatrixToken
	i *MazeImage
}

func (i MazeMatrix) Has(x, y int, t MatrixToken) bool {
	return t == (t & i.m[y][x])
}

// String prints the the matrix to the stdout in visula way
func (i MazeMatrix) String() string {
	buff := new(bytes.Buffer)
	for _, data := range i.m {
		for _, token := range data {
			switch(true) {
			case START == (START & token):
				buff.Write([]byte{'S'})
			case ROUTE == (ROUTE & token):
				buff.Write([]byte{'*'})
			case VISITED == (VISITED & token):
				buff.Write([]byte{'.'})
			case END == (END & token):
				buff.Write([]byte{'E'})
			case WALL == (WALL & token):
				buff.Write([]byte{'#'})
			case PATH == (PATH & token), BORDER == (BORDER & token):
				buff.Write([]byte{' '})
			}
		}
		buff.Write([]byte{'\n'})
	}
	return string(buff.Bytes())
}

// isWall check if the given colors is matching the config wallcolors
func (i *MazeMatrix) isWall(r, g, b uint8) bool {
	return i.i.wall_color.R == r && i.i.wall_color.G == g && i.i.wall_color.B == b;
}

// NewMazeMatrix creates a image matrix base on fetched body
func NewMazeMatrix(r io.Reader, m *MazeImage) (*MazeMatrix, error) {
	mazeMatrix := &MazeMatrix{i: m}
	img, err := gif.Decode(r)
	if err != nil {
		return nil, err
	}
	rect := img.Bounds()
	rgba := image.NewRGBA(rect)
	draw.Draw(rgba, rect, img, rect.Min, draw.Src)
	reader := bytes.NewReader(rgba.Pix)
	matrix := make([][]MatrixToken, rect.Max.Y)
	for y := 0; y < rect.Max.Y; y++ {
		for x := 0; x < rect.Max.X; x ++ {
			if len(matrix[y]) == 0 {
				matrix[y] = make([]MatrixToken, rect.Max.X)
			}
			part := make([]byte, 4)
			reader.Read(part)
			if y == 0 || x == 0 {
				matrix[y][x] = BORDER
			} else {
				if (mazeMatrix.isWall(part[0],part[1],part[2])) {
					matrix[y][x] = WALL
				} else {
					matrix[y][x] = PATH
				}
			}
		}
	}
	mazeMatrix.m = matrix
	return mazeMatrix, nil;
}

// ToFile will push the drawing created with DrawImage to the given fd
func (i MazeMatrix) ToFile(file *os.File) error {
	return gif.Encode(file, i.DrawImage(), &gif.Options{NumColors: 256})
}

// DrawImage draws a new image beased on matrix and config ration
func (i MazeMatrix) DrawImage() draw.Image {
	rect := image.Rect(0, 0, len(i.m[0]) * int(i.i.ratio), len(i.m)  * int(i.i.ratio))
	rgba := image.NewRGBA(rect)
	draw.Draw(rgba, rgba.Bounds(), &image.Uniform{color.White}, image.ZP, draw.Src)
	for y := 0; y < len(i.m) * int(i.i.ratio); y += int(i.i.ratio) {
		for x := 0; x < len(i.m[y/int(i.i.ratio)]) * int(i.i.ratio); x += int(i.i.ratio) {
			switch(i.m[y/int(i.i.ratio)][x/int(i.i.ratio)]) {
			case WALL:
				draw.Draw(rgba, image.Rect(x, y, x + int(i.i.ratio), y + int(i.i.ratio)), &image.Uniform{i.i.wall_color}, image.ZP, draw.Src)
			case PATH, BORDER:
				draw.Draw(rgba, image.Rect(x, y, x + int(i.i.ratio), y + int(i.i.ratio)), &image.Uniform{i.i.path_color}, image.ZP, draw.Src)
			}
		}
	}
	// rgba.Set()
	return rgba
}

type Walker struct {
	s   struct{			// start point
		  y int
		  x int
	  }
	e 	struct{			// end point
		  y int
		  x int
	  }
	b 	image.Rectangle	// maze bounds, to set borders
	m   *MazeMatrix
}

func NewWalker(m *MazeMatrix) *Walker {

	w := &Walker{
		b: image.Rect(1, 1, len(m.m[0]) - 2, len(m.m) - 2),
		m: m,
	}

	done := false
	check := func(p struct{ y, x int},walker *Walker) bool {
		if walker.m.Has(p.x, p.y, PATH) {
			if walker.s.x == 0 && walker.s.y == 0 {
				walker.s.y, walker.s.x = p.y, p.x
				walker.m.m[p.y][p.x] |= START
				return false
			}
			if walker.e.x == 0 && walker.e.y == 0 {
				walker.e.y, walker.e.x = p.y, p.x
				walker.m.m[p.y][p.x] |= END
				return true
			}
		}
		return false
	}

	pos := struct{
		y int
		x int
	}{1,1}

	for w.right(&pos.x) {
		done = check(pos, w)
	}

	for !done && w.down(&pos.y) {
		done = check(pos, w)
	}

	for !done && w.left(&pos.x) {
		done = check(pos, w)
	}

	for !done && w.up(&pos.y) {
		done = check(pos, w)
	}

	// add starting point as visited
	w.m.m[w.s.y][w.s.x] |= VISITED

	return w

}



func (w *Walker) Solve() {
	var wg sync.WaitGroup

	type solver struct {
		p    struct{ x, y int}
		t  []struct{ x, y int}
	}

	route := make(chan *solver, 1)
	solvers := make(chan *solver, 10)
	solvers <- &solver{p: struct{x,y int}{w.s.x, w.s.y}, t: []struct{x,y int}{struct{x,y int}{w.s.x, w.s.y}}}
	finished := make(chan bool, 10)


	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(jobs chan *solver, solvedRoute chan <- *solver, finished chan bool) {
			outer_loop:
			for {
				select {
				case <- finished:
					break outer_loop
				case job := <- jobs:
					for {

						free := w.getFree(job.p.x, job.p.y);

						if len(free) == 0 {
							break
						}

						//switch (len(free)) {
						//case 0:
						//	break
						//default:
							for i := 1; i < len(free); i++ {

								ns := &solver{p: struct{x,y int}{job.p.x, job.p.y}, t: job.t}

								switch(free[i]) {
								case RIGHT:
									w.right(&ns.p.x)
								case DOWN:
									w.down(&ns.p.y)
								case UP:
									w.up(&ns.p.y)
								case LEFT:
									w.left(&ns.p.x)
								}

								ns.t = append(ns.t, ns.p)
								w.m.m[ns.p.y][ns.p.x] |= VISITED
fmt.Println(len(jobs), len(finished))
								//if END == (END & w.m.m[ns.p.y][ns.p.x]) {
								////	solvedRoute <- ns
								////	for i := 0; i < 10; i++ {
								////		finished <- true
								////	}
								////	break outer_loop
								////}

								jobs <- ns
							}


							switch(free[0]) {
							case RIGHT:
								w.right(&job.p.x)
							case DOWN:
								w.down(&job.p.y)
							case UP:
								w.up(&job.p.y)
							case LEFT:
								w.left(&job.p.x)
							}

							job.t = append(job.t, job.p)
							w.m.m[job.p.y][job.p.x] |= VISITED

							if END == (END & w.m.m[job.p.y][job.p.x]) {
								solvedRoute <- job
								for i := 0; i < 10; i++ {
									finished <- true
								}
								break outer_loop
							}

						//}
						//if free := w.getFree(job.p.x, job.p.y); len(free) == 0 {
						//	break
						//} else {
						//
						//	for i := 1; i < len(free); i++ {
						//		jobs <- &solver{p: struct{x,y int}{job.p.x, job.p.y}, t: job.t}
						//	}
						//
						//	switch(free[0]) {
						//	case RIGHT:
						//		w.right(&job.p.x)
						//	case DOWN:
						//		w.down(&job.p.y)
						//	case UP:
						//		w.up(&job.p.y)
						//	case LEFT:
						//		w.left(&job.p.x)
						//	}
						//
						//	//job.t = append(job.t, struct{x,y int}{job.p.x, job.p.y})
						//	w.m.m[job.p.y][job.p.x] |= VISITED
						//
						//	if END == (END & w.m.m[job.p.y][job.p.x]) {
						//		solvedRoute <- job
						//		for i := 0; i < 10; i++ {
						//			finished <- true
						//		}
						//		break outer_loop
						//	}
						//}
					}
				}
			}
			wg.Done()
		}(solvers, route, finished)
	}

	wg.Wait()

	close(route)
	close(solvers)
	close(finished)

	r := <- route

	for _, t := range r.t {
		fmt.Printf("%#v", t)
		w.m.m[t.y][t.x] |= ROUTE

	}


	//for _, s := range r {
	//	w.m.m[s.y][s.x] |= ROUTE
	//}

}

func (m *Walker) getFree(x, y int) []Direction {

	directions := make([]Direction, 0)

	if x > m.b.Min.X && m.m.Has(x - 1, y, PATH) && !m.m.Has(x - 1, y, VISITED) {
		directions = append(directions, LEFT)
	}

	if x < m.b.Max.X && m.m.Has(x + 1, y, PATH) && !m.m.Has(x + 1, y, VISITED) {
		directions = append(directions, RIGHT)
	}

	if y > m.b.Min.Y && m.m.Has(x, y - 1, PATH) && !m.m.Has(x, y - 1, VISITED) {
		directions = append(directions, UP)
	}

	if y < m.b.Max.Y && m.m.Has(x, y + 1, PATH) && !m.m.Has(x, y + 1, VISITED) {
		directions = append(directions, DOWN)
	}

	return directions
}

func (m *Walker) left(x *int) bool {
	if *x > m.b.Min.X {
		*x--
		return true
	}
	return false
}

func (m *Walker) right(x *int) bool {
	if *x < m.b.Max.X {
		*x++
		return true
	}
	return false
}

func (m *Walker) up(y *int) bool {
	if *y > m.b.Min.Y {
		*y--
		return true
	}
	return false
}

func (m *Walker) down(y *int) bool {
	if *y < m.b.Max.Y {
		*y++
		return true
	}
	return false
}

func main(){

	maze := NewMazeImage(10, 20)
	maze.SetRatio(10)
	matrix, _  := maze.GetMatrix()
	fmt.Println(matrix)

	f, _ := os.Create("foo.gif")
	defer f.Close()
	matrix.ToFile(f)
	walker := NewWalker(matrix)

	walker.Solve()
	fmt.Printf("%#v\n", walker)
	fmt.Println(matrix)
	fmt.Println(walker.getFree(walker.s.x, walker.s.y))

}
