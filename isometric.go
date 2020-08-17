package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	rl "github.com/lachee/raylib-goplus/raylib"
)

var ( // MARK: var

	// maps
	blockmap     = make([]isoblock, blocknumber*2)
	blockvisible = make([]bool, blocknumber*2)
	treesmap     = make([]int, blocknumber*2)

	// images
	// terrain
	blockgrass1 = rl.NewRectangle(0, 0, 132, 98)

	// trees
	tree1 = rl.NewRectangle(510, 15, 97, 134)
	tree2 = rl.NewRectangle(653, 1, 87, 154)
	tree3 = rl.NewRectangle(807, 17, 64, 138)
	tree4 = rl.NewRectangle(926, 17, 84, 152)

	// isometric block grid
	gridon                            bool
	vertcount, horizcount, gridlayout int
	blockw                            = 132
	blockh                            = 66
	blocknumber                       int
	// core
	mouseblock     int
	framecount     int
	debugon        bool
	monh, monw     int
	monh32, monw32 int32
	mousepos       rl.Vector2
	imgs           rl.Texture2D
	camera         rl.Camera2D
)

type isoblock struct {
	xy, xy2, xy3, xy4, topleft rl.Vector2
}

/*

total blocks = 28 rows of 16 horizontal (vertical * 2)

*/

