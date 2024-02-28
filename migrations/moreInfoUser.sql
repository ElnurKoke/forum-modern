CREATE TABLE IF NOT EXISTS userInfo (
	user_id INT,
	cash INT DEFAULT 0,
	imageBack TEXT 'black.png',
	imageURL TEXT 'ava.jpg',
	posts INT DEFAULT 0,
    likes INT DEFAULT 0,
    friends INT DEFAULT 0,
	role INT DEFAULT 'user',
	bio TEXT DEFAULT 'hey i am new user',
	created_at DATE DEFAULT (datetime('now','localtime')),
    updated_at DATE DEFAULT (datetime('now','localtime')),
	FOREIGN KEY (user_id) REFERENCES user(id)
);