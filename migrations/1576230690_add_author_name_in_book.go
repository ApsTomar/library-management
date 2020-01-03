package migrations

import migrate "github.com/rubenv/sql-migrate"

func init() {
	instance.add(&migrate.Migration{
		Id: "1576230690",
		Up: []string{
			`
			ALTER TABLE book
				ADD COLUMN author_name varchar(255) NOT NULL;
			`,
		},
		//language=SQL
		Down: []string{
		},
	})
}
