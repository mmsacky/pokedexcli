package main

func main() {

	config := &config{
		commandsRegistry: getCommands(),
		nextUrl:          "",
		previousUrl:      "",
	}

	startRepl(config)

}
