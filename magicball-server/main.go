package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	driver = "mysql"
)

type dbConnection struct {
	user     string
	password string
	protocol string
	address  string
	port     string
	database string
}

func randomInt(min, max int) int {
	return rand.Intn(max-min) + min
}

func getSentence(db *sql.DB) (string, error) {
	var sentence string
	var count int

	err := db.QueryRow("SELECT COUNT(*) FROM answers").Scan(&count)
	if err != nil {
		return "", err
	}
	randomId := randomInt(1, count)

	err = db.QueryRow("SELECT sentence FROM answers WHERE id = ?", randomId).Scan(&sentence)
	if err != nil {
		return "", err
	}
	return sentence, nil
}

func toJson(s string) (string, error) {
	m := make(map[string]string)
	m["answer"] = s
	payload, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(payload), nil
}

func main() {

	listenPort := flag.String("p", "8080", "Service listening port")
	flag.Parse()

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

	defaultSocket := fmt.Sprintf("%s:%s", os.Getenv("POD_IP"), *listenPort)

	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		answer, err := getSentence(db)
		if err != nil {
			log.Fatal(err)
		}

		fullResponse := fmt.Sprintf("Magic 8 Ball said: %s.\n", answer)
		io.WriteString(w, fullResponse)
	})

	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {

		answer, err := getSentence(db)
		if err != nil {
			log.Fatal(err)
		}

		jsonAnswer, err := toJson(answer)
		if err != nil {
			log.Fatal(err)
		}

		fullResponse := fmt.Sprintf("%s\n", jsonAnswer)
		io.WriteString(w, fullResponse)
	})

	log.Fatal(http.ListenAndServe(defaultSocket, nil))
}