func raylib() { // MARK: raylib
	rl.InitWindow(monw, monh, "isometric")
	rl.SetExitKey(rl.KeyEnd)          // key to end the game and close window
	imgs = rl.LoadTexture("imgs.png") // load images
	rl.SetTargetFPS(30)
	// rl.HideCursor()
	// 	rl.ToggleFullscreen()
	for !rl.WindowShouldClose() {
		mousepos = rl.GetMousePosition()
		framecount++

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		rl.BeginMode2D(camera)

		count := 0
		drawx := 0
		drawy := 0
		blocktotal := 0

		// draw grid get xy
		for a := 0; a < gridlayout; a++ {

			if gridon {
				rl.DrawLine(drawx, drawy+(blockh/2), drawx+(blockw/2), drawy, rl.Fade(rl.Green, 0.1))                   // top left
				rl.DrawLine(drawx+(blockw/2), drawy, drawx+(blockw), drawy+(blockh/2), rl.Fade(rl.Green, 0.1))          // top right
				rl.DrawLine(drawx, drawy+(blockh/2), drawx+(blockw/2), drawy+blockh, rl.Fade(rl.Green, 0.1))            // bottom left
				rl.DrawLine(drawx+(blockw/2), drawy+(blockh), drawx+(blockw), drawy+(blockh/2), rl.Fade(rl.Green, 0.1)) // bottom right
			}

			blockxy := isoblock{}
			blockxy.xy = rl.NewVector2(float32(drawx), float32(drawy+(blockh/2)))           // left point
			blockxy.xy2 = rl.NewVector2(float32(drawx+(blockw/2)), float32(drawy))          // top point
			blockxy.xy3 = rl.NewVector2(float32(drawx+(blockw)), float32(drawy+(blockh/2))) // right point
			blockxy.xy4 = rl.NewVector2(float32(drawx+(blockw/2)), float32(drawy+(blockh))) // bottom point
			blockxy.topleft = rl.NewVector2(float32(drawx), float32(drawy))
			blockmap[blocktotal] = blockxy

			blockxy = isoblock{}
			blockxy.xy = rl.NewVector2(float32(drawx+(blockw/2)), float32(drawy+(blockh)))           // left point
			blockxy.xy2 = rl.NewVector2(float32(drawx+(blockw)), float32(drawy+(blockh/2)))          // top point
			blockxy.xy3 = rl.NewVector2(float32(drawx+(blockw+(blockw/2))), float32(drawy+(blockh))) // right point
			blockxy.xy4 = rl.NewVector2(float32(drawx+(blockw)), float32(drawy+(blockh+(blockh/2)))) // bottom
			blockxy.topleft = rl.NewVector2(float32(drawx+(blockw/2)), float32(drawy+(blockh/2)))
			blockmap[blocktotal+15] = blockxy

			blocktotal++
			count++
			drawx += 132

			if count == horizcount {
				blocktotal += 15
				count = 0
				drawx = 0
				drawy += 66
			}

		}

		// draw visible blocks
		for a := 0; a < blocknumber; a++ {
			if blockvisible[a] {
				checkblock := blockmap[a]
				rl.DrawTextureRec(imgs, blockgrass1, checkblock.topleft, rl.White)
				if treesmap[a] != 0 {
					v2 := rl.NewVector2(checkblock.xy2.X-40, checkblock.xy2.Y-100)
					switch treesmap[a] {
					case 1:
						rl.DrawTextureRec(imgs, tree1, v2, rl.Fade(rl.White, 0.7))
					case 2:
						rl.DrawTextureRec(imgs, tree2, v2, rl.Fade(rl.White, 0.7))
					case 3:
						rl.DrawTextureRec(imgs, tree3, v2, rl.Fade(rl.White, 0.7))
					case 4:
						rl.DrawTextureRec(imgs, tree4, v2, rl.Fade(rl.White, 0.7))
					}
				}
			}
			if mouseblock == a {
				checkblock := blockmap[a]
				rl.DrawTriangle(checkblock.xy, checkblock.xy4, checkblock.xy2, rl.Fade(rl.Red, 0.2))
				rl.DrawTriangle(checkblock.xy4, checkblock.xy3, checkblock.xy2, rl.Fade(rl.Red, 0.2))
			}
		}

		//	v2 := rl.NewVector2(10, 10)
		//	rl.DrawTextureRec(imgs, blockgrass1, v2, rl.White)

		rl.EndMode2D() // MARK: draw no camera
		update()
		rl.EndDrawing()
	}
	rl.CloseWindow()
}
func createsquare(sidelength, topblock int) { // MARK: createsquare

	//	area := sidelength*sidelength

	blockswitch := true
	blockswitch2 := true

	for b := 0; b < sidelength; b++ {
		nextblock := topblock
		for a := 0; a < sidelength; a++ {
			blockvisible[nextblock] = true
			if rolldice() == 6 {
				treesmap[nextblock] = rInt(1, 5)
			}
			if blockswitch {
				nextblock += 15
				blockswitch = false
			} else {
				nextblock += 14
				blockswitch = true
			}
		}
		if blockswitch2 {
			topblock += 16
			blockswitch = false
			blockswitch2 = false
		} else {
			topblock += 15
			blockswitch = true
			blockswitch2 = true
		}
	}

}
func getmouseblock() { // MARK: getmouseblock

	for a := 0; a < blocknumber; a++ {
		checkblock := blockmap[a]

		if mousepos.Y > checkblock.xy2.Y && mousepos.Y < checkblock.xy4.Y {
			if mousepos.X > checkblock.xy.X+30 && mousepos.X < checkblock.xy3.X-30 {
				mouseblock = a
			}
		}

	}

}
func update() { // MARK: update
	if debugon {
		debug()
	}
	input()
	getmouseblock()
}
func main() { // MARK: main
	rand.Seed(time.Now().UnixNano()) // random numbers
	rl.SetTraceLogLevel(rl.LogError) // hides info window
	rl.InitWindow(monw, monh, "isometric")
	setscreen()
	rl.CloseWindow()
	initialize()
	createlevel()
	raylib()
}
func input() { // MARK: input

	if rl.IsKeyPressed(rl.KeyKpAdd) {
		camera.Zoom += 0.1
	}
	if rl.IsKeyPressed(rl.KeyKpSubtract) {
		camera.Zoom -= 0.1
	}
	if rl.IsKeyPressed(rl.KeyKpMultiply) {
		if gridon {
			gridon = false
		} else {
			gridon = true
		}
	}
	if rl.IsKeyPressed(rl.KeyKpDecimal) {
		if debugon {
			debugon = false
		} else {
			debugon = true
		}
	}
}
func debug() { // MARK: debug
	rl.DrawRectangle(monw-300, 0, 500, monw, rl.Fade(rl.Blue, 0.4))
	rl.DrawFPS(monw-290, monh-100)

	mousexTEXT := fmt.Sprintf("%g", mousepos.X)
	mouseyTEXT := fmt.Sprintf("%g", mousepos.Y)
	vertcountTEXT := strconv.Itoa(vertcount)
	horizcountTEXT := strconv.Itoa(horizcount)
	blocknumberTEXT := strconv.Itoa(blocknumber)
	mouseblockTEXT := strconv.Itoa(mouseblock)

	checkblock := blockmap[17]

	block2xTEXT := fmt.Sprintf("%g", checkblock.xy.X)
	block2yTEXT := fmt.Sprintf("%g", checkblock.xy.Y)

	rl.DrawText(mousexTEXT, monw-290, 10, 10, rl.White)
	rl.DrawText("mouseX", monw-150, 10, 10, rl.White)
	rl.DrawText(mouseyTEXT, monw-290, 20, 10, rl.White)
	rl.DrawText("mouseY", monw-150, 20, 10, rl.White)
	rl.DrawText(vertcountTEXT, monw-290, 30, 10, rl.White)
	rl.DrawText("vertcount", monw-150, 30, 10, rl.White)
	rl.DrawText(horizcountTEXT, monw-290, 40, 10, rl.White)
	rl.DrawText("horizcount", monw-150, 40, 10, rl.White)
	rl.DrawText(blocknumberTEXT, monw-290, 50, 10, rl.White)
	rl.DrawText("blocknumber", monw-150, 50, 10, rl.White)
	rl.DrawText(block2xTEXT, monw-290, 60, 10, rl.White)
	rl.DrawText("block2x", monw-150, 60, 10, rl.White)
	rl.DrawText(block2yTEXT, monw-290, 70, 10, rl.White)
	rl.DrawText("block2y", monw-150, 70, 10, rl.White)
	rl.DrawText(mouseblockTEXT, monw-290, 80, 10, rl.White)
	rl.DrawText("mouseblock", monw-150, 80, 10, rl.White)

}
func createlevel() { // MARK: createlevel
	createsquare(7, 50)
}
func createmaps() { // MARK: createmaps

	treesmap = make([]int, blocknumber*2)
	blockmap = make([]isoblock, blocknumber*2)
	blockvisible = make([]bool, blocknumber*2)
}
func initialize() { // MARK: initialize

	vertcount = (monh / 66) + 1
	horizcount = (monw / 132) + 1
	gridlayout = horizcount * vertcount
	blocknumber = horizcount * (vertcount * 2)

	createmaps()

}
func setscreen() { // MARK: setscreen
	monh = rl.GetScreenHeight()
	monw = rl.GetScreenWidth()
	monh32 = int32(monh)
	monw32 = int32(monw)
	rl.SetWindowSize(monw, monh)
	camera.Zoom = 1.0
	camera.Target.X = 0
	camera.Target.Y = 0

} // random numbers
func rInt(min, max int) int {
	return rand.Intn(max-min) + min
}
func rInt32(min, max int) int32 {
	a := int32(rand.Intn(max-min) + min)
	return a
}
func rFloat32(min, max int) float32 {
	a := float32(rand.Intn(max-min) + min)
	return a
}
func flipcoin() bool {
	var b bool
	a := rInt(0, 10001)
	if a < 5000 {
		b = true
	}
	return b
}
func rolldice() int {
	a := rInt(1, 7)
	return a
}
