package main

func main() {
	c := NewCPU()
	c.Initialize()
	c.LoadGame("games/PONG")

	for {
		//Emulate cycle
		
	}
}
