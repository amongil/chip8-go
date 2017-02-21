package main

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

var fontset = [80]byte{
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

type key struct {
	x  int
	y  int
	ch rune
}

var K_1 = []key{{4, 4, '1'}}
var K_2 = []key{{7, 4, '2'}}
var K_3 = []key{{10, 4, '3'}}
var K_4 = []key{{13, 4, '4'}}
var K_5 = []key{{16, 4, '5'}}
var K_6 = []key{{19, 4, '6'}}
var K_7 = []key{{22, 4, '7'}}
var K_8 = []key{{25, 4, '8'}}
var K_9 = []key{{28, 4, '9'}}
var K_0 = []key{{31, 4, '0'}}

func print_tb(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func printf_tb(x, y int, fg, bg termbox.Attribute, format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	print_tb(x, y, fg, bg, s)
}

func drawChar() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyCtrlQ {
				break loop
			}

			if ev.Key == termbox.KeyArrowDown {
				printf_tb(33, 10, termbox.ColorBlue|termbox.AttrBold, termbox.ColorBlack, "Key pressed!")
			}
		}
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		termbox.Flush()
	}
}
