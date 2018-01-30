# GameOfLife
Golang implementation of Conway's Game of Life

## Getting program to run
What I did was:
* install golang to my system
* set the environment variable $GOPATH to $HOME/go
* chmod the file called goDependencies to executable
* run goDependencies, which will fetch all packages needed for my application
* then type the command :$ go run GameOfLife.go

## Using the program
There are several interactive commands:
* the mouse scroll wheel zooms the camera in and out
* the arrow keys move the camera
* the space bar pauses the program (this still isn't working properly for me)
* the n key advances the game one state (also not working properly for me)
* to reset the program you can close and restart or press r

## Things that break the program
This is the only thing as far as I know:
* If you choose a pattern besides random, you need to have at least 50 rows for glider gun and 25 for the other 2 or you get index out of bounds error
