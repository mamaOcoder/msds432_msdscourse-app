package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// JSONFILE resides in the current directory
var CSVFILE = "./msdscourse_data.csv"

type MSDSCourseCatalog []MSDSCourse

func readCSVFile(filepath string) (MSDSCourseCatalog, error) {
	_, err := os.Stat(filepath)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var data = MSDSCourseCatalog{}

	// CSV file read all at once
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}

	for _, line := range lines {
		temp := MSDSCourse{
			CID:     line[0],
			CNAME:   line[1],
			CPREREQ: line[2],
		}

		data = append(data, temp)
	}

	return data, nil
}

func buildDatabase() error {

	// Print out the courses added
	fmt.Println("Courses added from CSV:")
	data, err := ListCourses()
	if err != nil {
		fmt.Println(err)
		return err
	}
	for _, v := range data {
		fmt.Println(v)
	}

	return nil
}

func main() {

	// Make sure that the postgresql container has time to start before trying to add data
	var db *sql.DB
	var err error

	maxRetries := 10
	retryDelay := 5 * time.Second

	// Attempt to connect to the database with retries
	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", Hostname, Port, Username, Password, Database))
		if err == nil {
			break
		}
		log.Printf("Error connecting to database: %v\n", err)
		log.Printf("Retrying in %v...\n", retryDelay)
		time.Sleep(retryDelay)
	}

	// Check if connection was successful
	if err != nil {
		log.Fatalf("Failed to connect to database after %d retries: %v\n", maxRetries, err)
	}
	db.Close()

	// Database connection successful, proceed with your application logic
	fmt.Println("Connected to database successfully!")

	////////////////////
	// Build database
	////////////////////

	// Read csv file and add each to the database
	msds_data, err := readCSVFile(CSVFILE)
	if err != nil {
		fmt.Println("Error reading data from CSV:", err)
		return
	}

	// Loop through the courses and add each to the database
	for _, course := range msds_data {
		AddCourse(course)
	}

	// Build server
	mux := http.NewServeMux()
	s := &http.Server{
		Addr:         PORT,
		Handler:      mux,
		IdleTimeout:  10 * time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	}

	mux.Handle("/list", http.HandlerFunc(listHandler))
	mux.Handle("/insert/", http.HandlerFunc(insertHandler))
	mux.Handle("/insert", http.HandlerFunc(insertHandler))
	mux.Handle("/search", http.HandlerFunc(searchHandler))
	mux.Handle("/search/", http.HandlerFunc(searchHandler))
	mux.Handle("/delete/", http.HandlerFunc(deleteHandler))
	mux.Handle("/status", http.HandlerFunc(statusHandler))
	mux.Handle("/", http.HandlerFunc(defaultHandler))

	fmt.Println("Ready to serve at", PORT)
	err = s.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		return
	}
}
