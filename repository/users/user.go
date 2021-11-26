package users

type User struct {
	Id     string   `json:"id" bson:"_id"`
	Name   string   `json:"name" bson:"name,omitempty"` // omitempty คือถ้าไม่มีค่าก็ไม่ต้องใส่ key เข้าไปใน db
	Active []string `json:"active" bson:"active"`
	Equity float64  `json:"equity" bson:"equity,omitempty"`
}

type UserRepository interface {
	Create(User) (string, error)
	GetAll() ([]User, error)
	GetUserById(string) (*User, error)
	UpdateUserById(string, User) (*User, error)
	DeleteUserById(string) (int, error)
}
