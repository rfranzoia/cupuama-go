package products

// List retrieves all products
func List() []Products {
	list, err := model.List()
	if err != nil {
		return nil
	}
	return list
}

// Get retrieves an product by ID
func Get(id int64) Products {
	f, err := model.Get(id)
	if err != nil {
		return Products{}
	}
	return f
}

// Create add a new product
func Create(user Products) int64 {
	id, err := model.Create(user)
	if err != nil {
		return -1
	}
	return id
}

// Delete removes an product by ID
func Delete(id int64) {
	err := model.Delete(id)
	if err != nil {
		// do somehting here
	}
}

// Update changes the data of a product
func Update(user Products) Products {
	f, err := model.Update(user)
	if err != nil {
		return Products{}
	}
	return f
}
