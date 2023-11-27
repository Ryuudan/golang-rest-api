package database

import (
	"log"
	"os"

	_ "github.com/lib/pq" // Import the pq driver
	"github.com/ryuudan/golang-rest-api/ent/generated"
)

func PostgresClient() (*generated.Client, error) {
	// db := os.Getenv("POSTGRES_DB")
	// user := os.Getenv("POSTGRES_USER")
	// password := os.Getenv("POSTGRES_PASSWORD")

	// connectionString := fmt.Sprintf(
	// 	"host=localhost port=5432 user=%s dbname=%s password=%s sslmode=disable",
	// 	user, db, password,
	// )

	client, err := generated.Open("postgres", os.Getenv("POSTGRES_CONNECTION_STRING"))

	if err != nil {
		return nil, err
	}

	// Successfully connected to the database
	log.Println("✅ Sucessfully connected to the Postgres Database!")

	return client, nil
}
