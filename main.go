//Author: Design-BAB
//Date: 12/16/2025
//Description: It is my happy garden game project. The goal is to reach 268 lines of code
//Notes: Try out the sugesstions from end of 168 to 169

package main

import (
	"math/rand/v2"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// trying to implement some of clean code principles
const (
	MaxFlowers = 100    // fixed upper bound for flowers (Rule 2/3)
	MaxFrames  = 432000 // bounded main loop frames (Rule 2)
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
	TimeToMutate     time.Time
}

// New thing I didnt know about Go before, all variables and bools are automaticly zero or false.
func newGameState() *GameState {
	startTimeNow := time.Now()
	return &GameState{Width: 800, Height: 600, CenterX: 400, CenterY: 300, IsGardenHappy: true, StartTime: startTimeNow, LastFlowerTime: startTimeNow, LastWiltTime: startTimeNow, TimeToMutate: startTimeNow}
}

func validateGameState(cow *Actor, yourGame *GameState, textures map[string]rl.Texture2D) {
	if cow == nil {
		panic("validateGameState: cow is nil")
	}
	if yourGame == nil {
		panic("validateGameState: yourGame is nil")
	}
	if len(textures) == 0 {
		panic("validateGameState: textures is empty")
	}
	if cow.Speed <= 0 {
		panic("validateGameState: invalid cow speed")
	}
	if yourGame.Width <= 0 || yourGame.Height <= 0 {
		panic("validateGameState: invalid dimensions")
	}
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
	Vx               float32
	Vy               float32
	IsHappy          bool
	IsEvil           bool
	LastWateringTime time.Time
}

func newPlant(texture rl.Texture2D, x, y float32) *Plant {
	startTimeNow := time.Now()
	//making the rectangle a little bit wider to make watering more accessable
	return &Plant{Texture: texture, Rectangle: rl.Rectangle{X: x, Y: y, Width: float32(texture.Width) * 1.6, Height: float32(texture.Height)}, IsHappy: true, LastWateringTime: startTimeNow}
}

// the book seems to want me to seperate what and when it happens. I decied to keep it in one.
func growFlowers(flowerList []*Plant, flowerTexture rl.Texture2D, yourGame *GameState) []*Plant {
	if yourGame.IsOver == false && len(flowerList) < MaxFlowers {
		if time.Since(yourGame.LastFlowerTime) >= 4*time.Second {
			flowerNew := newPlant(flowerTexture, float32(rand.IntN(int(yourGame.Width-100))+50), float32(rand.IntN(int(yourGame.Height-250))+150))
			flowerList = append(flowerList, flowerNew)
			yourGame.LastFlowerTime = time.Now()
		}
	}
	return flowerList
}

// from pg.166
func wiltFlower(flowerList []*Plant, dry rl.Texture2D) []*Plant {
	//The loop will find one random happy flower and mark it as unhappy and dry. Then it will break out of the loop, stopping further iterations.
	for i := 0; i < len(flowerList); i++ {
		pick := rand.IntN(len(flowerList))
		if flowerList[pick].IsHappy == true && flowerList[pick].IsEvil == false {
			flowerList[pick].IsHappy = false
			flowerList[pick].Texture = dry
			break
		}
	}
	return flowerList
}

func mutate(flowerList []*Plant, evil rl.Texture2D) []*Plant {
	if len(flowerList) > 0 {
		pick := rand.IntN(len(flowerList))
		if flowerList[pick].IsEvil == false {
			flowerList[pick].IsEvil = true
			flowerList[pick].IsHappy = true
			flowerList[pick].Texture = evil
			flowerList[pick].Vx = velocity()
			flowerList[pick].Vy = velocity()
		}
	}
	return flowerList
}

// from pg.170
func velocity() float32 {
	//grabs a number between 0&1, represents direction
	randomDir := rand.IntN(2)
	randomVelocity := rand.IntN(2) + 2
	if randomDir == 0 {
		//turns it into negative
		randomVelocity *= -1
	}
	return float32(randomVelocity)
}

//from pg 172

func checkEvilFlowers(cow *Actor, flowerList []*Plant, yourGame *GameState) {
	for _, flower := range flowerList {
		if flower.IsEvil && rl.CheckCollisionRecs(cow.Rectangle, flower.Rectangle) {
			yourGame.IsOver = true
		}
	}
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
					if flowerToCheck.IsEvil == false {
						flowerToCheck.Texture = theFlowerTexture["normal"]
					}
				}
			}
		}
	}
}

