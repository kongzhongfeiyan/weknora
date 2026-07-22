package types

import "time"

type WitCompany struct {
	ID            int64     `json:"id"             gorm:"primaryKey;autoIncrement"`
	CompanyName   string    `json:"company_name"   gorm:"type:varchar(255);not null"`
	CompanyCode   string    `json:"company_code"   gorm:"type:varchar(100);not null;uniqueIndex"`
	Address       string    `json:"address"        gorm:"type:text"`
	ContactPerson string    `json:"contact_person" gorm:"type:varchar(100)"`
	CreatedAt     time.Time `json:"created_at"     gorm:"autoCreateTime"`
	CreatedBy     string    `json:"created_by"     gorm:"type:varchar(64)"`
}

func (WitCompany) TableName() string { return "witcompany" }
