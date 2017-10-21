package server

import (
	"text/template"

	"github.com/bluele/gcache"
	"github.com/jinzhu/gorm"

	"github.com/vbogretsov/maild/model"
)

type templateID struct {
	Lang string
	Name string
}

func dbLoader(db *gorm.DB) {
	return func(key interface{}) (interface{}, error) {
		tid := key.(templateID)

		tmt := model.Template{}
		err := db.
			Where("lang = ? and name = ?", tid.Lang, tid.Name).
			First(&tmp).Error

		if err != nil {
			return nil, err
		}

		item, err := template.New(tid.Lang + tid.Name).Parse(tmt.Text)
		if err != nil {
			return nil, err
		}

		return item, nil
	}
}

type templateCache struct {
	cache gcache.Cache
}

func (tc *templateCache) Get(lang, name string) (*template.Template, error) {
	data, err := tc.cache.Get(templateID{Lang: lang, Name: name})
	if err != nil {
		return nil, err
	}

	return data.(*template.Template)
}

func NewDbCache(maxSize int, db *gorm.DB) TemplateCache {
	cache := gcache.New(maxSize).
		LRU().
		LoaderFunc(dbLoader(db)).
		Build()
	return &templateCache{cache: cache}
}