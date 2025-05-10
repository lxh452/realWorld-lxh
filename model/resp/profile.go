package resp

type ProfileResp struct {
	Profile ProfileModel `json:"profile"`
}
type ProfileModel struct {
	Username  string  `json:"username"`
	Bio       string  `json:"bio"`
	Image     *string `json:"image"`
	Following bool    `json:"following"` //关注状态
}
