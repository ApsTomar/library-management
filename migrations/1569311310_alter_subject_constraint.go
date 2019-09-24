package migrations

import migrate "github.com/rubenv/sql-migrate"

func init() {
	instance.add(&migrate.Migration{
		Id: "1569311310",
		Up: []string{
			`
			ALTER TABLE subject_x_book DROP FOREIGN KEY subject_x_book_ibfk_3;
			`,
			`
			ALTER TABLE subject_x_book ADD CONSTRAINT FOREIGN KEY(subject_id) REFERENCES subject(id) ON DELETE CASCADE;
			`,
		},
		//language=SQL
		Down: []string{
		},
	})
}
