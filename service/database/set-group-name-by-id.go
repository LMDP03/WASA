package database

var querySetGroupNameById = `UPDATE Conversations SET name = ? WHERE id = ?;`

func (db *appdbimpl) SetGroupNameById(convid int, name string) error {

	_, err := db.c.Exec(querySetUserNameById, name, convid)
	if err != nil {
		return err
	}
	return nil
}
