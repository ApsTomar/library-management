package migrations

import migrate "github.com/rubenv/sql-migrate"

func init() {
	instance.add(&migrate.Migration{
		Id: "1569238933",
		Up: []string{
			`
			ALTER TABLE subject_x_book DROP FOREIGN KEY subject_x_book_ibfk_1;
			`,
			`
			ALTER TABLE subject MODIFY COLUMN id bigint(20) NOT NULL AUTO_INCREMENT;
			`,
			`
			ALTER TABLE subject_x_book ADD CONSTRAINT FOREIGN KEY(subject_id) REFERENCES subject(id);
			`,
		},
		//language=SQL
		Down: []string{
		},
	})
}
