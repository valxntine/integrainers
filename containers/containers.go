package containers

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
	"path/filepath"
)

type MySQLContainer struct {
	testcontainers.Container
	username string
	password string
	database string
}

type ServiceContainer struct {
	testcontainers.Container
}

type TestContainers struct {
	DB          *MySQLContainer
	MockService *ServiceContainer
	Service     *ServiceContainer
}

func CreateBookDB(ctx context.Context, n *testcontainers.DockerNetwork, aliases []string) (*MySQLContainer, error) {
	//TODO: ADD SCHEMA
	schema := filepath.Join("db", "initdb.sql")
	req := testcontainers.ContainerRequest{
		Name:         "book_db",
		Image:        "mysql/mysql-server:8.0.28",
		ExposedPorts: []string{"3306/tcp", "33060/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "book_db",
			"MYSQL_USER":          "book_db",
			"MYSQL_PASSWORD":      "book_db",
			"MYSQL_DATABASE":      "book_db",
		},
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      schema,
				ContainerFilePath: "/docker-entrypoint-initdb.d/" + filepath.Base(schema),
				FileMode:          0o755,
			},
		},
		WaitingFor: wait.ForListeningPort("3306/tcp"),
		Networks: []string{
			n.Name,
		},
		NetworkAliases: map[string][]string{
			n.Name: aliases,
		},
	}

	username := req.Env["MYSQL_USER"]
	password := req.Env["MYSQL_PASSWORD"]
	database := req.Env["MYSQL_DATABASE"]

	genericConReq := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	}

	container, err := testcontainers.GenericContainer(ctx, genericConReq)
	if err != nil {
		return &MySQLContainer{}, err
	}

	return &MySQLContainer{
		Container: container,
		username:  username,
		password:  password,
		database:  database,
	}, nil
}

func NewMockServiceContainer(ctx context.Context, n *testcontainers.DockerNetwork, aliases []string) (*ServiceContainer, error) {
	req := testcontainers.ContainerRequest{
		Name:  "library_service",
		Image: "mockserver/mockserver:mockserver-5.14.0",
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      "../e2e/mocks/library.json",
				ContainerFilePath: "/config/library.json",
				FileMode:          0o700,
			},
		},
		ExposedPorts: []string{"1080/tcp"},
		Env: map[string]string{
			"MOCKSERVER_LOG_LEVEL":                "DEBUG!",
			"MOCKSERVER_SERVER_PORT":              "1080",
			"MOCKSERVER_INITIALIZATION_JSON_PATH": "/config/library.json",
		},
		WaitingFor: wait.ForListeningPort("1080/tcp"),
		Networks: []string{
			n.Name,
		},
		NetworkAliases: map[string][]string{
			n.Name: aliases,
		},
	}

	genericConReq := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	}

	container, err := testcontainers.GenericContainer(ctx, genericConReq)
	if err != nil {
		return &ServiceContainer{}, err
	}

	return &ServiceContainer{
		Container: container,
	}, nil
}

func CreateServiceFromDocker(ctx context.Context, n *testcontainers.DockerNetwork, aliases []string) (*ServiceContainer, error) {
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			FromDockerfile: testcontainers.FromDockerfile{
				Context:       "../",
				Dockerfile:    "Dockerfile",
				PrintBuildLog: true,
				KeepImage:     false,
				BuildOptionsModifier: func(buildOptions *types.ImageBuildOptions) {
					buildOptions.Target = "test"
				},
			},
			Networks: []string{
				n.Name,
			},
			NetworkAliases: map[string][]string{
				n.Name: aliases,
			},
			ExposedPorts: []string{"8081/tcp"},
		},
		Started: true,
	})
	if err != nil {
		return nil, err
	}
	return &ServiceContainer{
		c,
	}, nil
}

func CreateDockerNetwork(ctx context.Context) (*testcontainers.DockerNetwork, error) {
	n, err := network.New(ctx)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func StartTestContainers(ctx context.Context) (*TestContainers, error) {
	n, err := CreateDockerNetwork(ctx)
	if err != nil {
		panic(err)
	}

	db, err := CreateBookDB(ctx, n, []string{"book-db"})
	if err != nil {
		panic(err)
	}

	mockSvc, err := NewMockServiceContainer(ctx, n, []string{"library-mock"})
	if err != nil {
		panic(err)
	}

	svc, err := CreateServiceFromDocker(ctx, n, []string{"book-service"})
	if err != nil {
		panic(err)
	}

	return &TestContainers{
		DB:          db,
		MockService: mockSvc,
		Service:     svc,
	}, nil
}
