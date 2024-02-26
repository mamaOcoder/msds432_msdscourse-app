package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// Connection details
var (
	Hostname = "postgres"
	Port     = 5432
	Username = "postgres"
	Password = "root"
	Database = "msds"
)

type MSDSCourse struct {
	CID     string `json:"courseI_D"`
	CNAME   string `json:"course_name"`
	CPREREQ string `json:"prerequisite"`
}

func openConnection() (*sql.DB, error) {
	// connection string
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		Hostname, Port, Username, Password, Database)

	// open database
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// The function returns the User ID of the username
// -1 if the user does not exist
func exists(cid string) int {
	cid = strings.ToLower(cid)

	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer db.Close()

	dID := -1
	statement := fmt.Sprintf(`SELECT "id" FROM "msdscoursecatalog" where "cid" = '%s'`, cid)
	rows, err := db.Query(statement)
	if err != nil {
		fmt.Println("Query error:", err)
		return -1
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			fmt.Println("Scan", err)
			return -1
		}
		dID = id
	}

	return dID
}

// AddUser adds a new user to the database
// Returns new User ID
// -1 if there was an error
func AddCourse(d MSDSCourse) error {
	d.CID = strings.ToLower(d.CID)

	db, err := openConnection()
	if err != nil {
		return fmt.Errorf("Unable to open DB connection")
	}
	defer db.Close()

	dID := exists(d.CID)
	if dID != -1 {
		return fmt.Errorf("Course already exists: %v", d.CID)
	}

	insertStatement := `insert into "msdscoursecatalog" ("cid", "cname", "cprereq")
	values ($1, $2, $3)`
	_, err = db.Exec(insertStatement, d.CID, d.CNAME, d.CPREREQ)
	if err != nil {
		return fmt.Errorf("db.Exec() error: %v", err)
	}

	return nil
}

// DeleteUser deletes an existing user
func DeleteCourse(cid string) error {
	db, err := openConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	// Does the ID exist?
	dID := exists(cid)
	if dID == -1 {
		return fmt.Errorf("Course with CID %s does not exist", cid)
	}

	// Delete
	deleteStatement := `delete from "msdscoursecatalog" where id=$1`
	_, err = db.Exec(deleteStatement, dID)
	if err != nil {
		return err
	}

	return nil
}

func SearchCourse(cid string) (MSDSCourse, error) {
	var data MSDSCourse

	db, err := openConnection()
	if err != nil {
		return data, err
	}
	defer db.Close()

	statement := fmt.Sprintf(`SELECT "id", "cid", "cname", "cprereq" FROM "msdscoursecatalog" where "cid" = '%s'`, cid)
	rows, err := db.Query(statement)

	for rows.Next() {
		var id int
		var cid string
		var cname string
		var cp string
		err = rows.Scan(&id, &cid, &cname, &cp)
		if err != nil {
			fmt.Println("Scan", err)
			return data, err
		}
		data = MSDSCourse{CID: cid, CNAME: cname, CPREREQ: cp}
	}
	defer rows.Close()

	return data, nil
}

// ListUsers lists all users in the database
func ListCourses() ([]MSDSCourse, error) {
	Data := []MSDSCourse{}
	db, err := openConnection()
	if err != nil {
		return Data, err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT * FROM "msdscoursecatalog"`)
	if err != nil {
		return Data, err
	}

	for rows.Next() {
		var id int
		var cid string
		var cname string
		var cp string
		err = rows.Scan(&id, &cid, &cname, &cp)
		temp := MSDSCourse{CID: cid, CNAME: cname, CPREREQ: cp}
		Data = append(Data, temp)
		if err != nil {
			return Data, err
		}
	}
	defer rows.Close()
	return Data, nil
}

// UpdateUser is for updating an existing user
func UpdateCourse(d MSDSCourse) error {
	db, err := openConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	dID := exists(d.CID)
	if dID == -1 {
		return errors.New("Course does not exist")
	}

	updateStatement := `update "msdscoursecatalog" set "cid"=$1, "cname"=$2, "cprereq"=$3 where "id"=$4`
	_, err = db.Exec(updateStatement, d.CID, d.CNAME, d.CPREREQ, dID)
	if err != nil {
		return err
	}

	return nil
}
