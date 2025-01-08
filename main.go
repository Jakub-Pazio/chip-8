package main

import (
	"crypto/rand"
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"os"
)

func checkPressed(key uint16) bool {
	switch key {
	case 1:
		return rl.IsKeyPressed(rl.KeyOne)
	case 2:
		return rl.IsKeyPressed(rl.KeyTwo)
	case 3:
		return rl.IsKeyPressed(rl.KeyThree)
	case 0xC:
		return rl.IsKeyPressed(rl.KeyFour)
	case 4:
		return rl.IsKeyPressed(rl.KeyQ)
	case 5:
		return rl.IsKeyPressed(rl.KeyW)
	case 6:
		return rl.IsKeyPressed(rl.KeyE)
	case 0xD:
		return rl.IsKeyPressed(rl.KeyR)
	case 7:
		return rl.IsKeyPressed(rl.KeyA)
	case 8:
		return rl.IsKeyPressed(rl.KeyS)
	case 9:
		return rl.IsKeyPressed(rl.KeyD)
	case 0xE:
		return rl.IsKeyPressed(rl.KeyF)
	case 0xA:
		return rl.IsKeyPressed(rl.KeyZ)
	case 0:
		return rl.IsKeyPressed(rl.KeyX)
	case 0xB:
		return rl.IsKeyPressed(rl.KeyC)
	case 0xF:
		return rl.IsKeyPressed(rl.KeyV)
	default:
		return false
	}
}

func GetKey() (uint8, bool) {
	if rl.IsKeyPressed(rl.KeyOne) {
		return 0x1, true
	}
	if rl.IsKeyPressed(rl.KeyTwo) {
		return 0x2, true
	}
	if rl.IsKeyPressed(rl.KeyThree) {
		return 0x3, true
	}
	if rl.IsKeyPressed(rl.KeyFour) {
		return 0xC, true
	}
	if rl.IsKeyPressed(rl.KeyQ) {
		return 0x4, true
	}
	if rl.IsKeyPressed(rl.KeyW) {
		return 0x5, true
	}
	if rl.IsKeyPressed(rl.KeyE) {
		return 0x6, true
	}
	if rl.IsKeyPressed(rl.KeyR) {
		return 0xD, true
	}
	if rl.IsKeyPressed(rl.KeyA) {
		return 0x7, true
	}
	if rl.IsKeyPressed(rl.KeyS) {
		return 0x8, true
	}
	if rl.IsKeyPressed(rl.KeyD) {
		return 0x9, true
	}
	if rl.IsKeyPressed(rl.KeyF) {
		return 0xE, true
	}
	if rl.IsKeyPressed(rl.KeyZ) {
		return 0xA, true
	}
	if rl.IsKeyPressed(rl.KeyX) {
		return 0x0, true
	}
	if rl.IsKeyPressed(rl.KeyC) {
		return 0xB, true
	}
	if rl.IsKeyPressed(rl.KeyV) {
		return 0xF, true
	}

	return 0, false // No key pressed
}

type Emul struct {
	Screen        [64 * 32]bool
	Memory        [4 * 1024]uint8
	PC            uint16
	IndexRegister uint16
	Stack         []uint16
	DelayTimer    uint8
	SoundTimer    uint8
	VarRegisters  [16]uint8
}

func New() Emul {
	e := Emul{Stack: make([]uint16, 0)}
	// Set up the fonts sprites
	fontSpriteArray := []uint8{
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80, // F
	}

	// save values into appropriate places in memory
	fontSpritePtr := 0x50
	for _, b := range fontSpriteArray {
		e.Memory[fontSpritePtr] = b
		fontSpritePtr++
	}
	return e
}

// Display draws n x 4 sized sprite on the screen at coordinates denoted by values in the
// registers specified by x and y parameters. Sprite comes from memory at address equal to IR
func (e *Emul) Display(x, y, n uint8) error {
	x = e.VarRegisters[x]
	y = e.VarRegisters[y]
	flipped := false
	for row := range n {
		fmt.Println("row", row)
		screenY := y + row
		if screenY > 31 {
			break
		}
		rowToDraw := e.Memory[e.IndexRegister+uint16(row)]
		fmt.Printf("%x\n", rowToDraw)
		for col := range uint8(8) {
			screenX := x + col
			if col > 63 {
				continue
			}
			fmt.Println(col, rowToDraw)
			pixel := (rowToDraw >> (7 - col)) & 0x1
			fmt.Println("pxl: ", pixel)
			if pixel != 0 {
				offset := uint16(screenY)*64 + uint16(screenX)
				fmt.Println("x:", screenX, "y:", screenY, "offset:", offset)
				if e.Screen[offset] {
					flipped = true
				}
				e.Screen[offset] = !e.Screen[offset] // Set the pixel on the screen
			}
		}
	}
	if flipped {
		e.VarRegisters[15] = 1
	}
	return nil
}

