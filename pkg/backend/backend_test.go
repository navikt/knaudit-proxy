package backend_test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	kp "github.com/navikt/knaudit-proxy/pkg/backend"

	dockertest "github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/require"
)

var backend *kp.OracleBackend //nolint: gochecknoglobals

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("creating pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("pinging docker: %s", err)
	}

	pool.MaxWait = 30 * time.Minute

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("getting working directory: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "gvenzl/oracle-free",
		Tag:        "slim-faststart",
		Env: []string{
			"ORACLE_PASSWORD=testpass",
		},
		Mounts: []string{
			fmt.Sprintf("%s/../../resources/container-entrypoint-initdb.d:/container-entrypoint-initdb.d/", wd),
		},
		ExposedPorts: []string{"1521/tcp"},
		Platform:     "linux/amd64",
	})
	if err != nil {
		log.Fatalf("running container: %s", err)
	}

	dsn := fmt.Sprintf("oracle://system:testpass@localhost:%s/FREEPDB1", resource.GetPort("1521/tcp"))

	_ = resource.Expire(240)

	backend = kp.NewOracleBackend(dsn)

	pool.MaxWait = 240 * time.Second
	if err = pool.Retry(func() error {
		log.Printf("trying to connect to database: %s", dsn)

		err = backend.Open()
		if err != nil {
			return fmt.Errorf("could not open database: %w", err)
		}

		err = backend.Ping()
		if err != nil {
			return fmt.Errorf("could not ping database: %w", err)
		}

		return nil
	}); err != nil {
		log.Fatalf("could not connect to database: %s", err)
	}

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestNewOracleBackend(t *testing.T) {
	t.Skipf("Test not implemented")

	t.Parallel()

	err := backend.Send("test")
	require.NoError(t, err)
}
