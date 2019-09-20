package migrations

import migrate "github.com/rubenv/sql-migrate"

func init() {
	instance.add(&migrate.Migration{
		Id: "1568808227",
		Up: []string{
			`
			CREATE TABLE account_type (
			  id bigint(20) NOT NULL AUTO_INCREMENT,
			  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
			  updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			  deleted_at timestamp NULL DEFAULT NULL,
			  role varchar(255) NOT NULL,
			  PRIMARY KEY (id),
			  UNIQUE KEY id (id),
			);
           `,
			`
            CREATE TABLE account (
			  id bigint(20) NOT NULL AUTO_INCREMENT,
			  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
			  updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			  deleted_at timestamp NULL DEFAULT NULL,
			  name varchar(255) NOT NULL,
			  email varchar(255) NOT NULL,
			  account_role_id bigint(20) NOT NULL,
			  password_hash varchar(255) NOT NULL,
			  PRIMARY KEY (id),
			  UNIQUE KEY id (id),
			  UNIQUE KEY email (email),
			  FOREIGN KEY (account_role_id) REFERENCES account_type (id) ON DELETE CASCADE
			);
			`,
			`
			CREATE TABLE author (
			  id bigint(20) NOT NULL AUTO_INCREMENT,
			  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
			  updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			  deleted_at timestamp NULL DEFAULT NULL,
			  name varchar(255) NOT NULL,
			  date_of_birth timestamp NOT NULL,
			  PRIMARY KEY (id),
			  UNIQUE KEY id (id),
			  UNIQUE KEY name_dob (name,date_of_birth),
			  FOREIGN KEY (category_id) REFERENCES category (id) ON DELETE CASCADE
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
			  available bool DEFAULT true,
			  available_date DEFAULT CURRENT_TIMESTAMP,
			  PRIMARY KEY (id),
			  UNIQUE KEY id (id),
			  FOREIGN KEY (category_id) REFERENCES category (id) ON DELETE CASCADE
			);
           `,
			`
			CREATE TABLE book_x_author (
			  book_id bigint(20) NOT NULL,
			  author_id bigint(20) NOT NULL,
			  PRIMARY KEY (book_id,author_id),
			  FOREIGN KEY (author_id) REFERENCES author (id) ON DELETE CASCADE,
			  FOREIGN KEY (book_id) REFERENCES book (id) ON DELETE CASCADE,
			);
           `,
		},
		//language=SQL
		Down: []string{
			`DROP TABLE account_type;`,
			`DROP TABLE account;`,
			`DROP TABLE author;`,
			`DROP TABLE book;`,
			`DROP TABLE book_x_author`,
		},
	})
}
