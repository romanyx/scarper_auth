DROP TYPE IF EXISTS user_status;
CREATE TYPE user_status AS ENUM ('NEW', 'VERIFIED');

CREATE TABLE IF NOT EXISTS users (
	id 		SERIAL PRIMARY KEY,
	account_id	CHAR (36) UNIQUE NOT NULL, /* global identifier */

	email 			VARCHAR (50) NOT NULL,
	password_hash 		VARCHAR (72) NOT NULL,
	status 			user_status NOT NULL,
	token		 	CHAR (32) NOT NULL, 

	/* timestamps */
	token_expired_at		TIMESTAMP NOT NULL,
	created_at			TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at			TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS users_email_idx ON users (email);
CREATE INDEX IF NOT EXISTS users_verification_token_idx ON users (verification_token);

CREATE TABLE IF NOT EXISTS resets (
	id	SERIAL PRIMARY KEY,

	token	VARCHAR (72) NOT NULL,

	user_id	INTEGER REFERENCES users (id), 

	/* timestamps */
	token_expired_at	TIMESTAMP NOT NULL,
	created_at		TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at		TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS resets_token_idx ON resets (token);