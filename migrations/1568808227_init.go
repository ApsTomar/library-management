package migrations

import migrate "github.com/rubenv/sql-migrate"

func init() {
	instance.add(&migrate.Migration{
		Id: "1568808227",
		Up: []string{
			`
            CREATE TABLE account (
			  id bigint(20) NOT NULL AUTO_INCREMENT,
			  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
			  updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			  deleted_at timestamp NULL DEFAULT NULL,
			  name varchar(255) NOT NULL,
			  email varchar(255) NOT NULL,
			  account_role varchar(255) NOT NULL,
			  password_hash varchar(255) NOT NULL,
			  PRIMARY KEY (id),
			  UNIQUE KEY id (id),
			  UNIQUE KEY email (email)
			);
			`,
			`
			CREATE TABLE author (
			  id bigint(20) NOT NULL AUTO_INCREMENT,
			  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
			  updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			  deleted_at timestamp NULL DEFAULT NULL,
			  name varchar(255) NOT NULL,
			  date_of_birth varchar(255) NOT NULL,
			  PRIMARY KEY (id),
			  UNIQUE KEY id (id),
			  UNIQUE KEY name_dob (name,date_of_birth)
			);
           `,
			`
			CREATE TABLE book (
			  id bigint(20) NOT NULL AUTO_INCREMENT,
			  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
			  updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			  deleted_at timestamp NULL DEFAULT NULL,
			  name varchar(255) NOT NULL,
			  subject varchar(255) NOT NULL,
			  author_id text,
			  available bool DEFAULT true,
			  available_date timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
			  PRIMARY KEY (id),
			  UNIQUE KEY id (id)
			);
           `,
			`
			CREATE TABLE book_x_author (
			  book_id bigint(20) NOT NULL,
			  author_id bigint(20) NOT NULL,
			  PRIMARY KEY (book_id,author_id),
			  FOREIGN KEY (author_id) REFERENCES author (id) ON DELETE CASCADE,
			  FOREIGN KEY (book_id) REFERENCES book (id) ON DELETE CASCADE
			);
           `,
		},
		//language=SQL
		Down: []string{
			`DROP TABLE account;`,
			`DROP TABLE author;`,
			`DROP TABLE book;`,
			`DROP TABLE book_x_author`,
		},
	})
}
