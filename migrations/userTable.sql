CREATE TABLE IF NOT EXISTS user (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	email TEXT UNIQUE,
	username TEXT UNIQUE,
	password TEXT,
	session_token TEXT DEFAULT NULL,
	expiresAt DATETIME DEFAULT NULL,
	imageBack TEXT DEFAULT 'black.png',
	imageURL TEXT DEFAULT 'ava.jpg',
	rol INT DEFAULT 'user',
	bio TEXT DEFAULT 'hey i am new user',
	created_at DATE DEFAULT (datetime('now','localtime')),
    updated_at DATE DEFAULT (datetime('now','localtime'))
);


-- DROP TABLE IF EXISTS likesComment;
-- DROP TABLE IF EXISTS likesPost;
-- DROP TABLE IF EXISTS comment;
-- DROP TABLE IF EXISTS post;