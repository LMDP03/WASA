package database

var queryGetUsersByName = `SELECT id, name FROM Users WHERE name LIKE '%?%' AND NOT id = ? OREDR BY name;`

func (db *appdbimpl) GetUsersByName(searchname string, userid int) ([]User, error) {

	var users []User

	rows, err := db.c.Query(queryGetUsersByName, searchname, userid)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if rows.Err() != nil {
			return nil, err
		}
		var user User
		if err := rows.Scan(&user.Id, &user.Name); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	defer func() { err = rows.Close() }()
	return users, err
}
