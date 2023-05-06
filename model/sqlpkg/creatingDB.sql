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
			userID INT NOT NULL,
			messageID INT NOT NULL,
			like BOOL NOT NULL,
			UNIQUE (userID, messageID),
			FOREIGN KEY (userID) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (messageID) REFERENCES posts(id) ON DELETE CASCADE
		);
		
		CREATE TABLE 'comments_likes' (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			userID INT NOT NULL,
			messageID INT NOT NULL,
			like BOOL NOT NULL,
			UNIQUE (userID, messageID),
			FOREIGN KEY (userID) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (messageID) REFERENCES posts(id) ON DELETE CASCADE
		);

		CREATE TABLE 'posts' (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			theme TEXT NOT NULL DEFAULT ('(No theme)'),
			content TEXT NOT NULL, 
			authorID INT NOT NULL,
			dateCreate TIMESTAMP NOT NULL,
			FOREIGN KEY (authorID) REFERENCES users(id) ON DELETE CASCADE
		);

		CREATE TABLE 'comments' (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			content TEXT NOT NULL, 
			authorID INT NOT NULL,
			dateCreate TIMESTAMP NOT NULL,
			postID INT NOT NULL,
			FOREIGN KEY (authorID) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (postID) REFERENCES posts(id) ON DELETE CASCADE
		);
		
		CREATE TABLE 'categories' (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			name TEXT NOT NULL 
		);
		
		CREATE TABLE 'post_categories' (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			categoryID INT NOT NULL, 
			postID INT NOT NULL,
			UNIQUE (categoryID, postID),
			FOREIGN KEY (categoryID) REFERENCES categories(id) ON DELETE CASCADE,
			FOREIGN KEY (postID) REFERENCES posts(id) ON DELETE CASCADE
		);

		CREATE INDEX userssession ON users (session);

		INSERT INTO users (name,email,password, dateCreate) VALUES (?,?,?,?);
		INSERT INTO categories (name) VALUES (?), (?), (?), (?);
		
        SELECT auth_user_add('webuser', 'webuser', 0);