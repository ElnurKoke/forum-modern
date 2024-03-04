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

INSERT INTO post(title, description, imageURL, author_id, category)
SELECT 'Welcome my new user', 
'I am not an admin, 
but I am a king 
login: great_king 
password: 123qweASD!@#', 
'test.png', 1, 'king message'
WHERE NOT EXISTS (
    SELECT * FROM post WHERE category = 'king message'
);

INSERT INTO post(title, description, imageURL, author_id, category)
SELECT 'Hello new user', 
'Admin, 
login: great_queen 
password: 123qweASD!@#',
'test.png', 1, 'king message 2'
WHERE NOT EXISTS (
    SELECT * FROM post WHERE category = 'king message 2'
);

INSERT INTO post(title, description, imageURL, author_id, category)
SELECT 'Hi new user', 
'Moderator, 
login: prince 
password: 123qweASD!@#',
'test.png', 1, 'king message 3'
WHERE NOT EXISTS (
    SELECT * FROM post WHERE category = 'king message 3'
);

UPDATE post SET status = 'done' WHERE id= 1 OR id=2 OR id=3;

-- DROP TABLE IF EXISTS likesComment;
-- DROP TABLE IF EXISTS likesPost;
-- DROP TABLE IF EXISTS comment;
-- DROP TABLE IF EXISTS post;