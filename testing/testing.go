package testing

import (
	"database/sql"
	"io/ioutil"
	"log"
	"net"
	"net/url"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	// imports the postgres sql driver
	_ "github.com/lib/pq"
)

const dbName = "minesweeper_test"
const dbUser = "minesweeper"

var testingDB *sql.DB
var pgURL *url.URL

// NullLogger builds a logger that discards every log
func NullLogger() *log.Logger {
	return log.New(ioutil.Discard, "", log.Ltime)
}

// DB connects to the testing database
func DB() (*sql.DB, error) {
	var err error
	if testingDB == nil {
		testingDB, err = sql.Open("postgres", pgURL.String())
	}
	return testingDB, err
}

// InitializeDB initializes spins up a clean postgres to run tests.
func InitializeDB(m *testing.M) {
	code := 0
	defer func() {
		os.Exit(code)
	}()

	log := log.New(os.Stdout, "", log.Ltime)

	schemaRaw, err := ioutil.ReadFile("../schema.sql")
	if err != nil {
		log.Fatalf("could read schema.sql file %s", err)
	}
	schemaSQL := string(schemaRaw[0:])

	pgURL = &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(dbUser, dbName),
		Path:   dbName,
	}
	q := pgURL.Query()
	q.Add("sslmode", "disable")
	pgURL.RawQuery = q.Encode()

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker %s", err)
	}

	pw, _ := pgURL.User.Password()
	runOpts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "11.2-alpine",
		Env: []string{
			"POSTGRES_USER=" + pgURL.User.Username(),
			"POSTGRES_PASSWORD=" + pw,
			"POSTGRES_DB=" + pgURL.Path,
		},
	}

	resource, err := pool.RunWithOptions(&runOpts)
	if err != nil {
		log.Fatalf("could start postgres container: %s", err)
	}
	defer func() {
		err = pool.Purge(resource)
		if err != nil {
			log.Fatalf("could not purge resource: %s", err)
		}
	}()

	pgURL.Host = resource.Container.NetworkSettings.IPAddress

	// Docker layer network is different on Mac
	if runtime.GOOS == "darwin" {
		pgURL.Host = net.JoinHostPort(resource.GetBoundIP("5432/tcp"), resource.GetPort("5432/tcp"))
	}

	logWaiter, err := pool.Client.AttachToContainerNonBlocking(docker.AttachToContainerOptions{
		Container:    resource.Container.ID,
		OutputStream: log.Writer(),
		ErrorStream:  log.Writer(),
		Stderr:       true,
		Stdout:       true,
		Stream:       true,
	})
	if err != nil {
		log.Fatalf("could not connect to postgres container log output: %s", err)
	}
	defer func() {
		err = logWaiter.Close()
		if err != nil {
			log.Fatalf("could not close container log: %s", err)
		}
		err = logWaiter.Wait()
		if err != nil {
			log.Fatalf("could not wait for container log to close: %s", err)
		}
	}()

	pool.MaxWait = 10 * time.Second
	err = pool.Retry(func() error {
		db, err := sql.Open("postgres", pgURL.String())
		if err != nil {
			return err
		}
		// create the schema
		log.Printf("initializing database\n")
		_, err = db.Exec(schemaSQL)
		return err
	})
	if err != nil {
		log.Fatalf("could not connect to postgres server: %s", err)
	}
	code = m.Run()
}
