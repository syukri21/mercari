package postgre

const CreateUser = `
	INSERT INTO users
	(username, email, password, activate_key, is_activated, created_at, updated_at)
	VALUES (:username, :email, :password, :activate_key, :is_activated,:created_at, :updated_at);
`
const GetUserPinByEmail = `SELECT activate_key FROM users WHERE email = $1 `

const ActivateUser = `
	UPDATE users
	SET is_activated=true, updated_at=:updated_at
	WHERE email=:email;
`

const GetUser = `SELECT username, email, password FROM users WHERE email = $1 AND is_activated = true`

const CreateLoginHistory = `INSERT INTO login_history
	(username, email, device_id,  created_at, updated_at,  login_at)
	VALUES (:username, :email, :device_id,  :created_at, :updated_at,  :login_at);`

const GetLoginHistories = `SELECT * FROM login_history WHERE email = :email   ORDER BY login_at DESC LIMIT :limit OFFSET :offset `

const DDLCreateTableUser = `
CREATE TABLE users (
	id serial PRIMARY KEY,
	username VARCHAR ( 50 ) NOT NULL,
	password TEXT NOT NULL,
	email VARCHAR ( 255 ) UNIQUE NOT NULL,
	activate_key VARCHAR (10) NOT NULL ,
	is_activated BOOLEAN NOT NULL DEFAULT false,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
`

const DDLCreateTableUserHistory = `
CREATE TABLE login_history (
	id serial PRIMARY KEY,
	username VARCHAR ( 50 ) NOT NULL,
	email VARCHAR ( 255 ) NOT NULL,
	device_id VARCHAR ( 255 ) NOT NULL,
	created_at TIMESTAMP NOT NULL,
	login_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP
);
`
