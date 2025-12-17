//Author: Design-BAB
//Date: 12/16/2025
//Description: It is my happy garden game project. The goal is to reach 268 lines of code
//Notes: Continue on suggestions on pg 166

package main

import (
	"math/rand/v2"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type GameState struct {
	Width               int32
	Height              int32
	CenterX             int32
	CenterY             int32
	IsOver              bool
	Finalized           bool
	IsGardenHappy       bool
	FangFlowerCollision bool
	TimeElapsed         int
	StartTime           time.Time
	//Go nor raylib has a scheduler func, so we need to calculate time ourselves
	LastFlowerTime   time.Time
	LastWiltTime     time.Time
	TimeSpentUnhappy time.Time
}

// New thing I didnt know about Go before, all variables and bools are automaticly zero or false.
func newGameState() *GameState {
	startTimeNow := time.Now()
	return &GameState{Width: 800, Height: 600, CenterX: 400, CenterY: 300, IsGardenHappy: true, StartTime: startTimeNow, LastFlowerTime: startTimeNow, LastWiltTime: startTimeNow}
}

type Actor struct {
	Texture          rl.Texture2D
	rl.Rectangle     // This gives Actor all the fields of rl.Rectangle (X, Y, Width, Height) and acts as the collision
	Speed            float32
	LastWateringTime time.Time
}

func newActor(texture rl.Texture2D, x, y float32) *Actor {
	startTimeNow := time.Now()
	return &Actor{Texture: texture, Rectangle: rl.Rectangle{X: x, Y: y, Width: float32(texture.Width), Height: float32(texture.Height)}, Speed: 5.0, LastWateringTime: startTimeNow}
}

type Plant struct {
	Texture rl.Texture2D
	//this is the collision box``
	rl.Rectangle     // This gives Actor all the fields of rl.Rectangle (X, Y, Width, Height)
	IsHappy          bool
	LastWateringTime time.Time
}

func newPlant(texture rl.Texture2D, x, y float32) *Plant {
	startTimeNow := time.Now()
	return &Plant{Texture: texture, Rectangle: rl.Rectangle{X: x, Y: y, Width: float32(texture.Width), Height: float32(texture.Height)}, IsHappy: true, LastWateringTime: startTimeNow}
}

// the book seems to want me to seperate what and when it happens. I decied to keep it in one.
func growFlowers(flowerList []*Plant, flowerTexture rl.Texture2D, yourGame *GameState) []*Plant {
	if yourGame.IsOver == false {
		if time.Since(yourGame.LastFlowerTime) >= 4*time.Second {
			flowerNew := newPlant(flowerTexture, float32(rand.IntN(int(yourGame.Width-100))+50), float32(rand.IntN(int(yourGame.Height-250))+150))
			flowerList = append(flowerList, flowerNew)
			yourGame.LastFlowerTime = time.Now()
		}
	}
	return flowerList
}

// from pg.166
func wiltFlower(flowerList []*Plant, dry rl.Texture2D, yourGame *GameState) []*Plant {
	//The loop will find one random happy flower and mark it as unhappy and dry. Then it will break out of the loop, stopping further iterations.
	for i := 0; i < len(flowerList); i++ {
		pick := rand.IntN(len(flowerList))
		if flowerList[pick].IsHappy {
			flowerList[pick].IsHappy = false
			flowerList[pick].Texture = dry
			break
		}
	}
	return flowerList
}

func getInput(cow *Actor, flowerList []*Plant, theCowTexture, theFlowerTexture map[string]rl.Texture2D, yourGame *GameState) {
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
		if rl.IsKeyDown(rl.KeySpace) {
			cow.Texture = theCowTexture["watering"]
			cow.LastWateringTime = time.Now()
			for _, flowerToCheck := range flowerList {
				if rl.CheckCollisionRecs(cow.Rectangle, flowerToCheck.Rectangle) {
					// Cow is touching the flower
					flowerToCheck.IsHappy = true
					flowerToCheck.Texture = theFlowerTexture["normal"]
				}
			}
		}
	}
}

