package basic

import "errors"

var (
	userList []User = make([]User, 0)
)

type User struct {
	Email        string `json:"e-mail,omitempty"`
	Username     string `json:"user_name,omitempty"`
	Passwordhash string `json:"passwordhash,omitempty"`
	Fullname     string `json:"fullname,omitempty"`
	CreateDate   string `json:"create_date,omitempty"`
	Role         int    `json:"role,omitempty"`
}

func GetUserObject(email string) (User, bool) {
	//needs to be replaces using Database
	for _, user := range userList {
		if user.Email == email {
			return user, true
		}
	}
	return User{}, false
}

// checks if the password hash is valid
func (u *User) ValidatePasswordHash(pswdhash string) bool {
	return u.Passwordhash == pswdhash
}

// this simply adds the user to the list
func AddUserObject(email string, username string, passwordhash string, fullname string, role int) bool {
	// declare the new user object
	newUser := User{
		Email:        email,
		Passwordhash: passwordhash,
		Username:     username,
		Fullname:     fullname,
		Role:         role,
	}
	// check if a user already exists
	for _, ele := range userList {
		if ele.Email == email || ele.Username == username {
			return false
		}
	}
	userList = append(userList, newUser)
	return true
}

// searches the user in the database.
func ValidateUser(email string, passwordHash string) (bool, error) {
	usr, exists := GetUserObject(email)
	if !exists {
		return false, errors.New("user does not exist")
	}
	passwordCheck := usr.ValidatePasswordHash(passwordHash)

	if !passwordCheck {
		return false, nil
	}
	return true, nil
}
