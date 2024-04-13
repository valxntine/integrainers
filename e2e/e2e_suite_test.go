package e2e

import (
	"context"
	"database/sql"
	"github.com/valxntine/integrainers/containers"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestEndToEnd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2E Test Suite")
}

var (
	db          *sql.DB
	c           *containers.TestContainers
	bookService string
)

var _ = BeforeSuite(func() {
	var err error

	ctx := context.Background()

	c, err = containers.StartTestContainers(ctx)
	Expect(err).ToNot(HaveOccurred())

	bookService = c.Service.URI

	db, err = sql.Open("mysql", c.DB.Connection)
	Expect(err).ToNot(HaveOccurred())

	_, err = db.Exec("TRUNCATE book;")
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
})
