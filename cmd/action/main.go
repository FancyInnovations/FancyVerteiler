package main

import "github.com/sethvargo/go-githubactions"

func main() {
	myInput := githubactions.GetInput("my_input")
	if myInput == "" {
		githubactions.Fatalf("missing input 'my_input'")
	}

	githubactions.Infof("Hello world: %s", myInput)

}
