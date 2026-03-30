package user

// UpdateUserInput 定義可更新的欄位
type UpdateUserInput struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Points *int   `json:"points"`
	Role   string `json:"role"`
}
