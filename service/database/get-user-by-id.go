package database

var queryGetUserById = `SELECT id, name FROM Users WHERE id = ?;`

func (db *appdbimpl) GetUserById(userid int) (User, error) {
	var user User
	err := db.c.QueryRow(queryGetUserById, userid).Scan(&user.Id, &user.Name)
	return user, err
}
