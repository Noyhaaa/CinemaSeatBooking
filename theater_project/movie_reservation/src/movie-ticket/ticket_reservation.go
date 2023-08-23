package movieticket

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/go-sql-driver/mysql"
	movieplace "movie.com/movie-place"
)

var db *sql.DB

type Movies struct {
	ID         int64
	Title      string
	Realisator string
	Leftticket int16
	Maxticket  int16
	Price      float32
	Room       int16
}

type Clients struct {
	ID              int64
	Titlename       int16
	Place_available string
}

type SqlValue struct {
	Key   string
	Value interface{}
}

func InitDB() {
	/*
		Initialise the database server
	*/
	// Capture connection properties.
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "theater",
	}
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	UpdateErr := SetTicket()
	if err != nil {
		log.Fatal(UpdateErr)
	}
}

func GeneratePlace() int16 {
	/*
		Generate a random ticket id.
		Use to give the number of remaining ticket just in case of the remaining ones
		is higher than the maximum ticket
		
		Returns
		---------------
		int16:
			new remaining ticket
	*/
	ticket := int16(rand.Intn(542))
	return ticket
}

func DataBaseRowsQuery(cmd string, args ...interface{}) (*sql.Rows, error) {
	/*
		Allow the user to query with no or several parameters 
		
		Parameters
		---------------
		cmd: string 
			Query to use 

		args: interface 
			all of the other parameters given to the DataBaseRowsQuery call.
		
		Returns
		---------------
		*sql.Rows:
			rows pointer get from the database result 
		
		error:
			Error returned by the database query
		
	*/
	if len(args) == 0 {
		rows, err := db.Query(cmd)
		if err != nil {
			return nil, fmt.Errorf("an occured with the command: %s", cmd)
		}
		return rows, nil
	} else {
		rows, err := db.Query(cmd, args...)
		if err != nil {
			return nil, fmt.Errorf("an occured with the command: %s with followign argument %v", cmd, args)
		}
		return rows, nil
	}
}

func GetAvailableTicket(movie string) (int16, error) {
	/*
		Get the number of ticket available
		
		Parameters
		---------------
		movie: string 
			Name of the movie
		
		Returns
		---------------
		int16:
			Maximum ticket for the movie
		
		error:
			Error returned by the database query
		
	*/
	rows, _ := DataBaseRowsQuery(
		"SELECT leftticket FROM movies WHERE title = ?", movie)
	defer rows.Close()
	var leftTicket int16
	for rows.Next() {
		if err := rows.Scan(&leftTicket); err != nil {
			return 0, fmt.Errorf("%v", err)
		}
		fmt.Printf("For %s, %d ticket left\n", movie, leftTicket)
	}
	return leftTicket, nil
}

func GetMaxTicket(movie string) (int16) {
	/*
		Get the maximum ticket for the movie
		
		Parameters
		---------------
		movie: string 
			Name of the movie
		
		Returns
		---------------
		int16:
			Maximum ticket for the movie
		
	*/
	rows, _ := DataBaseRowsQuery(
		"SELECT maxticket FROM movies WHERE title = ?", movie)
	defer rows.Close()
	var maxticket int16
	for rows.Next() {
		if err := rows.Scan(&maxticket); err != nil {
			fmt.Errorf("%v", err)
		}
		fmt.Printf("For %s, %d ticket max\n", movie, maxticket)
	}
	return maxticket
}

func SetTicket() error {
	/*
		Set ticket if the remainign ticket is higher than the limit.
		
		Returns
		---------------
		error:
			Error returned by the query	
	*/
	rows, _ := DataBaseRowsQuery("SELECT * FROM movies")
	defer rows.Close()
	for rows.Next() {
		var movies Movies
		if err := rows.Scan(&movies.ID, &movies.Title, &movies.Realisator, &movies.Leftticket, &movies.Maxticket, &movies.Price, &movies.Room); err != nil {
			return fmt.Errorf("%v", err)
		}
		for movies.Leftticket > movies.Maxticket {
			newTicket := GeneratePlace()
			_, err := db.Exec("UPDATE movies SET leftticket = ? WHERE id = ?", newTicket, movies.ID)
			if err != nil {
				return fmt.Errorf("update failed: %v", err)
			}
			movies.Leftticket = newTicket
		}
	}
	return nil
}

