package migrations

import migrate "github.com/rubenv/sql-migrate"

func init() {
	instance.add(&migrate.Migration{
		Id: "1569239890",
		Up: []string{
			`
			ALTER TABLE book MODIFY author_id varchar(255);
			`,
		},
		//language=SQL
		Down: []string{
		},
	})
}
