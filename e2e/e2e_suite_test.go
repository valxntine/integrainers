package e2e

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/valxntine/integrainers/containers"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestEndToEnd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2E Test Suite")
}

var (
	db *sql.DB
)

const (
	bookHost = "http://book-service:8081"
)

var _ = BeforeSuite(func() {
	var err error

	ctx := context.Background()

	_, err = containers.StartTestContainers(ctx)

	db, err = sql.Open("mysql", MakeConnectionString())
	Expect(err).ToNot(HaveOccurred())

	_, err = db.Exec("TRUNCATE book;")
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	Expect(db.Close()).To(Succeed())
})

func MakeConnectionString() string {
	// DOCKER
	u := os.Getenv("DB_USER")
	pw := os.Getenv("DB_PASSWORD")
	d := os.Getenv("DB_NAME")
	p := os.Getenv("DB_PORT")
	h := os.Getenv("DB_HOST")

	// Set defaults for local testing against the mysql container started by docker_up
	if u == "" {
		u = "book_db"
	}
	if pw == "" {
		pw = "book_db"
	}
	if d == "" {
		d = "book_db"
	}
	if p == "" {
		p = "3306"
	}
	if h == "" {
		h = "127.0.0.1"
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", u, pw, h, p, d)
}
