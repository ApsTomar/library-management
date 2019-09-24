package migrations

import migrate "github.com/rubenv/sql-migrate"

func init() {
	instance.add(&migrate.Migration{
		Id: "1569247400",
		Up: []string{
			`
			ALTER TABLE book_history
				ADD COLUMN returned boolean DEFAULT false;
			`,
		},
		//language=SQL
		Down: []string{
		},
	})
}
