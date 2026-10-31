package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const WINDOW_WIDTH int32 = 800
const WINDOW_HEIGHT int32 = 600
const TARGET_FPS int32 = 30

func main() {
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.SetTraceLogLevel(rl.LogError)
	rl.InitWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "GoBoard")
	rl.SetTargetFPS(TARGET_FPS)
	defer rl.CloseWindow()

	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	gb := NewGoboard()

	gb.ReadCommandFile("res/commands.txt")
	gb.WriteCommandFile("res/commands.txt.bak")

	gb.font = rl.LoadFontEx("res/LiberationMono-Regular.ttf", 20, nil, 0)
	defer rl.UnloadFont(gb.font)

	gb.key_sound = rl.LoadSound("res/key_sound.ogg")
	defer rl.UnloadSound(gb.key_sound)
	gb.select_sound = rl.LoadSound("res/select_sound.ogg")
	defer rl.UnloadSound(gb.select_sound)

	for !rl.WindowShouldClose() {
		if rl.IsWindowResized() {
			gb.WindowResized(rl.GetScreenWidth(), rl.GetScreenHeight())
		}

		if rl.IsWindowFocused() {
			if gb.Target_fps != TARGET_FPS {
				rl.SetTargetFPS(TARGET_FPS)
				gb.Target_fps = TARGET_FPS
			}
		} else {
			if gb.Target_fps != 1 {
				rl.SetTargetFPS(1)
				gb.Target_fps = 1
			}
		}

		if gb.HandleInput() < 0 {
			break
		}

		selection_height := (gb.selected * gb.font_size) + (10 * gb.font_size)
		if selection_height > gb.window_height {
			offset := -float32(selection_height - gb.window_height)
			gb.camera.Offset = rl.Vector2{X: 0.0, Y: offset}
		} else {
			gb.camera.Offset = rl.Vector2{X: 0.0, Y: 0.0}
		}

		rl.BeginDrawing()
		rl.BeginMode2D(gb.camera)
		if gb.copied == true {
			rl.PlaySound(gb.select_sound)
			rl.ClearBackground(rl.DarkGreen)
			gb.copied = false
		} else {
			gb.DrawList()
		}
		rl.EndMode2D()
		rl.EndDrawing()
	}

	gb.WriteCommandFile("res/commands.txt")
}
