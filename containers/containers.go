package containers

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
	"path/filepath"
	"time"
)

type logger struct{}

func (l logger) Accept(log testcontainers.Log) {
	fmt.Println(string(log.Content))
}

type MySQLContainer struct {
	testcontainers.Container
	Username   string
	Password   string
	Database   string
	Connection string
}

type ServiceContainer struct {
	testcontainers.Container
	URI string
}

type TestContainers struct {
	DB          *MySQLContainer
	MockService *ServiceContainer
	Service     *ServiceContainer
}

func CreateBookDB(ctx context.Context, n *testcontainers.DockerNetwork, aliases []string) (*MySQLContainer, error) {
	schema := filepath.Join("../", "db", "initdb.sql")
	req := testcontainers.ContainerRequest{
		Name:         "book_db",
		Image:        "mysql:8",
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

	//mysqlContainer, err := mysql.RunContainer(ctx,
	//	testcontainers.WithImage("mysql:8"),
	//	mysql.WithDatabase("book_db"),
	//	mysql.WithUsername("book_db"),
	//	mysql.WithPassword("book_db"),
	//	mysql.WithScripts(schema))
	//
	//if err != nil {
	//	panic(err)
	//}

	containerPort, err := container.MappedPort(ctx, "3306/tcp")
	if err != nil {
		return &MySQLContainer{}, err
	}

	host, err := container.Host(ctx)
	if err != nil {
		return &MySQLContainer{}, err
	}

	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, containerPort.Port(), database)

	return &MySQLContainer{
		Container:  container,
		Connection: conn,
	}, nil
}

func NewMockServiceContainer(ctx context.Context, n *testcontainers.DockerNetwork, aliases []string) (*ServiceContainer, error) {
	mockPath, _ := filepath.Abs(filepath.Join("..", "e2e", "mocks", "library.json"))
	req := testcontainers.ContainerRequest{
		Name:  "library_service",
		Image: "mockserver/mockserver:mockserver-5.14.0",
		Mounts: testcontainers.ContainerMounts{
			{
				Source: testcontainers.GenericVolumeMountSource{Name: "config"},
				Target: "/config/",
			},
		},
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      mockPath,
				ContainerFilePath: "/config/library.json",
				FileMode:          0o777,
			},
		},
		ExposedPorts: []string{"1080/tcp"},
		Env: map[string]string{
			"MOCKSERVER_LOG_LEVEL":                "DEBUG",
			"MOCKSERVER_SERVER_PORT":              "1080",
			"MOCKSERVER_INITIALIZATION_JSON_PATH": "/config/library.json",
		},
		WaitingFor: wait.ForLog("started on port: 1080"),
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
			},
			WaitingFor: wait.ForLog("App is running").WithPollInterval(10 * time.Second),
			LogConsumerCfg: &testcontainers.LogConsumerConfig{
				Opts:      []testcontainers.LogProductionOption{testcontainers.WithLogProductionTimeout(10 * time.Second)},
				Consumers: []testcontainers.LogConsumer{&logger{}},
			},
			Networks: []string{
				n.Name,
			},
			NetworkAliases: map[string][]string{
				n.Name: aliases,
			},
			ExposedPorts: []string{"8088/tcp"},
		},
		Started: true,
	})
	if err != nil {
		return nil, err
	}

	ip, _ := c.Host(ctx)

	port, _ := c.MappedPort(ctx, "8088")

	return &ServiceContainer{
		Container: c,
		URI:       fmt.Sprintf("http://%s:%s", ip, port.Port()),
	}, nil
}

func CreateDockerNetwork(ctx context.Context) (*testcontainers.DockerNetwork, error) {
	n, err := network.New(ctx,
		network.WithCheckDuplicate(),
		network.WithAttachable(),
		network.WithDriver("bridge"))
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

	db, err := CreateBookDB(ctx, n, []string{"book-db.app.internal"})
	if err != nil {
		panic(err)
	}

	mockSvc, err := NewMockServiceContainer(ctx, n, []string{"library-mock.app.internal"})
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