func SelectTicket(movie string, ticket int16) {
	/*
		Select a ticket/seat asked by the customer.
		Only the available seats can be buy.
		
		Parameters
		---------------
		movie: string 
			Name of the movie
		
		ticket: int16 
			ID of the movie the ticket/seat.
		
		
	*/
	query := fmt.Sprintf("SELECT place_available FROM `%s` WHERE id = ?", movie)
	var clients Clients
	// Empty set from the data base if the ticket does not exist
	empty_set := true
	err := db.QueryRow(query, ticket).Scan(&clients.Place_available)
	if err != nil {
        if err == sql.ErrNoRows {
            fmt.Println("L'ID de la place n'existe pas.")
			empty_set = false
        } else {
            fmt.Println("Erreur lors de la requête SQL:", err)
        }
    }
	if clients.Place_available != "False" && empty_set{
		query := fmt.Sprintf("UPDATE `%s` SET place_available = 'False' WHERE `%s` = %d", movie, movie, ticket)
		_, err := db.Exec(query)
		if err != nil {
			fmt.Printf("update failed: %v", err)
		}
		// -1 to the available ticket
		err = BuyTicket(movie)
		if err != nil {
			fmt.Printf("Impossible to buy the ticket because of: %v", err)
		}
		clients.Place_available = "False"
	} else {
		fmt.Printf("A pop up should be openned and a message should prevent the customer to select an another place.")
	}
	// UPDATE squares.png with red place in purpose
	CreatRoomSchema(movie)
}

func BuyTicket(movie string) error {
	/*
		Get name to the corresponding ID
		
		Parameters
		---------------
		movie: string 
			Name of the movie
		
		Returns
		---------------
		error
			error returned by the database query
	*/
	ticketAvailable, _ := GetAvailableTicket(movie)
	if ticketAvailable > 0 {
		_, err := db.Exec("UPDATE movies SET leftticket = ? WHERE title = ?", ticketAvailable-1, movie)
		if err != nil {
			return fmt.Errorf("update failed: %v", err)
		}
	} else if ticketAvailable < 1 || ticketAvailable > GetMaxTicket(movie){
		return fmt.Errorf("update failed, the place does not exist")
	}
	return nil
}

func MovieId() []int64 {
	/*
		Get available ID movies
		
		Returns
		---------------
		[]int64
			Correspond to the movies IDs
	*/
	rows, _ := DataBaseRowsQuery("SELECT id FROM movies;")
	defer rows.Close()
	var id []int64
	for rows.Next() {
		var movies Movies
		if err := rows.Scan(&movies.ID); err != nil {
			fmt.Printf("%v", err)
		}
		id = append(id, movies.ID)
	}
	fmt.Println(id)
	return id
}

func GetMovieName() map[string]int64 {
	/*
		Get name to the corresponding ID
		
		Returns
		---------------
		map[string]int64
			Map with movie name associated to an ID
	*/
	id := MovieId()
	var movieTitleId = make(map[string]int64)
	for _, elem := range id {
		rows, _ := DataBaseRowsQuery("SELECT title FROM movies WHERE id = ?;", elem)
		defer rows.Close()
		for rows.Next() {
			var movies Movies
			if err := rows.Scan(&movies.Title); err != nil {
				fmt.Printf("%v", err)
			}
			movieTitleId[movies.Title] = int64(elem)
		}
	}
	fmt.Println(movieTitleId)
	return movieTitleId
}

func CreatRoomSchema(movie string) {
	/*
		Get name to the corresponding ID
		
		Parameters
		---------------
		movie: string 
			Name of the movie 
	*/
	var availableSeats []string
	var clients Clients

	maxticket := GetMaxTicket(movie)
	query := fmt.Sprintf("SELECT `%s`, place_available FROM `%s`", movie, movie)
	rows, _ := DataBaseRowsQuery(query)
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&clients.Titlename, &clients.Place_available); err != nil {
			fmt.Printf("%v", err)
		}
		availableSeats = append(availableSeats, clients.Place_available)
	}
	movieplace.RoomSchema(movie, maxticket, availableSeats)
}

