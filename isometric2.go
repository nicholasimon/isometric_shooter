package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	rl "github.com/lachee/raylib-goplus/raylib"
)

var ( // MARK: var

	// fx
	cloud1drift, cloud2drift, cloud3drift                 float32
	cloud1drifton, cloud2drifton, cloud3drifton           bool
	cloud1speed, cloud2speed, cloud3speed                 float32
	cloud1active, cloud2active, cloud3active, startclouds bool
	cloud1v2                                              = rl.NewVector2(-500, -500)
	cloud2v2                                              = rl.NewVector2(-500, -500)
	cloud3v2                                              = rl.NewVector2(-500, -500)
	scanlineson, cloudson, createclouds                   bool
	cloud1lr, cloud2lr, cloud3lr                          bool
	cloudscount                                           int
	// player
	playerx, playery float32
	playerxy         rl.Vector2
	playerdirection  int
	// maps
	blocktiles   = make([]int, blocknumber*2)
	blockmap     = make([]isoblock, blocknumber*2)
	blockvisible = make([]bool, blocknumber*2)
	treesmap     = make([]int, blocknumber*2)
	// images
	// clouds
	cloud1 = rl.NewRectangle(0, 876, 456, 148)
	cloud2 = rl.NewRectangle(459, 878, 426, 146)
	cloud3 = rl.NewRectangle(0, 694, 398, 170)
	// ships
	ship1  = rl.NewRectangle(144, 288, 140, 112)
	ship1l = rl.NewRectangle(3, 294, 134, 106)
	ship1r = rl.NewRectangle(303, 294, 134, 106)
	// terrain
	grassblock1  = rl.NewRectangle(0, 0, 132, 98)
	grassblock2  = rl.NewRectangle(136, 2, 132, 98)
	grassblock3  = rl.NewRectangle(136, 101, 132, 98)
	grassblock4  = rl.NewRectangle(1, 101, 132, 98)
	grassblock5  = rl.NewRectangle(272, 5, 132, 130)
	grassblock6  = rl.NewRectangle(6, 437, 132, 83)
	grassblock7  = rl.NewRectangle(5, 521, 132, 98)
	grassblock8  = rl.NewRectangle(295, 439, 132, 98)
	grassblock9  = rl.NewRectangle(440, 442, 132, 98)
	grassblock10 = rl.NewRectangle(148, 416, 132, 131)
	// trees
	tree1 = rl.NewRectangle(510, 15, 97, 134)
	tree2 = rl.NewRectangle(653, 1, 87, 154)
	tree3 = rl.NewRectangle(807, 17, 64, 138)
	tree4 = rl.NewRectangle(926, 17, 84, 152)
	// isometric block grid
	drawblock, nextblock              int
	screenblocknumber                 int
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
	cameraclouds   rl.Camera2D
)

type isoblock struct {
	xy, xy2, xy3, xy4, topleft rl.Vector2
}

/*

total blocks = 28 rows of 16 horizontal (vertical * 2)

*/