func main() {
	e := New()

	// fetch program into memory
	// 1-chip8-logo.ch8
	// 2-ibm-logo.ch8
	// 3-corax+.ch8
	// 4-flags.ch8
	// 5-quirks.ch8
	// 6-keypad.ch8
	f, err := os.Open("4-flags.ch8")
	if err != nil {
		panic(err)
	}
	buf := make([]byte, 2024, 2024)
	f.Read(e.Memory[512:])
	fmt.Println(string(buf))

	rl.InitWindow(640, 320, "Chip-8")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)
	e.PC = 512

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		// Drawing each pixel from the "Screen" onto the window created by raylib
		for i, pixel := range e.Screen {
			if pixel == true {
				x := i % 64
				y := i / 64
				rl.DrawRectangle(int32(x*10), int32(y*10), 10, 10, rl.Black)
			}
		}

		rl.EndDrawing()
		// Fetch
		instruction := (uint16(e.Memory[e.PC]) << 8) | uint16(e.Memory[e.PC+1])
		e.PC += 2
		if e.PC > 4095 {
			e.PC = 0
		}
		// Decrement both timers
		e.SoundTimer--
		e.DelayTimer--
		// Decode & Execute
		code := instruction & 0xF000
		// Checking for xxxx in the instruction
		switch instruction {
		case 0x00E0: // CLN
			fmt.Println("clearing screen")
			for i := range e.Screen {
				e.Screen[i] = false
			}
		case 0x00EE: // RET {
			retAddr := e.Stack[len(e.Stack)-1]
			e.Stack = e.Stack[:len(e.Stack)-1]
			e.PC = retAddr
		}
		fmt.Println(code)
		// Checking for x___ in the instruction
		switch code {
		case 0x1000: // JP addr
			val := instruction & 0x0FFF
			e.PC = val
		case 0x2000: // CALL addr
			addr := instruction & 0x0FFF
			e.Stack = append(e.Stack, e.PC)
			e.PC = addr
		case 0x3000: // SE Vx, byte
			regN := (instruction & 0x0F00) >> 8
			regVal := e.VarRegisters[regN]
			if uint16(regVal) == (instruction & 0x00FF) {
				e.PC += 2
			}
		case 0x4000: // SNE Vx, byte
			regN := (instruction & 0x0F00) >> 8
			regVal := e.VarRegisters[regN]
			if uint16(regVal) != (instruction & 0x00FF) {
				e.PC += 2
			}
		case 0x5000: // SE Vx, Vy
			reg1 := (instruction & 0x0F00) >> 8
			reg1Val := e.VarRegisters[reg1]
			reg2 := (instruction & 0x0F0) >> 4
			reg2Val := e.VarRegisters[reg2]
			if reg1Val == reg2Val {
				e.PC += 2
			}
		case 0x6000:
			regN := (instruction & 0x0F00) >> 8
			val := uint8(instruction & 0x00FF)
			e.VarRegisters[regN] = val
		case 0x7000:
			regN := (instruction & 0x0F00) >> 8
			val := uint8(instruction & 0x00FF)
			e.VarRegisters[regN] += val
		case 0xA000: // LD I, addr
			val := instruction & 0x0FFF
			e.IndexRegister = val
		case 0xB000: // JP V0, addr
			val := instruction & 0x0FFF
			e.PC = val + uint16(e.VarRegisters[0])
		case 0xC000: //RND Vx, addr
			regN := (instruction & 0x0F00) >> 8
			val := byte(instruction)
			b := make([]byte, 1)
			if _, err := rand.Read(b); err != nil {
				panic(err)
			}
			e.VarRegisters[regN] = val & b[0]
		case 0xD000:
			xReg := uint8((instruction & 0x0F00) >> 8)
			yReg := uint8((instruction & 0x00F0) >> 4)
			n := uint8(instruction & 0x000F)
			e.Display(xReg, yReg, n)
		}
		code = instruction & 0xF00F
		// Checking for x__x in the instruction
		switch code {
		case 0x8000: // LD Vx, Vy
			xReg := (instruction & 0x0F00) >> 8
			yReg := (instruction & 0x0F0) >> 4
			e.VarRegisters[xReg] = e.VarRegisters[yReg]
		case 0x8001: // OR Vx, Vy
			xReg := (instruction & 0x0F00) >> 8
			yReg := (instruction & 0x0F0) >> 4
			val := e.VarRegisters[xReg] | e.VarRegisters[yReg]
			e.VarRegisters[xReg] = val
		case 0x8002: // AND Vx, Vy
			xReg := (instruction & 0x0F00) >> 8
			yReg := (instruction & 0x0F0) >> 4
			val := e.VarRegisters[xReg] & e.VarRegisters[yReg]
			e.VarRegisters[xReg] = val
		case 0x8003: // XOR Vx, Vy
			xReg := (instruction & 0x0F00) >> 8
			yReg := (instruction & 0x0F0) >> 4
			val := e.VarRegisters[xReg] ^ e.VarRegisters[yReg]
			e.VarRegisters[xReg] = val
		case 0x8004: // ADD Vx, Vy
			xReg := (instruction & 0x0F00) >> 8
			yReg := (instruction & 0x00F0) >> 4
			val := uint16(e.VarRegisters[xReg]) + uint16(e.VarRegisters[yReg])
			e.VarRegisters[xReg] = uint8(val)
			if val > 0xFF {
				e.VarRegisters[0xF] = 1
			} else {
				e.VarRegisters[0xF] = 0
			}
		case 0x8005: // SUB Vx, Vy
			xReg := (instruction & 0x0F00) >> 8
			yReg := (instruction & 0x00F0) >> 4
			oFlag := false
			x := uint16(e.VarRegisters[xReg])
			y := uint16(e.VarRegisters[yReg])
			if y > x {
				x += 0x100
				oFlag = true
			}
			e.VarRegisters[xReg] = uint8(x - y)
			if oFlag {
				e.VarRegisters[0xF] = 0
			} else {
				e.VarRegisters[0xF] = 1
			}
		case 0x8006: // SHR Vx {, Vy} Shift Right
			xReg := (instruction & 0x0F00) >> 8
			xVal := e.VarRegisters[xReg]
			rightBit := xVal & 1
			e.VarRegisters[xReg] = xVal >> 1
			if rightBit == 0x1 {
				e.VarRegisters[0xF] = 1
			} else {
				e.VarRegisters[0xF] = 0
			}
		case 0x8007: // SUBN Vx, Vy
			xReg := (instruction & 0x0F00) >> 8
			yReg := (instruction & 0x00F0) >> 4
			oFlag := false
			x := uint16(e.VarRegisters[xReg])
			y := uint16(e.VarRegisters[yReg])
			if x > y {
				y += 0x100
				oFlag = true
			}
			e.VarRegisters[xReg] = uint8(y - x)
			if oFlag {
				e.VarRegisters[0xF] = 0
			} else {
				e.VarRegisters[0xF] = 1
			}
		case 0x800E: // SHL Vx {, Vy} Shift Left
			xReg := (instruction & 0x0F00) >> 8
			xVal := e.VarRegisters[xReg]
			leftBit := xVal & (1 << 7)
			e.VarRegisters[xReg] = xVal << 1
			if leftBit != 0 {
				e.VarRegisters[0xF] = 1
			} else {
				e.VarRegisters[0xF] = 0
			}
		case 0x9000: // SE Vx, Vy
			reg1 := (instruction & 0x0F00) >> 8
			reg1Val := e.VarRegisters[reg1]
			reg2 := (instruction & 0x0F0) >> 4
			reg2Val := e.VarRegisters[reg2]
			if reg1Val != reg2Val {
				e.PC += 2
			}
		}
		code = instruction & 0xF0FF
		// Checking for x_xx in the instruction
		switch code {
		case 0xE09E: // SKP Vx
			key := (instruction & 0x0F00) >> 8
			pressed := checkPressed(key)
			if pressed {
				e.PC += 2
			}
		case 0xE0A1: // SKNP Vx
			key := (instruction & 0x0F00) >> 8
			pressed := checkPressed(key)
			if !pressed {
				e.PC += 2
			}
		case 0xF007: // LD Vx, DT
			reg := (instruction & 0x0F00) >> 8
			e.VarRegisters[reg] = e.DelayTimer
		case 0xF00A: // LD Vx, K
			reg := (instruction & 0x0F00) >> 8
			key, ok := GetKey()
			if !ok {
				e.PC -= 2
			} else {
				e.VarRegisters[reg] = key
			}
		case 0xF015: // LD DT, Vx
			reg := (instruction & 0x0F00) >> 8
			e.DelayTimer = e.VarRegisters[reg]
		case 0xF018: // LD ST, Vx
			reg := (instruction & 0x0F00) >> 8
			e.SoundTimer = e.VarRegisters[reg]
		case 0xF01E: // ADD I, Vx
			reg := (instruction & 0x0F00) >> 8
			e.IndexRegister += uint16(e.VarRegisters[reg])
		case 0xF029: // LD F, Vx
			reg := (instruction & 0x0F00) >> 8
			offset := e.VarRegisters[reg]
			e.IndexRegister = 0x50 + 5*uint16(offset)
		case 0xF033: // LD B, Vx
			reg := (instruction & 0x0F00) >> 8
			val := e.VarRegisters[reg]
			e.Memory[e.IndexRegister] = val / 100         // Hundreds digit
			e.Memory[e.IndexRegister+1] = (val / 10) % 10 // Tens digit
			e.Memory[e.IndexRegister+2] = val % 10        // Ones digit
		case 0xF055: // LD [I], Vx
			rg := (instruction & 0x0F00) >> 8
			for i := range rg + 1 {
				e.Memory[e.IndexRegister+i] = e.VarRegisters[i]
			}
			e.IndexRegister += rg + 1
		case 0xF065: // LD Vx, [I]
			rg := (instruction & 0x0F00) >> 8
			for i := range rg + 1 {
				e.VarRegisters[i] = e.Memory[e.IndexRegister+i]
			}
			e.IndexRegister += rg + 1
		}
	}
}
