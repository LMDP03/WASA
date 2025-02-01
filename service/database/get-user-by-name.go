package database

var queryGetUserByName = `SELECT id, name FROM Users WHERE name = ?;`

func (db *appdbimpl) GetUserByName(username string) (User, error) {
	var user User
	err := db.c.QueryRow(queryGetUserByName, username).Scan(&user.Id, &user.Name)
	return user, err
}