func raylib() { // MARK: raylib.
	rl.InitWindow(monw, monh, "isometric")
	rl.SetExitKey(rl.KeyEnd)          // key to end the game and close window
	imgs = rl.LoadTexture("imgs.png") // load images
	rl.SetTargetFPS(30)
	// rl.HideCursor()
	// 	rl.ToggleFullscreen()
	for !rl.WindowShouldClose() {
		mousepos = rl.GetMousePosition()

		checkblock := blockmap[nextblock]
		camera.Target.Y = checkblock.topleft.Y + 33

		framecount++

		if framecount%2 == 0 {
			if nextblock > horizcount*4 {
				nextblock -= 34
				camera.Target.Y -= 66
			}
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		rl.BeginMode2D(camera)

		// draw visible blocks
		drawblock = nextblock
		for a := 0; a < screenblocknumber; a++ {
			checkblock := blockmap[drawblock]
			checktiles := blocktiles[drawblock]
			rl.DrawTextureRec(imgs, grassblock1, checkblock.topleft, rl.White)

			switch checktiles {
			case 1:
				rl.DrawTextureRec(imgs, grassblock2, checkblock.topleft, rl.White)
			case 2:
				rl.DrawTextureRec(imgs, grassblock3, checkblock.topleft, rl.White)
			case 3:
				rl.DrawTextureRec(imgs, grassblock4, checkblock.topleft, rl.White)
			case 4:
				rl.DrawTextureRec(imgs, grassblock5, checkblock.topleft, rl.White)
			case 5:
				rl.DrawTextureRec(imgs, grassblock6, checkblock.topleft, rl.White)
			case 6:
				rl.DrawTextureRec(imgs, grassblock7, checkblock.topleft, rl.White)
			case 7:
				rl.DrawTextureRec(imgs, grassblock8, checkblock.topleft, rl.White)
			case 8:
				rl.DrawTextureRec(imgs, grassblock9, checkblock.topleft, rl.White)

			}

			if treesmap[drawblock] != 0 {
				v2 := rl.NewVector2(checkblock.xy2.X-40, checkblock.xy2.Y-100)
				switch treesmap[drawblock] {
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

			drawblock++
		}

		//	v2 := rl.NewVector2(10, 10)
		//	rl.DrawTextureRec(imgs, grassblock1, v2, rl.White)

		rl.EndMode2D() // MARK: draw no camera

		playershadowv2 := rl.NewVector2(playerx, playery+100)
		if playerdirection == 0 {
			rl.DrawTextureRec(imgs, ship1, playershadowv2, rl.Fade(rl.Black, 0.2))
			rl.DrawTextureRec(imgs, ship1, playerxy, rl.White)
		} else if playerdirection == 1 {
			rl.DrawTextureRec(imgs, ship1l, playershadowv2, rl.Fade(rl.Black, 0.2))
			rl.DrawTextureRec(imgs, ship1l, playerxy, rl.White)
		} else if playerdirection == 2 {
			rl.DrawTextureRec(imgs, ship1r, playershadowv2, rl.Fade(rl.Black, 0.2))
			rl.DrawTextureRec(imgs, ship1r, playerxy, rl.White)
		}

		// clouds
		if cloudson {
			rl.BeginMode2D(cameraclouds)
			cloud1shadowv2 := rl.NewVector2(cloud1v2.X, cloud1v2.Y+150)
			cloud2shadowv2 := rl.NewVector2(cloud2v2.X, cloud2v2.Y+150)
			cloud3shadowv2 := rl.NewVector2(cloud3v2.X, cloud3v2.Y+150)

			rl.DrawTextureRec(imgs, cloud1, cloud1shadowv2, rl.Fade(rl.Black, 0.1))
			rl.DrawTextureRec(imgs, cloud1, cloud1v2, rl.White)
			rl.DrawTextureRec(imgs, cloud2, cloud2shadowv2, rl.Fade(rl.Black, 0.1))
			rl.DrawTextureRec(imgs, cloud2, cloud2v2, rl.White)
			rl.DrawTextureRec(imgs, cloud3, cloud3shadowv2, rl.Fade(rl.Black, 0.1))
			rl.DrawTextureRec(imgs, cloud3, cloud3v2, rl.White)
			rl.EndMode2D()
		}
		// scan lines
		if scanlineson {
			for a := 0; a < monh; a++ {
				rl.DrawLine(0, a, monw, a, rl.Fade(rl.Black, 0.3))
				a++
			}
		}

		update()

		rl.EndDrawing()
	}
	rl.CloseWindow()
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
	playerxy = rl.NewVector2(playerx, playery)
	getmouseblock()
	clouds()

}
func clouds() { // MARK: clouds

	if cloudscount == 0 {
		if framecount%60 == 0 {
			cloud1lr = flipcoin()
			cloud2lr = flipcoin()
			cloud3lr = flipcoin()
			cloud1drifton = flipcoin()
			cloud2drifton = flipcoin()
			cloud3drifton = flipcoin()
			cloud1drift = rFloat32(1, 6)
			cloud2drift = rFloat32(1, 6)
			cloud3drift = rFloat32(1, 6)
			cloud1speed = rFloat32(8, 25)
			cloud2speed = rFloat32(8, 25)
			cloud3speed = rFloat32(8, 25)
			createclouds = true
			cloudscount = rInt(1, 4)

			zoomlevel := rFloat32(0, 11)
			zoomlevel = zoomlevel / 10
			zoomlevel++
			cameraclouds.Zoom = zoomlevel

		}
	}
	if cloudscount != 0 && createclouds {
		if cloud1lr {
			cloudx := rFloat32(-600, -450)
			cloudy := rFloat32(10, monh-200)
			cloud1v2 = rl.NewVector2(cloudx, cloudy)
			if cloudy > float32(monh/2) {
				cloud1drift = -(cloud1drift)
			}
		} else {
			cloudx := rFloat32(monw+10, monw+100)
			cloudy := rFloat32(10, monh-200)
			cloud1v2 = rl.NewVector2(cloudx, cloudy)
			if cloudy > float32(monh/2) {
				cloud1drift = -(cloud1drift)
			}
		}
		if cloud2lr {
			cloudx := rFloat32(-600, -450)
			cloudy := rFloat32(10, monh-200)
			cloud2v2 = rl.NewVector2(cloudx, cloudy)
			if cloudy > float32(monh/2) {
				cloud1drift = -(cloud1drift)
			}
		} else {
			cloudx := rFloat32(monw+10, monw+100)
			cloudy := rFloat32(10, monh-200)
			cloud2v2 = rl.NewVector2(cloudx, cloudy)
			if cloudy > float32(monh/2) {
				cloud1drift = -(cloud1drift)
			}
		}
		if cloud3lr {
			cloudx := rFloat32(-600, -450)
			cloudy := rFloat32(10, monh-200)
			cloud3v2 = rl.NewVector2(cloudx, cloudy)
			if cloudy > float32(monh/2) {
				cloud1drift = -(cloud1drift)
			}
		} else {
			cloudx := rFloat32(monw+10, monw+100)
			cloudy := rFloat32(10, monh-200)
			cloud3v2 = rl.NewVector2(cloudx, cloudy)
			if cloudy > float32(monh/2) {
				cloud1drift = -(cloud1drift)
			}
		}
		createclouds = false
	}

	if cloudscount != 0 && startclouds {
		if cloudscount == 1 {
			choose := rInt(1, 4)
			if choose == 1 {
				cloud1active = true
				cloud2active = false
				cloud3active = false
			} else if choose == 2 {
				cloud1active = false
				cloud2active = true
				cloud3active = false
			} else if choose == 3 {
				cloud1active = false
				cloud2active = false
				cloud3active = true
			}
		} else if cloudscount == 2 {
			choose := rInt(1, 4)
			choose2 := rInt(1, 4)
			if choose == 1 {
				cloud1active = true
				cloud2active = false
				cloud3active = false
			} else if choose == 2 {
				cloud1active = false
				cloud2active = true
				cloud3active = false
			} else if choose == 3 {
				cloud1active = false
				cloud2active = false
				cloud3active = true
			}
			if choose2 == 1 {
				if cloud1active == false {
					cloud1active = true
				} else {
					if flipcoin() {
						cloud2active = true
					} else {
						cloud3active = true
					}
				}
			} else if choose2 == 2 {
				if cloud2active == false {
					cloud2active = true
				} else {
					if flipcoin() {
						cloud1active = true
					} else {
						cloud3active = true
					}
				}
			} else if choose2 == 3 {
				if cloud3active == false {
					cloud3active = true
				} else {
					if flipcoin() {
						cloud1active = true
					} else {
						cloud2active = true
					}
				}
			}
		} else if cloudscount == 3 {
			cloud1active = true
			cloud2active = true
			cloud3active = true
		}

		startclouds = false
	}

	if cloud1active {
		if cloud1drifton {
			cloud1v2.Y += cloud1drift
		}
		if cloud1lr {
			cloud1v2.X += cloud1speed
			if cloud1v2.X > float32(monw) {
				cloud1active = false
				cloudscount--
			}
		} else {
			cloud1v2.X -= cloud1speed
			if cloud1v2.X < float32(-450) {
				cloud1active = false
				cloudscount--
			}
		}
	}

	if cloud2active {
		if cloud2drifton {
			cloud2v2.Y += cloud2drift
		}
		if cloud2lr {
			cloud2v2.X += cloud2speed
			if cloud2v2.X > float32(monw) {
				cloud2active = false
				cloudscount--
			}
		} else {
			cloud2v2.X -= cloud2speed
			if cloud2v2.X < float32(-450) {
				cloud2active = false
				cloudscount--
			}
		}
	}

	if cloud3active {
		if cloud3drifton {
			cloud3v2.Y += cloud3drift
		}
		if cloud3lr {
			cloud3v2.X += cloud3speed
			if cloud3v2.X > float32(monw) {
				cloud3active = false
				cloudscount--
			}
		} else {
			cloud3v2.X -= cloud3speed
			if cloud3v2.X < float32(-450) {
				cloud3active = false
				cloudscount--
			}
		}
	}

	if cloudscount < 0 {
		cloudscount = 0
	}

	if cloudscount == 0 {
		startclouds = true
	}

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

	if rl.IsKeyDown(rl.KeyLeft) {
		if playerx > 100 {
			playerx -= 60
		}
		playerdirection = 1
	} else if rl.IsKeyDown(rl.KeyRight) {
		if playerx < float32(monw-200) {
			playerx += 60
		}
		playerdirection = 2
	} else {
		playerdirection = 0
	}
	if rl.IsKeyDown(rl.KeyUp) {
		if playery > 100 {
			playery -= 60
		}
	}
	if rl.IsKeyDown(rl.KeyDown) {
		if playery < float32(monh-200) {
			playery += 60
		}
	}

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
	checkblock := blockmap[nextblock]
	checkblockxTEXT := fmt.Sprintf("%g", checkblock.topleft.X)
	checkblockyTEXT := fmt.Sprintf("%g", checkblock.topleft.Y)
	cameracloudszoomTEXT := fmt.Sprintf("%g", cameraclouds.Zoom)
	cameracloudsyTEXT := fmt.Sprintf("%g", cameraclouds.Target.Y)

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
	rl.DrawText(checkblockxTEXT, monw-290, 60, 10, rl.White)
	rl.DrawText("checkblockx", monw-150, 60, 10, rl.White)
	rl.DrawText(checkblockyTEXT, monw-290, 70, 10, rl.White)
	rl.DrawText("checkblocky", monw-150, 70, 10, rl.White)
	rl.DrawText(mouseblockTEXT, monw-290, 80, 10, rl.White)
	rl.DrawText("mouseblock", monw-150, 80, 10, rl.White)
	rl.DrawText(cameracloudszoomTEXT, monw-290, 90, 10, rl.White)
	rl.DrawText("cameracloudszoom", monw-150, 90, 10, rl.White)
	rl.DrawText(cameracloudsyTEXT, monw-290, 100, 10, rl.White)
	rl.DrawText("cameracloudsy", monw-150, 100, 10, rl.White)

}
func createlevel() { // MARK: createlevel

	count := 0
	drawx := 0
	drawy := 0
	blocktotal := 0

	for a := 0; a < blocknumber; a++ {

		blockxy := isoblock{}
		blockxy.xy = rl.NewVector2(float32(drawx), float32(drawy+(blockh/2)))           // left point
		blockxy.xy2 = rl.NewVector2(float32(drawx+(blockw/2)), float32(drawy))          // top point
		blockxy.xy3 = rl.NewVector2(float32(drawx+(blockw)), float32(drawy+(blockh/2))) // right point
		blockxy.xy4 = rl.NewVector2(float32(drawx+(blockw/2)), float32(drawy+(blockh))) // bottom point
		blockxy.topleft = rl.NewVector2(float32(drawx), float32(drawy))
		blockmap[blocktotal] = blockxy
		// add trees
		if rolldice() > 4 {
			treesmap[blocktotal] = rInt(1, 5)
		}

		blockxy = isoblock{}
		blockxy.xy = rl.NewVector2(float32(drawx+(blockw/2)), float32(drawy+(blockh)))           // left point
		blockxy.xy2 = rl.NewVector2(float32(drawx+(blockw)), float32(drawy+(blockh/2)))          // top point
		blockxy.xy3 = rl.NewVector2(float32(drawx+(blockw+(blockw/2))), float32(drawy+(blockh))) // right point
		blockxy.xy4 = rl.NewVector2(float32(drawx+(blockw)), float32(drawy+(blockh+(blockh/2)))) // bottom
		blockxy.topleft = rl.NewVector2(float32(drawx+(blockw/2)), float32(drawy+(blockh/2)))
		blockmap[blocktotal+17] = blockxy
		// add trees
		if rolldice() > 4 {
			treesmap[blocktotal+17] = rInt(1, 5)
		}

		blocktotal++
		count++
		drawx += 132

		if count == horizcount {
			blocktotal += 17
			count = 0
			drawx = 0
			drawy += 66
		}
	}

	leveltype()
}
func leveltype() { // MARK: leveltype()

	for a := 0; a < blocknumber; a++ {

		choose := rInt(1, 9)

		blocktiles[a] = choose

	}

}
func createmaps() { // MARK: createmaps

	blocktiles = make([]int, blocknumber*2)
	treesmap = make([]int, blocknumber*2)
	blockmap = make([]isoblock, blocknumber*2)
	blockvisible = make([]bool, blocknumber*2)
}
func initialize() { // MARK: initialize

	vertcount = (monh / 66)
	horizcount = (monw / 132) + 3
	gridlayout = horizcount * vertcount
	screenblocknumber = horizcount*(vertcount*2) + (horizcount * 4)
	blocknumber = horizcount * ((vertcount * 2) * 100)
	nextblock = blocknumber - screenblocknumber*2

	createmaps()
	playerx = float32(monw / 2)
	playery = float32(monh / 2)
	playerxy = rl.NewVector2(playerx, playery)
	scanlineson = true
	cloudson = true

}
func setscreen() { // MARK: setscreen
	monh = rl.GetScreenHeight()
	monw = rl.GetScreenWidth()
	monh32 = int32(monh)
	monw32 = int32(monw)
	rl.SetWindowSize(monw, monh)
	camera.Zoom = 1.0
	camera.Target.X = 66
	camera.Target.Y = 33

	cameraclouds.Zoom = 2.0

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
