package user

// User represents the users table
type User struct {
	ID        uint
	FirstName string
	LastName  string
	Email     string
	Password  string
	Phone     string
}
