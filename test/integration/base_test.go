package integration_test

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"testing"
// 	"time"

// 	_ "github.com/lib/pq"
// 	"github.com/ory/dockertest/v3"
// 	"github.com/ory/dockertest/v3/docker"
// )

// var db *sql.DB

// func TestMain(m *testing.M) {
// 	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
// 	pool, err := dockertest.NewPool("")
// 	if err != nil {
// 		log.Fatalf("Could not construct pool: %s", err)
// 	}

// 	err = pool.Client.Ping()
// 	if err != nil {
// 		log.Fatalf("Could not connect to Docker: %s", err)
// 	}

// 	// pulls an image, creates a container based on it and runs it
// 	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
// 		Repository: "postgres",
// 		Tag:        "15.4",
// 		Env: []string{
// 			"POSTGRES_PASSWORD=12345678",
// 			"POSTGRES_USER=testingdb",
// 			"POSTGRES_DB=testingdb",
// 			"listen_addresses = '*'",
// 		},
// 	}, func(config *docker.HostConfig) {
// 		// set AutoRemove to true so that stopped container goes away by itself
// 		config.AutoRemove = true
// 		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
// 	})
// 	if err != nil {
// 		log.Fatalf("Could not start resource: %s", err)
// 	}

// 	hostAndPort := resource.GetHostPort("5432/tcp")
// 	databaseUrl := fmt.Sprintf("postgres://testingdb:12345678@%s/testingdb?sslmode=disable", hostAndPort)

// 	log.Println("Connecting to database on url: ", databaseUrl)

// 	resource.Expire(20) // Tell docker to hard kill the container in 120 seconds

// 	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
// 	pool.MaxWait = 10 * time.Second
// 	if err = pool.Retry(func() error {
// 		db, err = sql.Open("postgres", databaseUrl)
// 		if err != nil {
// 			return err
// 		}
// 		return db.Ping()
// 	}); err != nil {
// 		log.Fatalf("Could not connect to docker: %s", err)
// 	}

// 	defer func() {
// 		if err := pool.Purge(resource); err != nil {
// 			log.Fatalf("Could not purge resource: %s", err)
// 		}
// 	}()

// 	// run tests
// 	m.Run()
// }

// func TestRealbob(t *testing.T) {
// 	// all tests
// 	fmt.Println("hi")
// }
