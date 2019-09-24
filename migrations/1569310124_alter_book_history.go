package migrations

import migrate "github.com/rubenv/sql-migrate"

func init() {
	instance.add(&migrate.Migration{
		Id: "1569310124",
		Up: []string{
			`
			ALTER TABLE book_history ADD CONSTRAINT FOREIGN KEY(user_id) REFERENCES account(id) ON DELETE CASCADE;
			`,
			`
			ALTER TABLE book_history ADD CONSTRAINT FOREIGN KEY(book_id) REFERENCES book(id) ON DELETE CASCADE;
			`,
		},
		//language=SQL
		Down: []string{
		},
	})
}