func DeleteMovie(movie string) {
	/*
		Delete movie from the database
		
		Parameters
		---------------
		movie: string 
			Name of the movie 
	*/
	_, err := db.Exec("DELETE FROM movies WHERE title = ?", movie)

	if err != nil {
		fmt.Printf("Delete failed: %v", err)
	}

	rows, _ := DataBaseRowsQuery("SELECT title FROM movies WHERE title = ?", movie)
	defer rows.Close()

	if !rows.Next() {
		fmt.Println("Movie has been deleted")
		return
	}
}

func CurrentlyInTheaters(table string) {
	/*
		Create a table in the database.
		The table is corresponding to the movies 
		currently in theater.
		
		Parameters
		---------------
		table: string 
			Name of the table 
	*/
	// Execute SQL requets

	DeleteTableQuery := fmt.Sprintf("DROP TABLE IF EXISTS %s;", table)

	_, err := db.Exec(DeleteTableQuery)

	if err != nil {
		fmt.Println("The request to delete the table failed: ", err)
	}

	createTableQuery := fmt.Sprintf(`
		CREATE TABLE %s (
		id         INT AUTO_INCREMENT NOT NULL,
		title      VARCHAR(128) NOT NULL,
		realisator VARCHAR(128) NOT NULL,
		leftticket SMALLINT,
		maxticket  SMALLINT,
		price      DECIMAL(5,2),
		room       SMALLINT, 
		PRIMARY KEY (id)
		);`, table)
	_, err = db.Exec(createTableQuery)

	if err != nil {
		fmt.Println("The request to create the movies currently in theater failed: ", err)
	}

	InsertInTable(table)
}

func CreatTableForMovies() {
	/*
		Create several tables corresponding to the movies
		currently in theater.
		Insert if the ticket is available into each movie table.
	*/
	movies_map := GetMovieName()
	for movie, _ := range movies_map {
		// Backsticks used because movie can contain space in the name.
		DeleteTableQuery := fmt.Sprintf("DROP TABLE IF EXISTS `%s`;", movie)

		_, err := db.Exec(DeleteTableQuery)

		if err != nil {
			fmt.Println("The request to delete the table failed: ", err)
		}

		// Drop the procedure
		DeleteProcQuery := fmt.Sprintf("DROP PROCEDURE IF EXISTS available_ticket;")

		_, err = db.Exec(DeleteProcQuery)

		if err != nil {
			fmt.Println("The request to delete the procedure failed: ", err)
		}

		// Create a table for each movie currently in theater
		createTableQuery := fmt.Sprintf("CREATE TABLE `%s` (id INT AUTO_INCREMENT, `%s` SMALLINT, place_available VARCHAR(128) NOT NULL, PRIMARY KEY (id))", movie, movie)
		_, err = db.Exec(createTableQuery)
		if err != nil {
			panic(err.Error())
		}

		// Prepare the instruction
		procedureQuery := fmt.Sprintf(`
			CREATE PROCEDURE available_ticket()
			BEGIN
				DECLARE i INT DEFAULT 1;
				WHILE i <= (SELECT maxticket FROM movies WHERE title = '%s') DO
					INSERT INTO %s (
						%s,
						place_available
					) VALUES (
						i, 'True'
					);
					SET i = i + 1;
				END WHILE;
			END;
		`, movie, "`"+movie+"`", "`"+movie+"`")

		_, err = db.Exec(procedureQuery)
		if err != nil {
			log.Fatal(err)
		}

		// Call the instruction
		_, err = db.Exec("CALL available_ticket()")
		if err != nil {
			log.Fatal(err)
		}
	}

}

func InsertInTable(table string) {
	/*
		Insert movies in the giventable
		
		Parameters
		---------------
		table: string 
			Name of the table 
	*/
	// Execute SQL request
	insertDataQuery := fmt.Sprintf(`
		INSERT INTO %s
		(title, realisator, leftticket, maxticket, price, room)
		VALUES
		('Elemental', 'Peter Sohn', 259, 259, 11.9, 1),
		('Asteroid City', 'Wes Anderson', 405, 405, 11.99, 2),
		('Spiderman Across the Spiderverse', 'Joaquim Dos Santos', 525, 525, 11.99, 3),
		('Indiana Jones', 'George Lucas', 179, 179, 11.99, 4);
	`, table)

	_, err := db.Exec(insertDataQuery)
	if err != nil {
		fmt.Println("Erreur lors de l'insertion des données:", err)
		return
	}
}
