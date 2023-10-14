package structure

type User struct {
	Id       int64
	Name     string
	Slug     string
	Email    string
	Image    string
	Cover    string
	Bio      string
	Website  string
	Location string
	Role     int64 //1 = Administrator, 2 = Editor, 3 = Author, 4 = Owner
}
