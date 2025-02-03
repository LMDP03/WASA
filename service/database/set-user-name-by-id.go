package database

var querySetUserNameById = `UPDATE Users SET name = ? WHERE id = ?;`

func (db *appdbimpl) SetUserNameById(username string, userid int) error {

	_, err := db.c.Exec(querySetUserNameById, username, userid)
	if err != nil {
		return err
	}
	return nil
}
