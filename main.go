package main

import (
	"fmt"
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi"
)

const (
	LCD_PLATE_RS    = 15
	LCD_PLATE_RW    = 14
	LCD_PLATE_EN    = 13
	LCD_PLATE_D4    = 12
	LCD_PLATE_D5    = 11
	LCD_PLATE_D6    = 10
	LCD_PLATE_D7    = 9
	LCD_PLATE_RED   = 6
	LCD_PLATE_GREEN = 7
	LCD_PLATE_BLUE  = 8
)

func main() {
	embd.InitI2C()
	defer embd.CloseI2C()
	mcp := NewMCP230XX(0x01, 16)
	mcp.Config(LCD_PLATE_RW, embd.Out)
	mcp.Output(LCD_PLATE_RW, byte(embd.Low))
	fmt.Println("yay")
}
