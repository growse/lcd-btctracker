package main

import (
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi"
)

const (
	MCP23017_IODIRA = 0x00
	MCP23017_IODIRB = 0x01
	MCP23017_GPIOA  = 0x12
	MCP23017_GPIOB  = 0x13
	MCP23017_GPPUA  = 0x0C
	MCP23017_GPPUB  = 0x0D
	MCP23017_OLATA  = 0x14
	MCP23017_OLATB  = 0x15
	MCP23008_GPIOA  = 0x09
	MCP23008_GPPUA  = 0x06
	MCP23008_OLATA  = 0x0A
)

type MCP230XX struct {
	i2c         embd.I2CBus
	address     byte
	num_gpios   int
	direction   byte
	outputvalue byte
}

func NewMCP230XX(address byte, num_gpios int) *MCP230XX {
	i2cbus := embd.NewI2CBus(address)
	mcp := MCP230XX{i2c: i2cbus, address: address, num_gpios: num_gpios}
	if num_gpios <= 8 {
		mcp.i2c.WriteByte(MCP23017_IODIRA, 0xFF)
		mcp.direction, _ = mcp.i2c.ReadByte(MCP23017_IODIRA)
		mcp.i2c.WriteByte(MCP23008_GPPUA, 0x00)
	} else if num_gpios > 8 {
		mcp.i2c.WriteByte(MCP23017_IODIRA, 0xFF)
		mcp.i2c.WriteByte(MCP23017_IODIRB, 0xFF)
		mcp.direction, _ = mcp.i2c.ReadByte(MCP23017_IODIRA)
		iodirb, _ := mcp.i2c.ReadByte(MCP23017_IODIRB)
		mcp.direction = mcp.direction | iodirb<<8
		mcp.i2c.WriteByte(MCP23017_GPPUA, 0x00)
		mcp.i2c.WriteByte(MCP23017_GPPUB, 0x00)
	}

	return &mcp
}

func (mcp *MCP230XX) Config(pin byte, mode embd.Direction) byte {
	if mcp.num_gpios <= 8 {
		mcp.direction = mcp.readAndChangePin(MCP23017_IODIRA, pin, byte(mode), 0)
	} else if mcp.num_gpios > 8 && mcp.num_gpios <= 16 {
		if pin < 8 {
			mcp.direction = mcp.readAndChangePin(MCP23017_IODIRA, pin, byte(mode), 0)
		} else {
			mcp.direction = mcp.direction | mcp.readAndChangePin(MCP23017_IODIRB, pin-8, byte(mode), 0)<<8
		}
	}
	return mcp.direction
}

func (mcp *MCP230XX) Output(pin byte, value byte) byte {
	if (mcp.num_gpios) <= 8 {
		olata, _ := mcp.i2c.ReadByte(MCP23008_OLATA)
		mcp.outputvalue = mcp.readAndChangePin(MCP23008_GPIOA, pin, value, olata)
		return mcp.outputvalue
	}
	if mcp.num_gpios <= 16 {
		if pin < 8 {
			olata, _ := mcp.i2c.ReadByte(MCP23017_OLATA)
			mcp.outputvalue = mcp.readAndChangePin(MCP23017_GPIOA, pin, value, olata)
		} else {
			olatb, _ := mcp.i2c.ReadByte(MCP23017_OLATB)
			mcp.outputvalue = mcp.readAndChangePin(MCP23017_GPIOB, pin-8, value, olatb) << 8
		}
	}
	return mcp.outputvalue
}

func (mcp *MCP230XX) readAndChangePin(port byte, pin byte, value byte, currValue byte) byte {
	if currValue == 0 {
		currValue, _ = mcp.i2c.ReadByte(port)
	}
	newValue := changeBit(currValue, pin, value)
	mcp.i2c.WriteByte(port, newValue)
	return newValue
}

func changeBit(bitmap byte, bit byte, value byte) byte {
	if value == 0 {
		return bitmap & (1 << bit)
	}
	if value == 1 {
		return bitmap | (1 << bit)
	}
	return bitmap
}
