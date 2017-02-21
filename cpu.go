package main

import (
	"io/ioutil"
	"os"
)

//CPU is a struct that represents the CHIP-8 CPU component
type CPU struct {
	// Memory Map:
	// +---------------+= 0xFFF (4095) End of Chip-8 RAM
	// |               |
	// |               |
	// |               |
	// |               |
	// |               |
	// | 0x200 to 0xFFF|
	// |     Chip-8    |
	// | Program / Data|
	// |     Space     |
	// |               |
	// |               |
	// |               |
	// +- - - - - - - -+= 0x600 (1536) Start of ETI 660 Chip-8 programs
	// |               |
	// |               |
	// |               |
	// +---------------+= 0x200 (512) Start of most Chip-8 programs
	// | 0x000 to 0x1FF|
	// | Reserved for  |
	// |  interpreter  |
	// +---------------+= 0x000 (0) Start of Chip-8 RAM

	//4KB (4,096 bytes) of RAM
	memory []byte

	//I register is generally used to store memory addresses, so only the lowest (rightmost) 12 bits are usually used.
	I uint16

	//General purpose 8-bit registers, usually referred to as Vx
	registers [16]byte

	//PC should be 16-bit, and is used to store the currently executing address.
	PC uint16

	//The stack is an array of 16 16-bit values, used to store the address that the
	//interpreter shoud return to when finished with a subroutine. Chip-8 allows for up to 16 levels of nested subroutines.
	stack [16]uint16

	//SP can be 8-bit, it is used to point to the topmost level of the stack.
	SP byte

	//DT The delay timer is active whenever the delay timer register (DT) is non-zero.
	//This timer does nothing more than subtract 1 from the value of DT at a rate of 60Hz. When DT reaches 0, it deactivates.
	DT byte

	//ST The sound timer is active whenever the sound timer register (ST) is non-zero.
	//This timer also decrements at a rate of 60Hz, however, as long as ST's value is greater than zero,
	//the Chip-8 buzzer will sound. When ST reaches zero, the sound timer deactivates.
	ST byte

	//instruction is a 16 bit register to store the instruction code to execute
	instruction uint16
}

//NewCPU returns an uninitialized CPU struct
func NewCPU() *CPU {
	return &CPU{}
}

//Initialize loads the fontset and sets the PC to 0x200
func (c *CPU) Initialize() {
	for i := 0; i < 80; i++ {
		c.memory[i] = fontset[i]
	}
	c.PC = 0x200
}

//LoadGame loads the game ROM into memory
func (c *CPU) LoadGame(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	//Load rom in buffer
	var rom []byte
	rom, err = ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	//Write rom to memory from the buffer
	for i := 0; i < len(rom); i++ {
		c.memory[i+512] = rom[i]
	}
	return nil
}

//NextCycle executes next instruction
func (c *CPU) NextCycle() {
	c.instruction = uint16(c.memory[c.PC])<<8 | uint16(c.memory[c.PC+1])
	switch c.instruction & 0xF000 {
	case 0x0000:
		switch c.instruction & 0x0FFF {
		// case 0x00E0: //CLS
		// 	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		case 0x00EE: //RET
			c.PC = c.stack[c.SP-1]
			c.SP--
		default: //SYS ADDR
			c.SP++
			c.stack[c.SP] = c.instruction & 0x0FFF
		}
	case 0x1000: //JP addr
		c.PC = c.instruction & 0x0FFF
	case 0x2000: //CALL addr
		c.SP++
		c.stack[c.SP] = c.PC
		c.PC = c.instruction & 0x0FFF
	case 0x3000: //SE Vx, byte
		if int(c.instruction&0x00FF) == int(c.registers[int((c.instruction&0x0F00)>>8)]) {
			c.PC += 2
		}
	case 0x4000: //SNE Vx, byte
		if int(c.instruction&0x00FF) != int(c.registers[int((c.instruction&0x0F00)>>8)]) {
			c.PC += 2
		}
	case 0x5000: //SE Vx, Vy
		if int(c.registers[int((c.instruction&0x0F00)>>8)]) != int(c.registers[int((c.instruction&0x00F0)>>8)]) {
			c.PC += 2
		}
	case 0x6000: //LD Vx, byte
		c.registers[int((c.instruction&0x0F00)>>8)] = byte(c.instruction & 0x00FF)
	case 0x7000: //ADD Vx, byte
		c.registers[int((c.instruction&0x0F00)>>8)] += byte(c.instruction & 0x0FFF)
	case 0x8000:
		switch c.instruction & 0x000F {
		case 0x0000: //LD Vx, Vy
			c.registers[int((c.instruction&0x0F00)>>8)] = c.registers[int((c.instruction&0x00F0)>>8)]
		case 0x0001: //OR Vx, Vy
			c.registers[int((c.instruction&0x0F00)>>8)] |= c.registers[int((c.instruction&0x00F0)>>8)]
		case 0x0002: //AND Vx, Vy
			c.registers[int((c.instruction&0x0F00)>>8)] &= c.registers[int((c.instruction&0x00F0)>>8)]
		case 0x0003: //XOR Vx, Vy
			c.registers[int((c.instruction&0x0F00)>>8)] ^= c.registers[int((c.instruction&0x00F0)>>8)]
		}
	}
}
