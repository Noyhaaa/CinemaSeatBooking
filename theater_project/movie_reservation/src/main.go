package main

import (
	"fmt"
	"movie.com/movie-ticket"
)

func main(){
	//Init the Database
	movieticket.InitDB()

	//Create a Table in the db with the movies currently in theater 
	movieticket.CurrentlyInTheaters("movies")

	//Create several table with the movies name 
	movieticket.CreatTableForMovies()

	ticket := movieticket.GeneratePlace()
	fmt.Printf("My generated number: %d\n", ticket)
	
	movieticket.SetTicket()

	movieticket.GetAvailableTicket("Asteroid City")

	movieticket.SelectTicket("Asteroid City", 128)	

	movieticket.SelectTicket("Elemental", 199)
	
	movieticket.SelectTicket("Asteroid City", 800)	

}