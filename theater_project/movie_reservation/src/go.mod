module movie.com/movie_reservation

go 1.20

replace movie.com/movie-ticket => ./movie-ticket

require (
	github.com/go-sql-driver/mysql v1.7.1
	movie.com/movie-ticket v0.0.0-00010101000000-000000000000
)
