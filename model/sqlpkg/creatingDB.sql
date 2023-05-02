CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			name TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			dateCreate TIMESTAMP NOT NULL,
			session TEXT,
			expirySession TIMESTAMP
		);
		
		CREATE TABLE 'posts_likes' (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			userID INT NOT NULL REFERENCES users(id),
			messageID INT NOT NULL REFERENCES posts(id),
			like BOOL NOT NULL,
			UNIQUE (userID, messageID)
		);
		
		CREATE TABLE 'comments_likes' (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			userID INT NOT NULL REFERENCES users(id),
			messageID INT NOT NULL REFERENCES comments(id),
			like BOOL NOT NULL,
			UNIQUE (userID, messageID)
		);

		CREATE TABLE 'posts' (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			theme TEXT NOT NULL DEFAULT ('(No theme)'),
			content TEXT NOT NULL, 
			authorID INT NOT NULL REFERENCES users(id),
			dateCreate TIMESTAMP NOT NULL
		);

		CREATE TABLE 'comments' (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			content TEXT NOT NULL, 
			authorID INT NOT NULL REFERENCES users(id),
			dateCreate TIMESTAMP NOT NULL,
			postID INT NOT NULL REFERENCES posts(id)
		);
		
		CREATE TABLE 'categories' (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			name TEXT NOT NULL 
		);
		
		CREATE TABLE 'post_categories' (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			categoryID INT NOT NULL REFERENCES categories(id), 
			postID INT NOT NULL REFERENCES posts(id),
			UNIQUE (categoryID, postID)
		);

		CREATE INDEX userssession ON users (session);

		INSERT INTO users (name,email,password, dateCreate) VALUES (?,?,?,?);
		INSERT INTO categories (name) VALUES (?), (?), (?), (?);
		
        SELECT auth_user_add('webuser', 'webuser', 0);