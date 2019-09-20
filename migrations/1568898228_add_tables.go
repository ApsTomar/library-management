package migrations

import migrate "github.com/rubenv/sql-migrate"

func init() {
	instance.add(&migrate.Migration{
		Id: "1568898228",
		Up: []string{
			`
			CREATE TABLE subject (
			  id bigint(20) NOT NULL,
			  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
			  updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			  deleted_at timestamp NULL DEFAULT NULL,
			  name varchar(255) NOT NULL,
			  PRIMARY KEY (id),
			  UNIQUE KEY id (id),
			);
			`,
			`
			CREATE TABLE subject_x_book (
			  book_id bigint(20) NOT NULL,
			  subject_id bigint(20) NOT NULL,
			  PRIMARY KEY (book_id,subject_id),
			  FOREIGN KEY (subject_id) REFERENCES subject (id) ON DELETE CASCADE,
			  FOREIGN KEY (book_id) REFERENCES book (id) ON DELETE CASCADE,
			);
			`,
			`
			CREATE TABLE book_history (
			  book_id bigint(20) NOT NULL,
			  user_id bigint(20) NOT NULL,
			  issue_date timestamp,
			  return_date timestamp,
			  PRIMARY KEY (book_id,user_id),
			  FOREIGN KEY (user_id) REFERENCES user (id) ON DELETE CASCADE,
			  FOREIGN KEY (book_id) REFERENCES book (id) ON DELETE CASCADE,
			);
			`,
		},
		//language=SQL
		Down: []string{
			`DROP TABLE subject;`,
			`DROP TABLE subject_x_book;`,
			`DROP TABLE book_history;`,
		},
	})
}
