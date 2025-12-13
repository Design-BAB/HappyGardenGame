//Author: Design-BAB
//Date: 12/12/2025
//Description: It is my happy garden game project. The goal is to reach 268 lines of code
//Continue from the suggestions on pg164

package main

import (
	"time"

	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
	//"time"
)

type GameState struct {
	Width               int32
	Height              int32
	CenterX             int32
	CenterY             int32
	IsOver              bool
	Finalized           bool
	HappyGarden         bool
	FangFlowerCollision bool
	TimeElapsed         int
	StartTime           time.Time
	//Go nor raylib has a scheduler func, so we need to calculate time ourselves
	LastFlowerTime time.Time
}

// New thing I didnt know about Go before, all variables and bools are automaticly zero or false.
func newGameState() *GameState {
	startTimeNow := time.Now()
	return &GameState{Width: 800, Height: 600, CenterX: 400, CenterY: 300, HappyGarden: true, StartTime: startTimeNow}
}

type Actor struct {
	Texture rl.Texture2D
	//this is the collision box``
	rl.Rectangle // This gives Actor all the fields of rl.Rectangle (X, Y, Width, Height)
	Speed        float32
}

func newActor(texture rl.Texture2D, x, y float32) *Actor {
	return &Actor{Texture: texture, Rectangle: rl.Rectangle{X: x, Y: y, Width: float32(texture.Width), Height: float32(texture.Height)}, Speed: 5.0}
}

type Plant struct {
	Texture rl.Texture2D
	//this is the collision box``
	rl.Rectangle // This gives Actor all the fields of rl.Rectangle (X, Y, Width, Height)
	Status       string
}

func newPlant(texture rl.Texture2D, x, y float32) *Plant {
	return &Plant{Texture: texture, Rectangle: rl.Rectangle{X: x, Y: y, Width: float32(texture.Width), Height: float32(texture.Height)}, Status: "Happy"}
}

// the book seems to want me to seperate what and when it happens. I decied to keep it in one.
func addFlower(flowerList, wiltFlower []*Plant, flowerTexture rl.Texture2D, yourGame *GameState) []*Plant {
	if yourGame.IsOver == false {
		if time.Since(yourGame.LastFlowerTime) >= 4*time.Second {
			flowerNew := newPlant(flowerTexture, float32(rand.IntN(int(yourGame.Width-100))+50), float32(rand.IntN(int(yourGame.Height-250))+150))
			flowerList = append(flowerList, flowerNew)
			yourGame.LastFlowerTime = time.Now()
		}
	}
	return flowerList
}

func checkWiltTimes() {
}

func wiltFlower() {
}

func checkFlowerCollision() {
}

func resetCow() {
}

func update(cow *Actor, yourGame *GameState) {
	//the book lists some variables, but i'm gonna skip that here
	if yourGame.IsOver == false {
		if rl.IsKeyDown(rl.KeyRight) {
			cow.X = cow.X + cow.Speed
		}
		if rl.IsKeyDown(rl.KeyLeft) {
			cow.X = cow.X - cow.Speed
		}
		if rl.IsKeyDown(rl.KeyUp) {
			cow.Y = cow.Y - cow.Speed
		}
		if rl.IsKeyDown(rl.KeyDown) {
			cow.Y = cow.Y + cow.Speed
		}
	}

}

func draw(cow *Actor, background *rl.Texture2D, yourGame *GameState) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)
	rl.DrawTexture(*background, 0, 0, rl.White)
	//rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.LightGray)
	if yourGame.IsOver == false {
		rl.DrawTexture(cow.Texture, int32(cow.X), int32(cow.Y), rl.White)
	}
	rl.EndDrawing()
}

func main() {
	//create a new game
	yourGame := newGameState()
	//creating window
	rl.InitWindow(yourGame.Width, yourGame.Height, "raylib [core] example - basic window")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)
	//this is just for background
	background := rl.LoadTexture("images/garden.png")
	defer rl.UnloadTexture(background)
	//time for actors
	cowTexture := rl.LoadTexture("images/cow.png")
	defer rl.UnloadTexture(cowTexture)
	cow := newActor(cowTexture, 100, 500)
	flowerTexture := rl.LoadTexture("images/flower.png")
	defer rl.UnloadTexture(flowerTexture)
	//below is from pg160
	flowerList := []*Plant{}
	wiltedList := []*Plant{}
	//Note: skipping step 9
	//fangFlowerList := []*Actor{}
	//this is the actual game loop
	for !rl.WindowShouldClose() {
		update(cow, yourGame)
		draw(cow, &background, yourGame)
	}
}
