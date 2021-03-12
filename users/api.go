package users

// List retrieves all users
func List() []Users {
	list, err := model.List()
	if err != nil {
		return nil
	}
	return list
}

// Get retrieves an user by login
func Get(login string) Users {
	u, err := model.Get(login)
	if err != nil {
		return Users{}
	}
	return u
}

// Create add a new user
func Create(user Users) {
	err := model.Create(user)
	if err != nil {
		// do something later
	}
}

// Delete removes an user by login
func Delete(login string) {
	err := model.Delete(login)
	if err != nil {
		// do somehting later
	}
}

// Update changes the data of an user
func Update(user Users) Users {
	user, err := model.Update(user)
	if err != nil {
		return Users{}
	}
	return user
}
