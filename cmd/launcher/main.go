package main

import (
	_ "embed"
	"flag"
	"log"
	"os"
	"snaky/src/render"
	"snaky/src/snaky"
)

// version
//go:embed version
var version string

var msg = `Snake Game - ` + version + `

Usage:
./snaky -x 30 -y 15 -s fast

Help:
-x		int 	- grid width, default:  26
-y		int 	- grid height, default: 14
-s		string 	- game speed, default:  medium
		(slow - medium - fast)
`

// usage
func usage() {
	print(msg)
	os.Exit(0)
}

// main
func main() {
	// command args
	w := flag.Int("x", 26, "grid width")
	h := flag.Int("y", 14, "grid height")
	s := flag.String("s", "medium", "game speed")

	flag.Usage = usage

	flag.Parse()

	// min gird width
	if *w < 26 {
		*w = 26
	}

	// min grid height
	if *h < 14 {
		*h = 14
	}

	// render
	render, err := render.NewRender()
	if err != nil {
		log.Println(err)
		return
	}
	defer render.Close()

	// game speed
	var speed snaky.Speed = snaky.Medium
	switch *s {
	case "slow":
		speed = snaky.Slow
	case "fast":
		speed = snaky.Fast
	}

	// game
	game, err := snaky.NewGame(*w, *h, speed, render)
	if err != nil {
		log.Println(err)
		return
	}

	// start
	log.Println(game.Run())
}
