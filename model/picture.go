package model

import (
	"github.com/easykoo/binding"

	. "github.com/easykoo/go-display/common"

	"errors"
	"net/http"
	"time"
)

var RelateErr = errors.New("此图片正在使用中!")

type Picture struct {
	Id          int       `form:"pictureId" xorm:"int(11) pk not null autoincr"`
	Name        string    `form:"name" xorm:"varchar(20) not null"`
	Url         string    `form:"url" xorm:"varchar(80) not null"`
	Description string    `form:"description" xorm:"varchar(100) not null"`
	CreateUser  string    `xorm:"varchar(20) default 'SYSTEM'"`
	CreateDate  time.Time `xorm:"datetime created"`
	UpdateUser  string    `xorm:"varchar(20) default 'SYSTEM'"`
	UpdateDate  time.Time `xorm:"datetime updated"`
	Version     int       `form:"version" xorm:"int(11) version"`
	Page        `xorm:"-"`
}

func (self *Picture) Exist() (bool, error) {
	return orm.Get(self)
}

func (self *Picture) Get() (*Picture, error) {
	picture := &Picture{}
	_, err := orm.Id(self.Id).Get(picture)
	return picture, err
}

func (self *Picture) Insert() error {
	_, err := orm.InsertOne(self)
	Log.Info(self.Name, " inserted")
	return err
}

func (self *Picture) Update() error {
	_, err := orm.Id(self.Id).Update(self)
	Log.Info("Picture ", self.Name, " updated")
	return err
}

func (self *Picture) Delete() error {
	if exist, err := orm.Get(&PadPicture{PictureId: self.Id}); exist {
		if err != nil {
			return err
		}
		return RelateErr
	}
	_, err := orm.Delete(self)
	Log.Info("Picture ", self.Name, " deleted")
	return err
}

func (self *Picture) DeletePictures(array []int) error {
	_, err := orm.In("id", array).Where("id not in (select picture_id from pad_picture)").Delete(&Picture{})
	Log.Info("Pictures: ", array, " deleted")
	return err
}

func (self *Picture) SelectAll() ([]Picture, error) {
	var pictures []Picture
	err := orm.Find(&pictures)
	return pictures, err
}

func (self *Picture) SearchByPage() ([]Picture, int64, error) {
	total, err := orm.Count(self)
	var pictures []Picture
	err = orm.OrderBy(self.GetSortProperties()[0].Column+" "+self.GetSortProperties()[0].Direction).Limit(self.GetPageSize(), self.GetDisplayStart()).Find(&pictures, self)
	return pictures, total, err
}

func (self *Picture) ArrayInfo(array string) ([]Picture, error) {
	var pictures []Picture
	err := orm.Where("id in ("+array+")").Find(&pictures, Picture{})
	return pictures, err
}

func (picture Picture) Validate(errors *binding.Errors, r *http.Request) {
	if len(picture.Description) > 20 {
		errors.Fields["name"] = "Length of name should be less than 20."
	}
}
