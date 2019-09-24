package migrations

import migrate "github.com/rubenv/sql-migrate"

func init() {
	instance.add(&migrate.Migration{
		Id: "1569241875",
		Up: []string{
			`
			ALTER TABLE subject ADD CONSTRAINT UNIQUE KEY name (name);
			`,
		},
		//language=SQL
		Down: []string{
		},
	})
}
