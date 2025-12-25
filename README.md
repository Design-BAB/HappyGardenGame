# 🌸 Happy Garden

A game project built with [raylib-go](https://github.com/gen2brain/raylib-go) and Go.  
The goal: keep your garden happy by watering flowers before they wilt or mutate into evil fang-flowers!

## 🎮 Gameplay

- You control the **cow gardener** 🐄.
- Move with the arrow keys:
  - **↑ ↓ ← →** to walk around.
- Press **Space** to water flowers.
- Flowers will:
  - 🌼 Grow over time.
  - 🥀 Wilt if left unattended.
  - 🌑 Mutate into evil fang-flowers that move around the garden.
- **Game Over** happens if:
  - You collide with an evil flower.
  - Your garden stays unhappy for too long (10 seconds).
  - Or you reach the frame/time limit.


## ✨ Features

- **Clean code principles** applied:
  - Bounded constants (`MaxFlowers`, `MaxFrames`).
  - Defensive validation (`validateGameState`).
  - Encapsulated actor/plant types.
- **Scheduling without a scheduler**:  
  Uses `time.Time` fields to manage growth, wilting, and mutation events.
- **Y-sorting rendering**:  
  Creates a pseudo-depth effect so flowers can appear in front or behind the cow.
- **Loss conditions**:  
  Multiple ways to lose keep the game tense and engaging.

## 🛠️ Tech Stack

- **Language:** Go
- **Graphics:** [raylib-go](https://github.com/gen2brain/raylib-go)
