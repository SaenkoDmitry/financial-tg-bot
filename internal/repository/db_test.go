package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func MigrateDb(dbURI string) error {
	db, err := sql.Open("postgres", dbURI)
	if err != nil {
		return err
	}

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "./../../migrations"); err != nil { // base schema and data
		panic(err)
	}

	if err := goose.Up(db, "./../../test_data"); err != nil { // test data
		panic(err)
	}

	return db.Close()
}

func SetupTestDatabase() (testcontainers.Container, *pgxpool.Pool) {
	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env: map[string]string{
			"POSTGRES_DB":       "route256",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_USER":     "postgres",
		},
	}

	// 2. Start PostgreSQL container
	dbContainer, err := testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		})
	if err != nil {
		log.Fatal(err)
	}

	// 3.1 Get host and port of PostgreSQL container
	host, _ := dbContainer.Host(context.Background())
	port, _ := dbContainer.MappedPort(context.Background(), "5432")

	// 3.2 Create db connection string and connect
	dbURI := fmt.Sprintf("postgres://postgres:postgres@%v:%v/route256?sslmode=disable", host, port.Port())
	connPool, _ := pgxpool.New(context.Background(), dbURI)

	// todo : generate data ?
	_ = MigrateDb(dbURI)

	return dbContainer, connPool
}
