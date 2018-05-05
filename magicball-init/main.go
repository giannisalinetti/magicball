package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
)

const (
	driver = "mysql"
)

var (
	createTableStatement = "CREATE TABLE answers(id INT AUTO_INCREMENT, sentence VARCHAR(255), PRIMARY KEY(id));"
	insertDataStatements = []string{
		"INSERT INTO answers (sentence) VALUES ('As I see it, yes');",
		"INSERT INTO answers (sentence) VALUES ('It is certain');",
		"INSERT INTO answers (sentence) VALUES ('It is decidedly so');",
		"INSERT INTO answers (sentence) VALUES ('Most likely');",
		"INSERT INTO answers (sentence) VALUES ('Yes');",
		"INSERT INTO answers (sentence) VALUES ('Ask again later');",
		"INSERT INTO answers (sentence) VALUES ('Better not tell you now');",
		"INSERT INTO answers (sentence) VALUES ('Cannot predict now');",
		"INSERT INTO answers (sentence) VALUES ('Reply hazy, try again');",
		"INSERT INTO answers (sentence) VALUES ('Concentrate and ask again');",
		"INSERT INTO answers (sentence) VALUES ('Do not count on it');",
		"INSERT INTO answers (sentence) VALUES ('My reply is no');",
		"INSERT INTO answers (sentence) VALUES ('My sources say no');",
		"INSERT INTO answers (sentence) VALUES ('Outlook not so good');",
		"INSERT INTO answers (sentence) VALUES ('Very doubtful');",
	}
)

type dbConnection struct {
	user     string
	password string
	protocol string
	address  string
	port     string
	database string
}

func main() {

	myConn := &dbConnection{
		os.Getenv("APPDB_USER"),
		os.Getenv("APPDB_PASS"),
		os.Getenv("APPDB_PORT_3306_TCP_PROTO"),
		os.Getenv("APPDB_SERVICE_HOST"),
		os.Getenv("APPDB_SERVICE_PORT"),
		os.Getenv("APPDB_NAME"),
	}

	ds := fmt.Sprintf("%s:%s@%s(%s:%s)/%s", myConn.user, myConn.password, myConn.protocol, myConn.address, myConn.port, myConn.database)

	db, err := sql.Open(driver, ds)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Print("Connection successful.")
	}

	// Create the table
	_, err = db.Exec(createTableStatement)
	if err != nil {
		errMessage := err.Error()
		tableExistsCode := "1050"
		re := regexp.MustCompile(tableExistsCode)
		if re.FindString(errMessage) == tableExistsCode {
			log.Print("Table already exists, exiting.")
			os.Exit(0) // If the table exists, we're done
		} else {
			log.Fatal(err)
		}
	}

	// Insert rows, this need more checks
	for _, statement := range insertDataStatements {
		res, err := db.Exec(statement)
		if err != nil {
			log.Fatal(err)
		} else {
			id, _ := res.LastInsertId()
			rows, _ := res.RowsAffected()
			log.Printf("Last insert id: %v, Rows affected: %v", id, rows)
		}
	}

	log.Print("Database initialization completed.")
}
