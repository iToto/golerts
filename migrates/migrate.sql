/* SQLEditor (Postgres)*/

/**
 DROP TABLE IF EXISTS
*/

DROP TABLE IF EXISTS token;
DROP TABLE IF EXISTS notification;
DROP TABLE IF EXISTS "user";

/**
 CREATE TABLE
*/
CREATE TABLE "user"
(
id UUID NOT NULL UNIQUE ,
email VARCHAR(255) NOT NULL UNIQUE ,
password TEXT,
PRIMARY KEY (id)
);

CREATE TABLE token
(
id UUID NOT NULL UNIQUE ,
token VARCHAR(255) UNIQUE ,
user_id UUID,
status BOOL DEFAULT TRUE,
PRIMARY KEY (id)
);

CREATE TABLE notification
(
id UUID NOT NULL UNIQUE ,
message TEXT,
user_id UUID,
PRIMARY KEY (id)
);

/**
 CREATE INDEX
 ADD FOREIGN KEY
*/
CREATE INDEX user_id_idx ON "user"(id);

CREATE INDEX user_email_idx ON "user"(email);

CREATE INDEX token_idx ON token(token);

ALTER TABLE token ADD FOREIGN KEY (user_id) REFERENCES "user" (id);

ALTER TABLE notification ADD FOREIGN KEY (user_id) REFERENCES "user" (id);