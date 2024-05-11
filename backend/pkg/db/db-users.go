package db

type UserData struct {
	UserID        int    `json:"user_id"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	Firstname     string `json:"firstname"`
	Lastname      string `json:"lastname"`
	DateOfBirth   string `json:"date_of_birth"`
	Avatar        string `json:"avatar"`
	Nickname      string `json:"nickname"`
	AboutMe       string `json:"about_me"`
	ProfilePublic bool   `json:"profile_public"`
	CreatedAt     string `json:"created_at"`
}

func GetUsersFromDB() ([]UserData, error) {
	var users []UserData

	// Execute the SQL query to fetch all users' data
	rows, err := DB.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the result set and populate the users slice
	for rows.Next() {
		var user UserData
		if err := rows.Scan(&user.UserID, &user.Email, &user.Password, &user.Firstname, &user.Lastname, &user.DateOfBirth, &user.Avatar, &user.Nickname, &user.AboutMe, &user.ProfilePublic, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
