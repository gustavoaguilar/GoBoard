package main

import (
	"bufio"
	"log"
	"os"
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type GoBoard struct {
	commands      []string
	selected      int32
	copied        bool
	font          rl.Font
	font_size     int32
	window_width  int32
	window_height int32
	camera        rl.Camera2D
	Target_fps    int32
	key_sound     rl.Sound
	select_sound  rl.Sound
}

func NewGoboard() GoBoard {
	gb := GoBoard{}
	gb.selected = 1
	gb.font_size = 20
	gb.camera = rl.Camera2D{Zoom: 1}
	gb.window_width = WINDOW_WIDTH
	gb.window_height = WINDOW_HEIGHT
	gb.Target_fps = TARGET_FPS

	return gb
}

func (gb *GoBoard) ReadCommandFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Error: %v\n", err.Error())
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		gb.commands = append(gb.commands, line)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error: %v\n", err.Error())
	}
}

func (gb *GoBoard) WriteCommandFile(path string) {
	file, err := os.Create(path)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	defer file.Close()

	for _, value := range gb.commands {
		line := value + "\n"
		file.Write([]byte(line))
	}
}

func (gb *GoBoard) WindowResized(width int, height int) {
	gb.window_width = int32(width)
	gb.window_height = int32(height)
}

func (gb GoBoard) PlayKeySound() {
	new_pitch := 1.0 + float32(float32(rl.GetRandomValue(0, 10))/100)
	rl.SetSoundPitch(gb.key_sound, float32(new_pitch))
	rl.PlaySound(gb.key_sound)
}

func (gb *GoBoard) HandleInput() int32 {
	if rl.IsKeyPressed(rl.KeyQ) {
		return -1
	}

	scroll_movement := rl.GetMouseWheelMoveV()
	if rl.IsKeyPressed(rl.KeyDown) || rl.IsKeyPressedRepeat(rl.KeyDown) ||
		rl.IsKeyPressed(rl.KeyS) || rl.IsKeyPressedRepeat(rl.KeyS) ||
		rl.IsKeyPressed(rl.KeyJ) || rl.IsKeyPressedRepeat(rl.KeyJ) ||
		scroll_movement.Y < 0 {
		gb.PlayKeySound()
		if gb.selected < int32(len(gb.commands)-1) {
			gb.selected++
		}
	}
	if rl.IsKeyPressed(rl.KeyUp) || rl.IsKeyPressedRepeat(rl.KeyUp) ||
		rl.IsKeyPressed(rl.KeyW) || rl.IsKeyPressedRepeat(rl.KeyW) ||
		rl.IsKeyPressed(rl.KeyK) || rl.IsKeyPressedRepeat(rl.KeyK) ||
		scroll_movement.Y > 0 {
		gb.PlayKeySound()
		if gb.selected >= 1 {
			gb.selected--
		}
	}

	if rl.IsKeyPressed(rl.KeySpace) || rl.IsKeyPressed(rl.KeyEnter) {
		gb.PlayKeySound()
		rl.SetClipboardText(gb.commands[gb.selected])
		gb.copied = true
	}

	if rl.IsKeyPressed(rl.KeyI) {
		gb.PlayKeySound()
		if gb.selected < 0 {
			gb.selected = 0
		}
		gb.commands = slices.Insert(gb.commands, int(gb.selected), rl.GetClipboardText())
	}
	if rl.IsKeyPressed(rl.KeyD) {
		gb.PlayKeySound()
		if len(gb.commands) > 0 {
			gb.commands = slices.Delete(gb.commands, int(gb.selected), int(gb.selected)+1)
			if gb.selected == int32(len(gb.commands)) {
				gb.selected--
			}
		}
	}

	return 0
}

func (gb *GoBoard) DrawList() {
	color := rl.LightGray
	mouse_pos := rl.GetMousePosition()

	rl.ClearBackground(rl.Black)

	for i, x := range gb.commands {
		//text_size := rl.MeasureText(x, gb.font_size)
		text_size := rl.MeasureTextEx(gb.font, x, float32(gb.font_size), 0)
		x_offset := gb.window_width/2 - int32(text_size.X)/2

		// Rows background
		if bool(i&0x01 == 1) {
			rl.DrawRectangle(0, int32(i)*gb.font_size, gb.window_width, gb.font_size, rl.Color{R: 0, G: 20, B: 0, A: 255})
		}

		if gb.selected == int32(i) {
			color = rl.Black
			//rl.DrawRectangle(0, int32(i)*gb.font_size, gb.window_width, gb.font_size, rl.Green)
			rl.DrawRectangle(x_offset-3, int32(i)*gb.font_size, int32(text_size.X)+6, gb.font_size, rl.Green)
		} else if rl.CheckCollisionPointRec(mouse_pos,
			rl.NewRectangle(float32(x_offset),
				float32(i*int(gb.font_size)+int(gb.camera.Offset.Y)),
				float32(rl.MeasureText(x, gb.font_size)),
				float32(gb.font_size))) {
			color = rl.Green
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				rl.SetClipboardText(x)
				gb.PlayKeySound()
				gb.selected = int32(i)
				gb.copied = true
			}
		} else {
			color = rl.DarkGreen
		}

		rl.DrawTextEx(
			gb.font,
			x,
			rl.Vector2{X: float32(x_offset), Y: float32(int32(i) * gb.font_size)},
			float32(gb.font_size),
			0,
			color)
	}

}
