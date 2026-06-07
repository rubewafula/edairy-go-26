package models

type UserStoreAssignment struct {
	BaseModel
	UserID  uint64 `gorm:"column:user_id"`
	StoreID uint64 `gorm:"column:store_id"`
	User    User   `gorm:"foreignKey:UserID"`
	Store   Store  `gorm:"foreignKey:StoreID"`
}

func (UserStoreAssignment) TableName() string {
	return "user_store_assignments"
}
