		INSERT INTO users (name,email,password, dateCreate) VALUES ("test1","test1@forum", ? ,"2023-03-20 09:41:04.656479916+00:00");
		INSERT INTO users (name,email,password, dateCreate) VALUES ("test2","test2@forum", ? ,"2023-03-20 09:52:07.656479916+00:00");
		INSERT INTO posts (theme,content,authorID, dateCreate) VALUES ("cats", "cats are cute", 1, "2023-03-20 15:41:42.656479916+00:00");
		INSERT INTO posts (theme,content,authorID, dateCreate) VALUES ("dogs", "dogs are funny", 2, "2023-03-21 14:41:04.656479916+00:00");
		INSERT INTO posts (theme,content,authorID, dateCreate) VALUES ("My cat", "She is the best", 3, "2023-03-22 10:41:23.656479916+00:00");
		INSERT INTO posts (theme,content,authorID, dateCreate) VALUES ("My dog"," He is the best", 2, "2023-03-22 11:41:14.656479916+00:00");
		INSERT INTO posts (theme,content,authorID, dateCreate) VALUES ("My parrot", "My parrot is a pirate", 1, "2023-03-23 11:41:52.656479916+00:00");
		INSERT INTO posts (theme,content,authorID, dateCreate) VALUES ("Seamus", "My dog is such a cheeky monkey", 1, "2023-03-24 14:41:00.656479916+00:00");
		INSERT INTO posts (theme,content,authorID, dateCreate) VALUES ("Wise Kaa", "How many monkeys can a python swallow in one gulp?", 3, "2023-03-24 21:25:25.656479916+00:00");
		
		INSERT INTO comments (content,authorID, dateCreate,postID) VALUES ("No, mine", 1, "2023-03-22 11:21:05.656479916+00:00",3);
		INSERT INTO comments (content,authorID, dateCreate,postID) VALUES ("25, ish", 2, "2023-03-27 09:43:13.656479916+00:00",7);
		
		INSERT INTO posts_likes (userID, messageID, like) VALUES (2, 4, 0);
		INSERT INTO posts_likes (userID, messageID, like) VALUES (2, 3, 0);
		INSERT INTO posts_likes (userID, messageID, like) VALUES (2, 2, 1);
		INSERT INTO posts_likes (userID, messageID, like) VALUES (1, 2, 1);
		INSERT INTO posts_likes (userID, messageID, like) VALUES (1, 3, 1);
		INSERT INTO posts_likes (userID, messageID, like) VALUES (1, 1, 1);
		
		INSERT INTO comments_likes (userID, messageID, like) VALUES (1, 1, 1);
		INSERT INTO comments_likes (userID, messageID, like) VALUES (3, 1, 0);
		
		INSERT INTO post_categories (categoryID, postID) VALUES (1,1);
		INSERT INTO post_categories (categoryID, postID) VALUES (3,1);
		INSERT INTO post_categories (categoryID, postID) VALUES (2,2);
		INSERT INTO post_categories (categoryID, postID) VALUES (3,2);
		INSERT INTO post_categories (categoryID, postID) VALUES (1,3);
		INSERT INTO post_categories (categoryID, postID) VALUES (3,3);
		INSERT INTO post_categories (categoryID, postID) VALUES (2,4);
		INSERT INTO post_categories (categoryID, postID) VALUES (3,4);
		INSERT INTO post_categories (categoryID, postID) VALUES (3,5);
		INSERT INTO post_categories (categoryID, postID) VALUES (2,6);
		INSERT INTO post_categories (categoryID, postID) VALUES (4,6);
		INSERT INTO post_categories (categoryID, postID) VALUES (4,7);
		

