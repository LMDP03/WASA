package database

var queryFindUserById = `SELECT id, name FROM Users WHERE id = ?;`

func (db *appdbimpl) GetUserById(userid int32) (User, error) {
	var user User
	err := db.c.QueryRow(queryFindUserByName, userid).Scan(&user.Id, &user.Name)
	return user, err
}
