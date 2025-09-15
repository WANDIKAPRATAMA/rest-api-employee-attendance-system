package dto

type SignupRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name" validate:"required"`
}

type SigninRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ChangeRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=employee admin "`
}

type UpdateProfileRequest struct {
	FullName  string `json:"full_name" validate:"omitempty,min=2,max=255"`
	Phone     string `json:"phone" validate:"omitempty,phone"`
	AvatarURL string `json:"avatar_url" validate:"omitempty,url"`
	Address   string `json:"address" validate:"omitempty,max=500"`
}
