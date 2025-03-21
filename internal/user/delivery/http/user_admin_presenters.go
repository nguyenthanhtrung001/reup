package http

import (
	"reup/internal/models"
	"reup/internal/user/usecase"
	"reup/pkg/paginator"
	"reup/pkg/response"
)

type listRequest struct {
	Name string `json:"name"`
}

type listMetaResponse struct {
	paginator.PaginatorResponse
}

type listItemResponse struct {
	ID         string             `json:"id"`
	Email      string             `json:"email"`
	Phone      string             `json:"phone"`
	FullName   string             `json:"fullname"`
	CreatedAt  string             `json:"created_at"`
	Credit     float64            `json:"credit"`
	CreditUsed float64            `json:"credit_used"`
	Discount   map[int]float64    `json:"discount"`
	LastOrder  *response.DateTime `json:"last_order,omitempty"`
	Priority   bool               `json:"priority"`
}

type listResponse struct {
	Users []listItemResponse `json:"users"`
	Meta  listMetaResponse   `json:"meta"`
}

func newListItemResponse(user models.User) listItemResponse {
	items := listItemResponse{
		ID:         user.ID.Hex(),
		Email:      user.Email,
		Phone:      user.Phone,
		FullName:   user.FullName,
		CreatedAt:  user.CreatedAt.Format("2006-01-02 15:04:05"),
		Credit:     user.Credit,
		CreditUsed: user.CreditUsed,
		Discount:   user.Discount,
		Priority:   user.Priority,
	}

	if user.LastOrderAt != nil {
		time := response.DateTime(*user.LastOrderAt)
		items.LastOrder = &time
	}

	return items
}

func newListResponse(lo usecase.ListOutput) listResponse {
	items := make([]listItemResponse, 0, len(lo.Users))
	//indexed order by userid
	for _, v := range lo.Users {
		items = append(items, newListItemResponse(v))
	}

	return listResponse{
		Users: items,
		Meta: listMetaResponse{
			PaginatorResponse: lo.Paginator.ToResponse(),
		},
	}
}

type detailAdminResponse struct {
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
	ApiKey     string          `json:"api_key,omitempty"`
	Priority   bool            `json:"priority"`
}

func newDetailAdminResponse(user models.User) detailAdminResponse {
	profile := detailAdminResponse{
		ID:       user.ID.Hex(),
		Email:    user.Email,
		Phone:    user.Phone,
		FullName: user.FullName,
		Group: groupConfig{
			ID:   user.GroupID.Hex(),
			Name: user.GroupName,
			Role: user.GroupRole,
		},
		Other: otherConfig{
			ClientIP: getClientIP(),
		},
		Credit:     user.Credit,
		CreditUsed: user.CreditUsed,
		Discount:   user.Discount,
		ApiKey:     user.ApiKey,
		Priority:   user.Priority,
	}
	return profile
}

type updateAdminRequest struct {
	FullName string           `json:"fullname"`
	Phone    string           `json:"phone"`
	Email    string           `json:"email"`
	GroupID  string           `json:"group_id"`
	Discount *map[int]float64 `json:"discount"`
	Priority bool             `json:"priority"`
}
