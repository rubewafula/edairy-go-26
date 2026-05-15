package models

type CustomerType struct {
	BaseModel
	Name        string `gorm:"type:varchar(125);uniqueIndex;not null"`
	Description string `gorm:"type:varchar(255)"`
}

func (CustomerType) TableName() string {
	return "customer_types"
}
