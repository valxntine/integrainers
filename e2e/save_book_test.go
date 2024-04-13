package e2e

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"io"
	"net/http"
	"strings"
)

var _ = Describe("Invalidate", func() {

	BeforeEach(func() {
		_, err := db.Exec("TRUNCATE book;")
		Expect(err).ToNot(HaveOccurred())

		_, err = db.Exec(`INSERT INTO book(id, isbn, name, author, pages) VALUES(1, "312", "Valentines Book", "Valentine", 321)`)
		Expect(err).ToNot(HaveOccurred())
	})

	It("updates the rescan status and notifies shopper", func() {
		url := fmt.Sprintf("%s/api/v1/books", bookService)
		req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(fmt.Sprintf(`{"isbn": %d}`, 123)))
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		Expect(err).ToNot(HaveOccurred())

		responsePayload, err := io.ReadAll(resp.Body)
		Expect(err).NotTo(HaveOccurred())

		fmt.Printf("BODY -> %s\n", string(responsePayload))

		defer resp.Body.Close()

		Expect(resp.StatusCode).To(Equal(201))

		res := db.QueryRow(fmt.Sprintf("SELECT isbn, name, author, pages FROM book WHERE id=%d", 2))
		var isbn string
		var name string
		var author string
		var pages int32

		err = res.Scan(&isbn, &name, &author, &pages)
		Expect(err).ToNot(HaveOccurred())
		Expect(isbn).To(Equal("123"))
		Expect(name).To(Equal("Amys Book"))
		Expect(author).To(Equal("Amy Bull"))
		Expect(pages).To(Equal(int32(666)))

		expectedJson := `{
				 "message": "Successfully saved book",
				 "isbn": 123
				}`

		Expect(string(responsePayload)).To(MatchJSON(expectedJson))

	})
})
