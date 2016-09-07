-- SQL SEEDS

DELETE FROM notification;
DELETE FROM token;
DELETE FROM "user";

INSERT INTO "user" (id, email, password) VALUES ('c93791f0-29d8-4e29-b944-ec2bc5692069', 'foo@example.com', 'password');
INSERT INTO token (id, token, user_id) VALUES ('c14b368b-1467-461b-85bf-3feda8339f8e', 'atesttoken', 'c93791f0-29d8-4e29-b944-ec2bc5692069');
INSERT INTO notification (id, message, user_id) VALUES ('409c4e13-10e5-4530-b35b-b57c01d01e9c', 'A Notification', 'c93791f0-29d8-4e29-b944-ec2bc5692069');