func update(cow *Actor, flowerList []*Plant, theCowTexture, theFlowerTexture map[string]rl.Texture2D, yourGame *GameState) []*Plant {
	//the book lists some variables, but i'm gonna skip that here
	if yourGame.IsOver == false {
		//This is where we are going to schedule the stuff
		//the book had this in reset_cow function, but i just did it here because why not
		//the book suggest 500, did 350 because it looks better
		if time.Since(cow.LastWateringTime) >= 350*time.Millisecond && cow.Texture == theCowTexture["watering"] {
			cow.Texture = theCowTexture["normal"]
		}
		if time.Since(yourGame.LastWiltTime) >= 3*time.Second {
			flowerList = wiltFlower(flowerList, theFlowerTexture["dry"], yourGame)
			yourGame.LastWiltTime = time.Now()
		}

		//this is where we are going to consider the lossing conditions
		numberOfUnhappyPlants := 0
		for _, flowerToCheck := range flowerList {
			if flowerToCheck.IsHappy == false {
				yourGame.IsGardenHappy = false
				numberOfUnhappyPlants++
			}
		}
		if numberOfUnhappyPlants == 0 {
			yourGame.IsGardenHappy = true
			yourGame.TimeSpentUnhappy = time.Time{}
		}
		if yourGame.IsGardenHappy == false {
			if yourGame.TimeSpentUnhappy.IsZero() {
				startTimeNow := time.Now()
				yourGame.TimeSpentUnhappy = startTimeNow
			} else {
				if time.Since(yourGame.TimeSpentUnhappy) >= 10*time.Second {
					yourGame.IsOver = true
				}
			}
		}
	}
	flowerList = growFlowers(flowerList, theFlowerTexture["normal"], yourGame)
	return flowerList

}

func draw(cow *Actor, background *rl.Texture2D, flowerList []*Plant, yourGame *GameState) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)
	rl.DrawTexture(*background, 0, 0, rl.White)
	//rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.LightGray)
	if yourGame.IsOver == false {
		//The whole perspective thing is weird raw, so we need two loops to determine whether the flower should be infront or behind the cow
		for _, flowerToDisplay := range flowerList {
			if flowerToDisplay.Y < cow.Y+35 {
				rl.DrawTexture(flowerToDisplay.Texture, int32(flowerToDisplay.X), int32(flowerToDisplay.Y), rl.White)
			}
		}
		rl.DrawTexture(cow.Texture, int32(cow.X), int32(cow.Y), rl.White)
		for _, flowerToDisplay := range flowerList {
			if flowerToDisplay.Y >= cow.Y+35 {
				rl.DrawTexture(flowerToDisplay.Texture, int32(flowerToDisplay.X), int32(flowerToDisplay.Y), rl.White)
			}
		}
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
	//new texture code
	theCowTexture := map[string]rl.Texture2D{
		"normal":   rl.LoadTexture("images/cow.png"),
		"watering": rl.LoadTexture("images/cow-water.png"),
	}
	cow := newActor(theCowTexture["normal"], 100, 500)

	theFlowerTexture := map[string]rl.Texture2D{
		"normal": rl.LoadTexture("images/flower.png"),
		"dry":    rl.LoadTexture("images/flower-wilt.png"),
	}
	//Usage Example:
	//cow.Texture = cowTexture["normal"]
	//flower.Texture = theFlowerTexture["wilted"]

	// Cleanup:
	for _, texture := range theCowTexture {
		defer rl.UnloadTexture(texture)
	}
	for _, texture := range theFlowerTexture {
		defer rl.UnloadTexture(texture)
	}

	//below is from pg160
	flowerList := []*Plant{}
	//Note: skipping step 9
	//fangFlowerList := []*Actor{}
	//this is the actual game loop
	for !rl.WindowShouldClose() {
		getInput(cow, flowerList, theCowTexture, theFlowerTexture, yourGame)
		flowerList = update(cow, flowerList, theCowTexture, theFlowerTexture, yourGame)
		draw(cow, &background, flowerList, yourGame)
	}
}
