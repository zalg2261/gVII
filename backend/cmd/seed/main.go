package main

import (
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/zalg2261/bioskop/backend/internal/db"
	"github.com/zalg2261/bioskop/backend/internal/models"
)

func main() {
	// Try to load .env from backend directory
	if err := godotenv.Load("../../.env"); err != nil {
		if err := godotenv.Load("../.env"); err != nil {
			log.Println(".env not loaded, relying on environment")
		}
	}

	// Connect to database
	db.Connect()

	// Clear existing movies and related data first
	log.Println("Clearing existing movies and showtimes...")
	db.DB.Exec("DELETE FROM showtimes;")
	db.DB.Exec("DELETE FROM movies;")
	
	log.Println("Seeding database with new movies...")

	// Insert movies
	movies := []models.Movie{
		{Title: "Your Name", Genre: "Romance, Drama", Duration: 106, Synopsis: "Two teenagers share a profound, magical connection upon discovering they are swapping bodies. Things manage to become even more complicated when the boy and girl decide to meet in person."},
		{Title: "Chainsaw Man: Reze arc", Genre: "Action, Horror", Duration: 120, Synopsis: "Denji faces off against the bomb devil Reze in this intense arc filled with action and emotional depth."},
		{Title: "The Tale of Princess Kaguya", Genre: "Fantasy, Drama", Duration: 137, Synopsis: "A tiny nymph found inside a bamboo stalk grows into a beautiful and desirable young woman, who orders her suitors to prove their love by completing near-impossible tasks."},
		{Title: "Perfect Blue", Genre: "Psychological Thriller", Duration: 81, Synopsis: "A pop singer gives up her career to become an actress, but she slowly goes insane when she starts being stalked by an obsessed fan and what seems to be a ghost of her past."},
		{Title: "Neon Genesis Evangelion 3.0+1.0", Genre: "Sci-Fi, Mecha", Duration: 155, Synopsis: "The final installment of the Rebuild of Evangelion series, where Shinji and his friends must confront the truth about the Human Instrumentality Project."},
		{Title: "Violet Evergarden: The Movie", Genre: "Drama, Slice of Life", Duration: 140, Synopsis: "Several years after the war, Violet Evergarden continues to work as an Auto Memory Doll, writing letters for others while searching for the meaning of love."},
		{Title: "Howl's Moving Castle", Genre: "Fantasy, Adventure", Duration: 119, Synopsis: "When an unconfident young woman is cursed with an old body by a spiteful witch, her only chance of breaking the spell lies with a self-indulgent yet insecure young wizard and his companions in his legged, walking home."},
		{Title: "Detective Conan: One-Eyed Flashback", Genre: "Mystery, Crime", Duration: 110, Synopsis: "Conan Edogawa investigates a mysterious case involving a one-eyed suspect and a series of flashbacks that reveal crucial clues."},
		{Title: "The Tatami Galaxy", Genre: "Comedy, Psychological", Duration: 90, Synopsis: "An unnamed third-year university student embarks on a journey through parallel universes, exploring different club activities and romantic pursuits in search of the perfect college life."},
	}

	for _, movie := range movies {
		if err := db.DB.Create(&movie).Error; err != nil {
			log.Printf("Error creating movie %s: %v", movie.Title, err)
		} else {
			log.Printf("Created movie: %s", movie.Title)
		}
	}

	// Insert branches
	branches := []models.Branch{
		{Name: "XXI Bandung PVJ", City: "Bandung", Address: "Jl. Sukajadi No. 123, Bandung"},
		{Name: "XXI Jakarta Kota Kasablanka", City: "Jakarta", Address: "Jl. Casablanca Raya No. 88, Jakarta Selatan"},
		{Name: "CGV Grand Indonesia", City: "Jakarta", Address: "Jl. MH Thamrin No. 1, Jakarta Pusat"},
		{Name: "XXI Surabaya Tunjungan Plaza", City: "Surabaya", Address: "Jl. Basuki Rahmat No. 8-12, Surabaya"},
		{Name: "Cinema 21 Yogyakarta Malioboro", City: "Yogyakarta", Address: "Jl. Malioboro No. 52-58, Yogyakarta"},
	}

	for _, branch := range branches {
		if err := db.DB.Create(&branch).Error; err != nil {
			log.Printf("Error creating branch %s: %v", branch.Name, err)
		} else {
			log.Printf("Created branch: %s", branch.Name)
		}
	}

	// Get created movies and branches for showtimes
	var createdMovies []models.Movie
	var createdBranches []models.Branch
	db.DB.Find(&createdMovies)
	db.DB.Find(&createdBranches)

	if len(createdMovies) > 0 && len(createdBranches) > 0 {
		// Insert showtimes for all movies
		showtimes := []models.Showtime{}
		
		// Create showtimes for each movie (at least 1 showtime per movie)
		for i, movie := range createdMovies {
			branchIdx := i % len(createdBranches)
			studioNum := (i % 3) + 1
			hoursOffset := (i * 2) + 1
			
			showtimes = append(showtimes, models.Showtime{
				MovieID:    movie.ID,
				BranchID:   createdBranches[branchIdx].ID,
				Studio:    fmt.Sprintf("Studio %d", studioNum),
				ShowTime:  time.Now().Add(time.Duration(hoursOffset*24) * time.Hour).Add(14 * time.Hour),
				SeatsTotal: 50,
				SeatsLeft:  50,
				Price:     50000,
				Status:    "ACTIVE",
			})
		}

		for _, showtime := range showtimes {
			if err := db.DB.Create(&showtime).Error; err != nil {
				log.Printf("Error creating showtime: %v", err)
			} else {
				log.Printf("Created showtime for movie ID %d", showtime.MovieID)
			}
		}
	}

	log.Println("Database seeded successfully!")
}

