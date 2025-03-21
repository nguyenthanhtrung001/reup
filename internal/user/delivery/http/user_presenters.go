package http

import (
	"regexp"

	"reup/internal/user/usecase"
	"reup/pkg/response"
)

type registerRequest struct {
	FullName  string `form:"fullname" json:"fullname" binding:"required"`
	Email     string `form:"email" json:"email" binding:"required"`
	Phone     string `form:"phone" json:"phone" binding:"required"`
	Password  string `form:"password" json:"password" binding:"required"`
	RefererID string `form:"referer_id" json:"referer_id"`
}

type loginRequest struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type updateProfileRequest struct {
	FullName string `form:"fullname" json:"fullname" binding:"required"`
}

type updatePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type profileRequest struct{}

type loginResponse struct {
	Token string          `json:"token"`
	User  profileResponse `json:"user"`
}
type groupConfig struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}
type otherConfig struct {
	ClientIP string `json:"client_ip"`
}
type refererInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type profileResponse struct {
	ID         string          `json:"id"`
	Email      string          `json:"email"`
	Phone      string          `json:"phone"`
	FullName   string          `json:"fullname"`
	Group      groupConfig     `json:"group"`
	Other      otherConfig     `json:"other"`
	Referer    *refererInfo    `json:"referer,omitempty"`
	Credit     float64         `json:"credit"`
	CreditUsed float64         `json:"credit_used"`
	Discount   map[int]float64 `json:"discount"`
}

func newProfileResponse(uo usecase.ProfileOutput) profileResponse {
	profile := profileResponse{
		ID:       uo.ID,
		Email:    uo.Email,
		Phone:    uo.Phone,
		FullName: uo.FullName,
		Group: groupConfig{
			ID:   uo.Group.ID,
			Name: uo.Group.Name,
			Role: uo.Group.Role,
		},
		Other: otherConfig{
			ClientIP: getClientIP(),
		},
		Credit:     uo.Credit,
		CreditUsed: uo.CreditUsed,
		Discount:   uo.Discount,
	}

	if uo.Referer.ID != "" {
		profile.Referer = &refererInfo{
			ID:   uo.Referer.ID,
			Name: uo.Referer.Name,
		}
	}
	return profile
}

func newLoginResponse(lo usecase.LoginOutput) loginResponse {
	return loginResponse{
		Token: lo.Token,
		User:  newProfileResponse(lo.User),
	}
}

func (reg registerRequest) validateEmail() bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	match, err := regexp.MatchString(emailRegex, reg.Email)
	if err != nil {
		return false
	}
	return match
}

func (reg registerRequest) validatePhone() bool {
	// Regular expression to match Vietnamese phone numbers
	var re = regexp.MustCompile(`^(03|05|07|08|09|024|028)[0-9]{7,8}$`)
	return re.MatchString(reg.Phone)
}

func (reg registerRequest) validatePassword() bool {
	return validPassword(reg.Password)
}

func validPassword(pwd string) bool {
	return len(pwd) >= 8
}

type adminItemResponse struct {
	ID          string `json:"id"`
	ServiceType string `json:"service_type"`

	ServiceViewType string `json:"service_view_type,omitempty"` // search, suggest, external || just for service_type = view
	MaxViewTime     *int   `json:"max_view_time,omitempty"`     // just for service_type = view
	MinViewTime     *int   `json:"min_view_time,omitempty"`     // just for service_type = view

	Platform       string            `json:"platform"`
	Category       string            `json:"category"`
	Type           string            `json:"type"`
	ViewType       string            `json:"view_type"`
	ServiceID      int               `json:"service_id"`
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	PricePer10     float64           `json:"price_per_10"`
	MaxThreads     int               `json:"max_threads,omitempty"`
	MaxThreads3000 int               `json:"max_threads_3000,omitempty"`
	MaxThreads5000 int               `json:"max_threads_5000,omitempty"`
	Priority       bool              `json:"priority"`
	Enabled        bool              `json:"enabled"`
	Min            int               `json:"min,omitempty"`
	Max            int               `json:"max,omitempty"`
	Geo            string            `json:"geo"`
	RestApi        bool              `json:"rest_api"`
	CreatedAt      response.DateTime `json:"created_at"`
	UpdatedAt      response.DateTime `json:"updated_at"`
}
