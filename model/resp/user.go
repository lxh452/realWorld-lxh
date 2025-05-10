package resp

type UserResp struct {
	User UserModel `json:"user"`
}
type UserModel struct {
	Id       uint    `json:"-"`
	Email    string  `json:"email"`
	Token    string  `json:"token"`
	Username string  `json:"username"`
	Bio      string  `json:"bio"`
	Image    *string `json:"image"`
}

func (receiver UserResp) TableName() string {
	return "users"
}

func (receiver UserModel) TableName() string {
	return "users"
}
