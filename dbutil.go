package main

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "log"
    "os"
    "fmt"
)

var db *sql.DB

//Setup the DB. If it doesn't exist, create it.
func Setup() {
    if !Exists("data.db") {
      	log.Println("Table not found - creating it.")
        defer CreateTable()
    }
    log.Println("Opening connection to the database.")
    database, err := sql.Open("sqlite3", "data.db")
    if err != nil {
        log.Fatal(err)
    }
    db = database
}

//Check if the file at the specified path exists.
func Exists(path string) bool {
    _, err := os.Stat(path)
    if err == nil {
        return true
    }
    if os.IsNotExist(err) {
        return false
    } else {
        return true
    }
}

//Create the necessary tables.
func CreateTable() {
    log.Println("Creating new table.")
    _, err := db.Exec("CREATE TABLE keys (key VARCHAR(14));")
    if err != nil {
        log.Fatal(err)
    }
    _, err = db.Exec("CREATE TABLE links (key VARCHAR(14), shorten VARCHAR(4), url VARCHAR(80), clicks INT);")
    if err != nil {
        log.Fatal(err)
    }

    key := GenRandomShorten(14)
    query := fmt.Sprintf("INSERT INTO keys VALUES (\"%s\")", key)
    output := fmt.Sprintf("Generated a new key - %s - KEEP THIS SAFE.", key)
    log.Println(output)
    Execute(query)
}

//Execute the specified statement in the DB.
func Execute(stmt string) {
    _, err := db.Exec(stmt)
    if err != nil {
        log.Fatal(err)
    }
}

//Returns the DB in use.
func Database() *sql.DB {
    return db
}

//Close the connection to the Database.
func Close() {
    if db != nil {
        db.Close()
    }
}