func updateFangFlower(flowerList []*Plant, yourGame *GameState) []*Plant {
	for i := 0; i < len(flowerList); i++ {
		if flowerList[i].IsEvil {
			if flowerList[i].X <= 0 || flowerList[i].X >= float32(yourGame.Width)-flowerList[i].Width {
				flowerList[i].Vx *= -1 // Reverse horizontal direction\
				flowerList[i].X = flowerList[i].X + flowerList[i].Vx
			} else {
				flowerList[i].X = flowerList[i].X + flowerList[i].Vx
			}
			if flowerList[i].Y <= 0 || flowerList[i].Y >= float32(yourGame.Height-250)+150 {
				flowerList[i].Vy *= -1 // Reverse vertical direction
				flowerList[i].Y = flowerList[i].Y + flowerList[i].Vy
			} else {
				flowerList[i].Y = flowerList[i].Y + flowerList[i].Vy
			}
		}
	}
	return flowerList
}

func update(cow *Actor, flowerList []*Plant, theCowTexture, theFlowerTexture map[string]rl.Texture2D, yourGame *GameState) []*Plant {
	//the book lists some variables, but i'm gonna skip that here
	validateGameState(cow, yourGame, theCowTexture)
	if yourGame.IsOver == false {
		//This is where we are going to schedule the stuff
		//the book had this in reset_cow function, but i just did it here because why not
		//the book suggest 500, did 350 because it looks better
		if time.Since(cow.LastWateringTime) >= 350*time.Millisecond && cow.Texture == theCowTexture["watering"] {
			cow.Texture = theCowTexture["normal"]
		}
		if time.Since(yourGame.LastWiltTime) >= 3*time.Second {
			flowerList = wiltFlower(flowerList, theFlowerTexture["dry"])
			yourGame.LastWiltTime = time.Now()
		}
		if time.Since(yourGame.TimeToMutate) >= 20*time.Second {
			flowerList = mutate(flowerList, theFlowerTexture["evil"])
			yourGame.TimeToMutate = time.Now()
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
		flowerList = growFlowers(flowerList, theFlowerTexture["normal"], yourGame)
		flowerList = updateFangFlower(flowerList, yourGame)
		checkEvilFlowers(cow, flowerList, yourGame)
	}
	return flowerList
}

func draw(cow *Actor, background *rl.Texture2D, flowerList []*Plant, yourGame *GameState) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)
	rl.DrawTexture(*background, 0, 0, rl.White)

	if yourGame.IsOver == false {
		// Y-sorting for perspective
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
	} else {
		// Centering Logic
		msg1 := "GARDEN LOST!"
		w1 := rl.MeasureText(msg1, 40)
		rl.DrawText(msg1, yourGame.CenterX-(w1/2), yourGame.CenterY, 40, rl.Red)

		if yourGame.IsGardenHappy == false && time.Since(yourGame.TimeSpentUnhappy) >= 10*time.Second {
			msg2 := "You left your flowers to be unhappy too long"
			w2 := rl.MeasureText(msg2, 20)
			rl.DrawText(msg2, yourGame.CenterX-(w2/2), yourGame.CenterY+50, 20, rl.DarkGray)
		}

		msg3 := "Press ESC to exit"
		w3 := rl.MeasureText(msg3, 20)
		rl.DrawText(msg3, yourGame.CenterX-(w3/2), yourGame.CenterY+100, 20, rl.DarkGray)
	}
	rl.EndDrawing()
}

func main() {
	//create a new game
	yourGame := newGameState()
	//creating window
	rl.InitWindow(yourGame.Width, yourGame.Height, "Happy Garden")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)
	//this is just for background
	background := rl.LoadTexture("images/garden.png")
	defer rl.UnloadTexture(background)
	//time for actors
	theCowTexture := map[string]rl.Texture2D{
		"normal":   rl.LoadTexture("images/cow.png"),
		"watering": rl.LoadTexture("images/cow-water.png"),
	}
	cow := newActor(theCowTexture["normal"], 100, 500)

	theFlowerTexture := map[string]rl.Texture2D{
		"normal": rl.LoadTexture("images/flower.png"),
		"dry":    rl.LoadTexture("images/flower-wilt.png"),
		"evil":   rl.LoadTexture("images/fangflower.png"),
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
	frames := 0
	for !rl.WindowShouldClose() && frames < MaxFrames {
		getInput(cow, flowerList, theCowTexture, theFlowerTexture, yourGame)
		flowerList = update(cow, flowerList, theCowTexture, theFlowerTexture, yourGame)
		draw(cow, &background, flowerList, yourGame)
		frames++
	}
}
