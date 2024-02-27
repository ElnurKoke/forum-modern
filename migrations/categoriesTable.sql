CREATE TABLE IF NOT EXISTS hashtags (
    hashtag_id INTEGER PRIMARY KEY AUTOINCREMENT,
    hashtag TEXT
);

DELETE FROM hashtags;

INSERT INTO hashtags (hashtag) VALUES ('Art');
INSERT INTO hashtags (hashtag) VALUES ('Animal');
INSERT INTO hashtags (hashtag) VALUES ('Anime');

INSERT INTO hashtags (hashtag) VALUES ('Book');
INSERT INTO hashtags (hashtag) VALUES ('Cars');
INSERT INTO hashtags (hashtag) VALUES ('Education');

INSERT INTO hashtags (hashtag) VALUES ('Food');
INSERT INTO hashtags (hashtag) VALUES ('Game');
INSERT INTO hashtags (hashtag) VALUES ('Legend');

INSERT INTO hashtags (hashtag) VALUES ('Marvel');
INSERT INTO hashtags (hashtag) VALUES ('Medicine');
INSERT INTO hashtags (hashtag) VALUES ('Movie');

INSERT INTO hashtags (hashtag) VALUES ('Psychology');
INSERT INTO hashtags (hashtag) VALUES ('Nature');
INSERT INTO hashtags (hashtag) VALUES ('News');

INSERT INTO hashtags (hashtag) VALUES ('Technology');
INSERT INTO hashtags (hashtag) VALUES ('Sport');
INSERT INTO hashtags (hashtag) VALUES ('Other');