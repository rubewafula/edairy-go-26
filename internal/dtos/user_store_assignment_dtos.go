package dtos

type CreateUserStoreAssignmentRequest struct {
	UserID  uint64 `json:"user_id" validate:"required"`
	StoreID uint64 `json:"store_id" validate:"required"`
}

type UpdateUserStoreAssignmentRequest struct {
	UserID  uint64 `json:"user_id" validate:"required"`
	StoreID uint64 `json:"store_id" validate:"required"`
}

type UserStoreAssignmentResponse struct {
	ID        uint64 `json:"id"`
	UserID    uint64 `json:"user_id"`
	UserName  string `json:"user_name"`
	StoreID   uint64 `json:"store_id"`
	StoreName string `json:"store_name"`
}
