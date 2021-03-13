package fruits

// List retrieves all fruits
func List() []Fruits {
	list, err := model.List()
	if err != nil {
		return nil
	}
	return list
}

// Get retrieves an fruit by ID
func Get(id int64) Fruits {
	f, err := model.Get(id)
	if err != nil {
		return Fruits{}
	}
	return f
}

// Create add a new fruit
func Create(user Fruits) int64 {
	id, err := model.Create(user)
	if err != nil {
		return -1
	}
	return id
}

// Delete removes a fruit by ID
func Delete(id int64) {
	err := model.Delete(id)
	if err != nil {
		// do somehting here
	}
}

// Update changes the data of a fruit
func Update(user Fruits) Fruits {
	f, err := model.Update(user)
	if err != nil {
		return Fruits{}
	}
	return f
}
