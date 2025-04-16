package models

import (
	"Yakudza/pkg/database"
)

type Links struct {
	ID       uint   `json:"id"`
	Title    string `gorm:"not null" json:"title"`
	Link     string `gorm:"not null" json:"link"`
	Position uint   `gorm:"default:0" json:"position"`
}

func AllLinks() ([]*Links, error) {
	db := database.GetDB()
	var links []*Links

	err := db.Find(&links).Error

	if err != nil {
		return nil, err
	}

	return links, nil
}

func (l *Links) Create() error {
	db := database.GetDB()

	return db.Create(&l).Error
}

func (l *Links) FindID() error {
	db := database.GetDB()

	return db.Where(&Links{ID: l.ID}).Find(&l).Error
}

func (l *Links) Delete() error {
	db := database.GetDB()

	if l.ID == 10 {
		return nil
	}

	return db.Delete(&l).Error
}

func (l *Links) DeleteByID() error {
	dbase := database.GetDB()
	if l.ID == 10 {
		return nil
	}
	result := dbase.Delete(&Links{ID: l.ID})

	return result.Error
}

func (l *Links) Update() error {
	dbase := database.GetDB()
	if l.ID == 10 {
		return nil
	}
	result := dbase.Model(&l).Updates(&l)
	return result.Error
}
