/* SQLEditor (Postgres)*/

/**
 DROP TABLE IF EXISTS
*/

DROP TABLE IF EXISTS "user";
DROP TABLE IF EXISTS token;
DROP TABLE IF EXISTS notification;

/**
 CREATE TABLE
*/
CREATE TABLE "user"
(
email VARCHAR(255) NOT NULL UNIQUE ,
password TEXT,
PRIMARY KEY (email)
);

CREATE TABLE token
(
token VARCHAR(255) UNIQUE ,
user_email VARCHAR(255),
status BOOL DEFAULT true,
PRIMARY KEY (token)
);

CREATE TABLE notification
(
id UUID NOT NULL UNIQUE ,
message TEXT,
user_email VARCHAR(255),
PRIMARY KEY (id)
);


/**
 CREATE INDEX
 ADD FOREIGN KEY
*/
CREATE INDEX user_email_idx ON "user"(email);
CREATE INDEX token_idx ON token(token);
ALTER TABLE token ADD FOREIGN KEY (user_email) REFERENCES "user" (email);
ALTER TABLE notification ADD FOREIGN KEY (user_email) REFERENCES "user" (email);
