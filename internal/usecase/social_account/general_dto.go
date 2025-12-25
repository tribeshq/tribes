package social_account

type SocialAccountOutputDTO struct {
	Id        uint   `json:"id"`
	UserId    uint   `json:"user_id"`
	Username  string `json:"username"`
	Platform  string `json:"platform"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
