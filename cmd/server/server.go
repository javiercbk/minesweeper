package main

//go:generate sqlboiler psql

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/javiercbk/minesweeper/http"
)

const defaultLogFilePath = "minesweeper-server.log"
const defaultJWTSecret = "minesweep"
const defaultAddress = "0.0.0.0"
const defaultDBName = "minesweep"
const defaultDBUser = "minesweep"

func main() {
	db := sql.DB{}
	var logFilePath, address, jwtSecret, dbName, dbHost, dbUser, dbPass string
	flag.StringVar(&logFilePath, "l", defaultLogFilePath, "the log file location")
	flag.StringVar(&address, "a", defaultAddress, "the http server address")
	flag.StringVar(&jwtSecret, "jwt", defaultJWTSecret, "the jwt secret")
	flag.StringVar(&dbName, "dbn", defaultDBName, "the database name")
	flag.StringVar(&dbHost, "dbh", defaultDBUser, "the database host")
	flag.StringVar(&dbUser, "dbu", "", "the database user")
	flag.StringVar(&dbPass, "dbp", "", "the database password")
	flag.Parse()
	logFile, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("error opening lof file: %s", err)
		os.Exit(1)
	}
	defer logFile.Close()
	logger := log.New(logFile, "applog: ", log.Lshortfile|log.LstdFlags)
	err = connectPostgres(dbName, dbHost, dbUser, dbPass, &db)
	if err != nil {
		logger.Printf("error connecting to postgres: %s", err)
		os.Exit(1)
	}
	cnf := http.Config{
		Address:   address,
		JWTSecret: jwtSecret,
	}
	err = http.Serve(cnf, logger, &db)
	if err != nil {
		logger.Fatalf("could not start server %s\n", err)
	}
}

func connectPostgres(dbName, dbHost, dbUser, dbPass string, db *sql.DB) error {
	var err error
	postgresOpts := fmt.Sprintf("dbname=%s host=%s user=%s password=%s", dbName, dbHost, dbUser, dbPass)
	db, err = sql.Open("postgres", postgresOpts)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	return nil
}
