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

type answers struct {
	id       int
	sentence string
}

type dbConnection struct {
	user     string
	password string
	protocol string
	address  string
	port     string
	database string
}

var (
	createTableStatement = "CREATE TABLE answers(id INT AUTO_INCREMENT, sentence VARCHAR(255), PRIMARY KEY(id));"
	insertDataStatements = []answers{
		{1, "INSERT INTO answers (id, sentence) VALUES (?, 'As I see it, yes');"},
		{2, "INSERT INTO answers (id, sentence) VALUES (?, 'It is certain');"},
		{3, "INSERT INTO answers (id, sentence) VALUES (?, 'It is decidedly so');"},
		{4, "INSERT INTO answers (id, sentence) VALUES (?, 'Most likely');"},
		{5, "INSERT INTO answers (id, sentence) VALUES (?, 'Yes');"},
		{6, "INSERT INTO answers (id, sentence) VALUES (?, 'Ask again later');"},
		{7, "INSERT INTO answers (id, sentence) VALUES (?, 'Better not tell you now');"},
		{8, "INSERT INTO answers (id, sentence) VALUES (?, 'Cannot predict now');"},
		{9, "INSERT INTO answers (id, sentence) VALUES (?, 'Reply hazy, try again');"},
		{10, "INSERT INTO answers (id, sentence) VALUES (?, 'Concentrate and ask again');"},
		{11, "INSERT INTO answers (id, sentence) VALUES (?, 'Do not count on it');"},
		{12, "INSERT INTO answers (id, sentence) VALUES (?, 'My reply is no');"},
		{13, "INSERT INTO answers (id, sentence) VALUES (?, 'My sources say no');"},
		{14, "INSERT INTO answers (id, sentence) VALUES (?, 'Outlook not so good');"},
		{15, "INSERT INTO answers (id, sentence) VALUES (?, 'Very doubtful');"},
	}
)

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
		res, err := db.Exec(statement.sentence, statement.id)
		if err != nil {
			log.Fatal(err)
		} else {
			rows, _ := res.RowsAffected()
			log.Printf("Last insert id: %v, Rows affected: %v", statement.id, rows)
		}
	}

	log.Print("Database initialization completed.")
}
