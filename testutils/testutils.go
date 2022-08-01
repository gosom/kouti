//go:build test
// +build test

package testutils

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	TestDbName   = "test"
	TestDbUser   = "test"
	TestDbPasswd = "test"
)

func GetTestDSN(port string) string {
	dsn := "postgres://%s:%s@localhost:%s/test?sslmode=disable&pool_max_conns=10"
	return fmt.Sprintf(dsn, TestDbName, TestDbPasswd, port)
}

func SpinPostgresContainer(ctx context.Context) (testcontainers.Container, string) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:14.4-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       TestDbUser,
			"POSTGRES_USER":     TestDbUser,
			"POSTGRES_PASSWORD": TestDbPasswd,
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections"),
	}

	postgresContainer, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{ContainerRequest: req, Started: true},
	)
	if err != nil {
		panic(err)
	}
	p, _ := postgresContainer.MappedPort(ctx, "5432")
	// added this sleep here :(
	// better to find a way to proper wait
	time.Sleep(5 * time.Second)
	return postgresContainer, p.Port()
}
