package movieplace

import (
	"fmt"

	"github.com/fogleman/gg"
)

const (
	Width  = 720
	Height = 600
)

func RoomSchema(movie string, maxplace int16, availablePlaces []string){
	/*
		Allow the user to query with no or several parameters 
		
		Parameters
		---------------
		movie: string 
			Name of the movie

		maxplace: int16 
			number of the max ticket 

		availablePlaces: []string
			values (False or True) showing if which seat is available or not 
			to create the schema with the seats available or not available for
			the given movie.
	*/

	dc := gg.NewContext(Width, Height)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// Valeur de la variable pour déterminer la couleur
	heightCoef := 0.05

	// Place 
	var place int
	var availablePlace string

	// Loop to print the maximum place for a movie
	for place = 0; place < int(maxplace); place++ {
		var x, y, size float64
		availablePlace = availablePlaces[place]
		temp_place := place % 26
		// Calcul de la position et de la taille du carré
		if temp_place < 5 || temp_place > 20 {
			if temp_place == 0{
				heightCoef += 0.05
			}
			if temp_place > 20 {
				x = float64(temp_place*25 + 40 + 5) // index + two-square offset + offset between 2 square
			} else {
				x = float64(temp_place*25 + 5)
			}				
		} else {
			if temp_place == 5 || temp_place == 20 {
				x = float64(temp_place*25 + 25)
			} else {
				x = float64(temp_place*25 + 20+ 5) // index + one-square offset + offset between 2 square
			}
		}

		y = float64(Height*heightCoef - 50)
		size = float64(20)
		// Select the color that you want
		switch {
		case availablePlace == "False":
			dc.SetRGB(1, 0, 0) // Red
		case availablePlace == "True":
			dc.SetRGB(0, 0, 1) // Blue
		default:
			dc.SetRGB(0, 0, 0) // Black
		}

		dc.DrawRectangle(x, y, size, size)
		dc.Stroke()

		// Write into the square
		text := fmt.Sprintf("%d", place +1)
		dc.SetRGB(0, 0, 0) // set up  text color to white
		dc.DrawStringAnchored(text, x+size/2, y+size/2, 0.5, 0.5)
	}

	// Save the remaining place schema into movie_name.png
	dc.SavePNG(movie + ".png")
}