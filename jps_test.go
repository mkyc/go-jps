package jps

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestFind_simple(t *testing.T) {
	type args struct {
		obstacles [][]bool
		start     Point
		goal      Point
	}
	tests := []struct {
		name    string
		args    args
		want    []Point
		wantErr bool
	}{
		{
			name: "1_simple",
			args: args{
				obstacles: [][]bool{
					{false, false, false},
					{false, true, false},
					{false, true, false},
					{false, true, false},
					{false, false, false},
				},
				start: Point{0, 0},
				goal:  Point{4, 2},
			},
			want: []Point{
				{0, 0},
				{1, 0},
				{2, 0},
				{3, 0},
				{4, 0},
				{4, 1},
				{4, 2},
			},
			wantErr: false,
		},
		{
			name: "2_direct",
			args: args{
				obstacles: [][]bool{
					{false, false, false},
					{false, true, false},
					{false, true, false},
					{false, true, false},
					{false, false, false},
				},
				start: Point{0, 0},
				goal:  Point{4, 1},
			},
			want: []Point{
				{0, 0},
				{1, 0},
				{2, 0},
				{3, 0},
				{4, 1},
			},
			wantErr: false,
		},
		{
			name: "3_simple_with_frame",
			args: args{
				obstacles: [][]bool{
					{true, true, true, true, true, true, true},
					{true, true, true, true, true, true, true},
					{true, true, false, false, false, true, true},
					{true, true, false, false, false, true, true},
					{true, true, false, false, false, true, true},
					{true, true, true, true, true, true, true},
					{true, true, true, true, true, true, true},
				},
				start: Point{2, 2},
				goal:  Point{4, 4},
			},
			want: []Point{
				{2, 2},
				{3, 3},
				{4, 4},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Find(tt.args.obstacles, tt.args.start, tt.args.goal)
			generateImage(t, tt.args.obstacles, tt.args.start, tt.args.goal, got, fmt.Sprintf("./test_output/Find_simple_%s_result.png", tt.name))
			if (err != nil) != tt.wantErr {
				t.Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("Find() got = %v", got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Find() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFind_maps(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		start    Point
		goal     Point
		wantErr  bool
	}{
		{
			name:     "1_jps",
			filename: "./test_assets/map1.png",
			start:    Point{160, 140},
			goal:     Point{170, 170},
			wantErr:  false,
		},
		{
			name:     "2_jps",
			filename: "./test_assets/map1.png",
			start:    Point{280, 70},
			goal:     Point{390, 265},
			wantErr:  false,
		},
		{
			name:     "3_direct",
			filename: "./test_assets/map1.png",
			start:    Point{140, 195},
			goal:     Point{250, 300},
			wantErr:  false,
		},
		{
			name:     "4_start_outside_map",
			filename: "./test_assets/map1.png",
			start:    Point{-1, -1},
			goal:     Point{250, 300},
			wantErr:  true,
		},
		{
			name:     "5_goal_outside_map",
			filename: "./test_assets/map1.png",
			start:    Point{140, 195},
			goal:     Point{500, 500},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obstacles := readPNG(t, tt.filename)
			got, err := Find(obstacles, tt.start, tt.goal)
			generateImage(t, obstacles, tt.start, tt.goal, got, fmt.Sprintf("./test_output/Find_maps_%s_result.png", tt.name))
			if (err != nil) != tt.wantErr {
				t.Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("Find() got = %v", got)
		})
	}
}

func TestTryFind_simple(t *testing.T) {
	type args struct {
		obstacles [][]bool
		start     Point
		goal      Point
	}
	tests := []struct {
		name    string
		args    args
		want    []Point
		wantErr bool
	}{
		{
			name: "1_simple",
			args: args{
				obstacles: [][]bool{
					{true, false, false},
					{false, true, false},
					{false, true, false},
					{false, false, false},
					{false, false, false},
				},
				start: Point{0, 0},
				goal:  Point{4, 2},
			},
			want: []Point{
				{1, 0},
				{2, 0},
				{3, 1},
				{4, 2},
			},
			wantErr: false,
		},
		{
			name: "2_direct",
			args: args{
				obstacles: [][]bool{
					{true, false, false},
					{true, true, false},
					{false, true, false},
					{false, true, false},
					{false, false, false},
				},
				start: Point{0, 0},
				goal:  Point{4, 0},
			},
			want: []Point{
				{2, 0},
				{3, 0},
				{4, 0},
			},
			wantErr: false,
		},
		{
			name: "3_simple_with_frame",
			args: args{
				obstacles: [][]bool{
					{true, true, true, true, true, true, true},
					{true, true, true, true, true, true, true},
					{true, true, true, false, false, true, true},
					{true, true, false, false, false, true, true},
					{true, true, false, false, false, true, true},
					{true, true, true, true, true, true, true},
					{true, true, true, true, true, true, true},
				},
				start: Point{2, 2},
				goal:  Point{4, 4},
			},
			want: []Point{
				{3, 3},
				{4, 4},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TryFind(tt.args.obstacles, tt.args.start, tt.args.goal)
			generateImage(t, tt.args.obstacles, tt.args.start, tt.args.goal, got, fmt.Sprintf("./test_output/TryFind_simple_%s_result.png", tt.name))
			if (err != nil) != tt.wantErr {
				t.Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("Find() got = %v", got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Find() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTryFind_maps(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		start    Point
		goal     Point
		wantErr  bool
	}{
		{
			name:     "1_jps",
			filename: "./test_assets/map1.png",
			start:    Point{150, 130},
			goal:     Point{190, 190},
			wantErr:  false,
		},
		{
			name:     "2_jps",
			filename: "./test_assets/map1.png",
			start:    Point{300, 90},
			goal:     Point{400, 275},
			wantErr:  false,
		},
		{
			name:     "3_direct",
			filename: "./test_assets/map1.png",
			start:    Point{120, 175},
			goal:     Point{270, 320},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obstacles := readPNG(t, tt.filename)
			got, err := TryFind(obstacles, tt.start, tt.goal)
			generateImage(t, obstacles, tt.start, tt.goal, got, fmt.Sprintf("./test_output/TryFind_maps_%s_result.png", tt.name))
			if (err != nil) != tt.wantErr {
				t.Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("Find() got = %v", got)
		})
	}
}

func TestMustFind_simple(t *testing.T) {
	type args struct {
		obstacles [][]bool
		start     Point
		goal      Point
	}
	tests := []struct {
		name string
		args args
		want []Point
	}{
		{
			name: "1_simple",
			args: args{
				obstacles: [][]bool{
					{false, false, false},
					{false, true, false},
					{false, true, false},
					{false, true, false},
					{false, false, false},
				},
				start: Point{-1, -1},
				goal:  Point{4, 2},
			},
			want: []Point{
				{0, 0},
				{1, 0},
				{2, 0},
				{3, 0},
				{4, 0},
				{4, 1},
				{4, 2},
			},
		},
		{
			name: "2_direct",
			args: args{
				obstacles: [][]bool{
					{true, false, false},
					{true, true, false},
					{false, true, false},
					{false, true, false},
					{false, false, false},
				},
				start: Point{-1, 0},
				goal:  Point{4, 0},
			},
			want: []Point{
				{2, 0},
				{3, 0},
				{4, 0},
			},
		},
		{
			name: "3_simple_with_frame",
			args: args{
				obstacles: [][]bool{
					{true, true, true, true, true, true, true},
					{true, true, true, true, true, true, true},
					{true, true, false, false, false, true, true},
					{true, true, false, false, false, true, true},
					{true, true, false, false, false, true, true},
					{true, true, true, true, true, true, true},
					{true, true, true, true, true, true, true},
				},
				start: Point{-1, -1},
				goal:  Point{4, 4},
			},
			want: []Point{
				{2, 2},
				{3, 3},
				{4, 4},
			},
		},
		{
			name: "4_no_path_to_goal",
			args: args{
				obstacles: [][]bool{
					{false, false, true, false, false, true, false},
					{false, false, true, false, false, true, false},
					{false, false, true, false, false, true, false},
					{false, false, true, false, false, true, false},
					{false, false, true, false, false, true, false},
					{false, false, true, false, false, true, false},
					{false, false, false, false, false, true, false},
				},
				start: Point{3, 0},
				goal:  Point{3, 6},
			},
			want: []Point{
				{3, 0},
			},
		},
		{
			name: "5_no_path_from_start",
			args: args{
				obstacles: [][]bool{
					{false, true, false, false, false, false, false},
					{false, true, false, false, false, true, false},
					{false, true, false, false, false, true, false},
					{false, true, false, false, false, true, false},
					{false, true, false, false, false, true, false},
					{false, true, false, false, false, true, false},
					{false, true, false, false, false, true, false},
				},
				start: Point{3, 0},
				goal:  Point{3, 6},
			},
			want: []Point{
				{3, 0},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := MustFind(tt.args.obstacles, tt.args.start, tt.args.goal)
			generateImage(t, tt.args.obstacles, tt.args.start, tt.args.goal, got, fmt.Sprintf("./test_output/MustFind_simple_%s_result.png", tt.name))
			t.Logf("Find() got = %v", got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Find() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustFind_maps(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		start    Point
		goal     Point
		wantPath bool
	}{
		{
			name:     "1_jps",
			filename: "./test_assets/map1.png",
			start:    Point{150, -10},
			goal:     Point{190, 190},
			wantPath: true,
		},
		{
			name:     "2_jps",
			filename: "./test_assets/map1.png",
			start:    Point{300, 90},
			goal:     Point{500, 400},
			wantPath: true,
		},
		{
			name:     "3_jps",
			filename: "./test_assets/map1.png",
			start:    Point{120, 20},
			goal:     Point{280, 330},
			wantPath: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obstacles := readPNG(t, tt.filename)
			got := MustFind(obstacles, tt.start, tt.goal)
			generateImage(t, obstacles, tt.start, tt.goal, got, fmt.Sprintf("./test_output/MustFind_maps_%s_result.png", tt.name))
			t.Logf("Find() got = %v", got)
			if tt.wantPath && len(got) == 1 {
				t.Errorf("Find() got = %v, want non empty path", got)
			}
			if !tt.wantPath && len(got) > 1 {
				t.Errorf("Find() got = %v, want empty path", got)
			}
		})
	}
}

func FuzzFind(f *testing.F) {
	obstacles := readPNG(f, "./test_assets/passable_map1.png")
	f.Fuzz(func(t *testing.T, startX, startY, goalX, goalY int) {
		t.Logf("startX: %d, startY: %d, goalX: %d, goalY: %d", startX, startY, goalX, goalY)
		start := Point{startX % len(obstacles), startY % len(obstacles[0])}
		goal := Point{goalX % len(obstacles), goalY % len(obstacles[0])}
		_, _ = Find(obstacles, start, goal)
	})
}

func BenchmarkFind(b *testing.B) {
	obstacles := readPNG(b, "./test_assets/map1.png")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Find(obstacles, Point{280, 70}, Point{390, 265})
	}
}

type TB interface {
	Fatal(args ...interface{})
}

func readPNG(tb TB, filename string) [][]bool {
	file, err := os.Open(filename)
	if err != nil {
		tb.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			tb.Fatal(err)
		}
	}(file)

	img, err := png.Decode(file)
	if err != nil {
		tb.Fatal(err)
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	bwArray := make([][]bool, width)
	for i := range bwArray {
		bwArray[i] = make([]bool, height)
	}

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// Convert to 8-bit grayscale
			grayscaleValue := uint8((r + g + b) / 3 >> 8)
			bwArray[x][y] = grayscaleValue < 120
		}
	}
	return bwArray
}

func generateImage(tb TB, obstacles [][]bool, start, goal Point, path []Point, filePath string) {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		tb.Fatal(err)
	}
	height := len(obstacles[0])
	width := len(obstacles)

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.White)
		}
	}

	for x, row := range obstacles {
		for y, isObstacle := range row {
			if isObstacle {
				img.Set(x, y, color.Black)
			}
		}
	}

	//blue path
	for _, p := range path {
		img.Set(p.X, p.Y, color.RGBA{0, 0, 255, 255})
	}

	//yellow corners
	img.Set(0, 0, color.RGBA{255, 255, 0, 255})
	img.Set(0, height-1, color.RGBA{255, 255, 0, 255})
	img.Set(width-1, 0, color.RGBA{255, 255, 0, 255})
	img.Set(width-1, height-1, color.RGBA{255, 255, 0, 255})

	if isPointPassable(obstacles, start) {
		img.Set(start.X, start.Y, color.RGBA{0, 255, 0, 255}) // green for start
	} else if !isPointOutsideMap(obstacles, start) && isPointInsideObstacle(obstacles, start) {
		img.Set(start.X, start.Y, color.RGBA{0, 128, 0, 255}) // dark green for obstacle
	}
	if isPointPassable(obstacles, goal) {
		img.Set(goal.X, goal.Y, color.RGBA{255, 0, 0, 255}) // red for goal
	} else if !isPointOutsideMap(obstacles, goal) && isPointInsideObstacle(obstacles, goal) {
		img.Set(goal.X, goal.Y, color.RGBA{128, 0, 0, 255}) // dark red for obstacle
	}

	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			tb.Fatal(err)
		}
	}(file)

	err = png.Encode(file, img)
	if err != nil {
		panic(err)
	}
}
