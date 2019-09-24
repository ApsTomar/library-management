package migrations

import migrate "github.com/rubenv/sql-migrate"

func init() {
	instance.add(&migrate.Migration{
		Id: "1569237627",
		Up: []string{
			`
			ALTER TABLE author MODIFY COLUMN date_of_birth varchar(255);
			`,
		},
		//language=SQL
		Down: []string{
		},
	})
}
