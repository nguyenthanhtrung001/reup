package models

type Scope struct {
	UserID    string `json:"user_id"`
	GroupID   string `json:"group_id"`
	GroupRole string `json:"group_role"`
}

func (s Scope) IsEmpty() bool {
	return s.UserID == ""
}

func (s Scope) IsAdmin() bool {
	return s.GroupRole == "admin"
}

func (s Scope) IsSuperAdmin() bool {
	return s.GroupRole == "superadmin"
}
