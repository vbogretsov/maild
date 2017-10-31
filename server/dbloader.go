package server

import (
	"text/template"

	"github.com/jinzhu/gorm"

	"github.com/vbogretsov/maild/model"
)

type templateData struct {
	Lang  string
	Name  string
	Value string
}

func (templateData) TableName() string {
	return "templates"
}

func NewDbLoader(db *gorm.DB) model.TemplateLoader {
	return func(key *model.TemplateID) (*template.Template, error) {
		tml := templateData{}
		err := db.
			Where("lang = ? and name = ?", key.Lang, key.Name).
			First(&tml).Error

		if err != nil {
			return nil, err
		}

		item, err := template.New(key.Lang + key.Name).Parse(tml.Value)
		if err != nil {
			return nil, err
		}

		return item, nil
	}
}
