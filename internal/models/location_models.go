package models

type County struct {
	BaseModel
	Name string `gorm:"uniqueIndex;column:name"`
	Code string `gorm:"uniqueIndex;column:code"`
}

func (County) TableName() string {
	return "counties"
}

type SubCounty struct {
	BaseModel
	CountyID uint64 `gorm:"column:county_id"`
	Name     string `gorm:"column:name"`
}

func (SubCounty) TableName() string {
	return "sub_counties"
}

type Ward struct {
	BaseModel
	SubCountyID uint64 `gorm:"column:sub_county_id"`
	Name        string `gorm:"column:name"`
}

func (Ward) TableName() string {
	return "wards"
}
