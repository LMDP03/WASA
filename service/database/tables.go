package database

var sql_USERS = `CREATE TABLE IF NOT EXISTS Users
(
	id INTEGER NOT NULL AUTO_INCREMENT,
	name STRING NOT NULL UNIQUE,
	PRIMARY KEY(id)
)
`

var sql_CONVERSATIONS = `CREATE TABLE IF NOT EXISTS Conversations
(
	id INTEGER NOT NULL AUTO_INCREMENT,
	name STRING NOT NULL,
	group BOOLEAN,
	last_message INTEGER,
	PRIMARY KEY (id)
	CONSTRAINT fk_Conversations
		FOREIGN KEY (last_message) REFERENCES Messages(msgId)
			ON DELETE CASCADE
)`

var sql_PARTICIPANTS = `CREATE TABLE IF NOT EXISTS Participants
(
	convId INTEGER NOT NULL,
	userId INTEGER NOT NULL,
	PRIMARY KEY (convId, userId),
	CONSTRAINT fk_Participants
		FOREIGN KEY (convId) REFERENCES Conversation(id)
			ON DELETE CASCADE
		FOREIGN KEY (userId) REFERNCES User(id)
			ON DELETE CASCADE
)`

var sql_MESSAGES = `CREATE TABLE IF NOT EXISTS Messages
(
	convId INTEGER NOT NULL,
	senderId INTEGER NOT NULL,
	msgId INTEGER NOT NULL AUTO_INCREMENT,
	text TEXT,
	image STRING,
	timeStamp DATETIME DEFAULT CURRENT_TIMESTAMP,
	responseTo INTEGER NOT NULL.const,
	checkMark ENUM("sent", "received", "read"),
	PRIMARY KEY (convId, msgId),
	CONSTRAINT fk_Messages
		FOREIGN KEY (ConvId) REFERENCES Conversations(id)
			ON DELETE CASCADE
		FOREIGN KEY (senderId) REFERNCES User(id)
			ON DELETE CASCADE
)`

var sql_REACTIONS = `CREATE TABLE IF NOT EXISTS Reactions
(
	convId INTEGER NOT NULL,
	senderId INTEGER NOT NULL,
	msgId INTEGER NOT NULL,
	emoji STRING,
	PRIMARY KEY (convId, msgId, senderId),
	CONSTRAINT fk_Reactions
		FOREIGN KEY (convId, msgId) REFERENCES Messages(convId, msgId)
			ON DELETE CASCADE
    	FOREIGN KEY (SenderId) REFERENCES Users(id)
      		ON DELETE CASCADE
	) `
