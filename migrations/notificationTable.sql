DROP TABLE IF EXISTS notification;

CREATE TABLE IF NOT EXISTS notification (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INT,
	comment_id INT,
    author TEXT,
    reactauthor TEXT,
    message TEXT,
    activity INT DEFAULT 1,
    created_at DATE,
    FOREIGN KEY (post_id) REFERENCES post(id) ON DELETE CASCADE
);
