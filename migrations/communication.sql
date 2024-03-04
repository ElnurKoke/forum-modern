CREATE TABLE IF NOT EXISTS communication (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	from_user_id INT,
	for_whom_role TEXT NULL,
	about_user_id INT NULL,
	post_id INT NULL,
	comment_id INT NULL,
	message TEXT ,
	created_at DATE DEFAULT (datetime('now','localtime')),
    FOREIGN KEY (from_user_id) REFERENCES user(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS askrole (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	from_user_id INT,
	old_role TEXT,
	new_role TEXT,
	for_whom_role TEXT,
	created_at DATE DEFAULT (datetime('now','localtime')),
    FOREIGN KEY (from_user_id) REFERENCES user(id) ON DELETE CASCADE
);