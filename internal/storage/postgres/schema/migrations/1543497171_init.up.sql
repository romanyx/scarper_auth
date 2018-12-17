DROP TYPE IF EXISTS user_status;
CREATE TYPE user_status AS ENUM ('NEW', 'VERIFIED');

CREATE TABLE IF NOT EXISTS users (
	id 		SERIAL PRIMARY KEY,
	account_id	CHAR (36) UNIQUE NOT NULL, /* global identifier */

	email 			VARCHAR (50) NOT NULL,
	password_hash 		VARCHAR (72) NOT NULL,
	status 			user_status NOT NULL,
	token		 	CHAR (36) UNIQUE, 

	/* timestamps */
	created_at			TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at			TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS users_email_idx ON users (email);
CREATE INDEX IF NOT EXISTS users_token_idx ON users (token);

CREATE TABLE IF NOT EXISTS resets (
	id	SERIAL PRIMARY KEY,

	token	VARCHAR (36) UNIQUE NOT NULL,

	user_id	INTEGER REFERENCES users (id), 

	/* timestamps */
	expired_at		TIMESTAMP NOT NULL,
	created_at		TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at		TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS resets_token_idx ON resets (token);
