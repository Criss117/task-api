package users

type UsersRepository struct {
	Users []*User
}

func NewUsersRepository() *UsersRepository {

	return &UsersRepository{
		Users: []*User{
			NewUser("Cristian", "cristian@gmail.com", "holamundo"),
		},
	}
}

func (r *UsersRepository) AddUser(user *User) {
	r.Users = append(r.Users, user)
}

func (r *UsersRepository) GetUserByEmail(email string) *User {
	for _, user := range r.Users {
		if user.Email == email {
			return user
		}
	}
	return nil
}

func (r *UsersRepository) GetUserByID(id string) *User {
	for _, user := range r.Users {
		if user.ID == id {
			return user
		}
	}
	return nil
}
