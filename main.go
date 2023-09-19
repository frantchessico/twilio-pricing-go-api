package main

import (
	"encoding/csv"
	// "encoding/json"
	"fmt"
	"log"
	"os"
    "github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

// Structure to represent a CSV line
type CSVRow struct {
	ISO         string `json:"ISO"`
	Country     string `json:"Country"`
	Description string `json:"Description"`
	PriceMsg    string `json:"Price / msg"`
}

// Function for CSV search
func searchCSV(filename string, searchColumn string, searchValue string) []CSVRow {
	var result []CSVRow

	// Open the CSV file for reading
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Read all lines from the CSV
	lines, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Map the CSV rows to the CSVRow structure
	for _, line := range lines {
		row := CSVRow{
			ISO:         line[0],
			Country:     line[1],
			Description: line[2],
			PriceMsg:    line[3],
		}
		// Check if the lookup column contains the desired value
		if searchColumn == "ISO" && row.ISO == searchValue {
			result = append(result, row)
		} else if searchColumn == "Country" && row.Country == searchValue {
			result = append(result, row)
		} else if searchColumn == "Description" && row.Description == searchValue {
			result = append(result, row)
		} else if searchColumn == "Price / msg" && row.PriceMsg == searchValue {
			result = append(result, row)
		}
	}

	return result
}

func getCountries(filename string) []string {
	var countries []string

	// Open the CSV file for reading
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Read all lines from the CSV
	lines, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Extract unique country names from CSV
	seen := make(map[string]bool)
	for _, line := range lines {
		country := line[1]
		if !seen[country] {
			countries = append(countries, country)
			seen[country] = true
		}
	}

	return countries
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Create a Fiber app
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	// Define routes
	app.Post("/search_csv", func(c *fiber.Ctx) error {
		// Parse the JSON request body
		var requestBody map[string]string
		if err := c.BodyParser(&requestBody); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Error parsing request body",
			})
		}
		searchColumn := requestBody["search_column"]
		searchValue := requestBody["search_value"]

		// Perform the search in the CSV file
		filename := "./datas.csv"
		searchResult := searchCSV(filename, searchColumn, searchValue)

		// Return results as JSON
		return c.JSON(searchResult)
	})

	app.Get("/countries", func(c *fiber.Ctx) error {
		// Get the list of unique countries
		filename := "./datas.csv"
		countryList := getCountries(filename)

		// Return country list as JSON
		return c.JSON(countryList)
	})

	// Get the PORT variable from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // Default port if PORT is not defined in .env
	}

	// Start the Fiber app
	fmt.Printf("Server running on port %s...\n", port)
	log.Fatal(app.Listen(":" + port))
}
