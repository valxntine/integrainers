package e2e

import (
	"database/sql"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"io"
	"net/http"
	"strings"
)

var _ = Describe("Invalidate", func() {

	var (
		db  *sql.DB
		err error
	)
	db, err = sql.Open("mysql", MakeConnectionString())
	if err != nil {
		fmt.Errorf("WOOPS -> %w", err)
	}

	BeforeEach(func() {
		_, err = db.Exec("TRUNCATE book;")
		Expect(err).ToNot(HaveOccurred())

		_, err = db.Exec(`INSERT INTO book(id, iban, name, author, pages) VALUES(1, "312", "Valentines Book", "Valentine", 321)`)
		Expect(err).ToNot(HaveOccurred())
	})

	It("updates the rescan status and notifies shopper", func() {
		url := fmt.Sprintf("%s/api/v1/books/", bookHost)
		req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(fmt.Sprintf(`{"iban":"%s"}`, "123")))

		resp, err := http.DefaultClient.Do(req)
		Expect(err).ToNot(HaveOccurred())

		defer resp.Body.Close()

		Expect(resp.StatusCode).To(Equal(200))

		res := db.QueryRow(fmt.Sprintf("SELECT iban, name, author, pages FROM book WHERE id=%d", 2))
		var iban string
		var name string
		var author string
		var pages int32

		err = res.Scan(&iban, &name, &author, &pages)
		Expect(err).ToNot(HaveOccurred())
		Expect(iban).To(Equal("123"))
		Expect(name).To(Equal("Amys book"))
		Expect(author).To(Equal("Amy Bull"))
		Expect(pages).To(Equal("666"))

		expectedJson := `{
				 "message": "Successfullu saved book"
				 "iban": 123
				}`

		responsePayload, err := io.ReadAll(resp.Body)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(responsePayload)).To(MatchJSON(expectedJson))

	})
})
