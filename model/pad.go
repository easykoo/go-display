package model

import (
	"github.com/easykoo/binding"

	. "github.com/easykoo/go-display/common"

	"net/http"
	"time"
)

type Pad struct {
	Id         int       `form:"padId" xorm:"int(11) pk not null autoincr"`
	Name       string    `form:"name" xorm:"varchar(20) not null"`
	PictureStr string    `form:"pictureStr" xorm:"-"`
	Pictures   []Picture `xorm:"-"`
	//PictureArray []Picture `json:"pictureArray" xorm:"-"`
	Interval    int       `form:"interval", xorm:"int(1) default 5"`
	Color       string    `form:"color" xorm:"varchar(20) not null"`
	Description string    `form:"description" xorm:"varchar(100) not null"`
	CreateUser  string    `xorm:"varchar(20) default 'SYSTEM'"`
	CreateDate  time.Time `xorm:"datetime created"`
	UpdateUser  string    `xorm:"varchar(20) default 'SYSTEM'"`
	UpdateDate  time.Time `xorm:"datetime updated"`
	Version     int       `form:"version" xorm:"int(11) version"`
	Page        `xorm:"-"`
}

type PadPicture struct {
	PadId     int `xorm:"int(11) not null "`
	PictureId int `xorm:"int(11) not null "`
}

func (self *Pad) Exist() (bool, error) {
	return orm.Get(self)
}

func (self *Pad) GetById() (*Pad, error) {
	pad := &Pad{}
	_, err := orm.Id(self.Id).Get(pad)
	pad.LoadPictures()
	return pad, err
}
func (self *Pad) Get() (*Pad, error) {
	exist, err := orm.Get(self)
	if !exist {
		return nil, err
	}
	self.LoadPictures()
	return self, err
}

func (self *Pad) Insert() error {
	_, err := orm.InsertOne(self)
	Log.Info(self.Name, " inserted")
	return err
}

func (self *Pad) Update() error {
	_, err := orm.Id(self.Id).Update(self)
	Log.Info("Pad ", self.Name, " updated")
	return err
}

func (self *Pad) UpdatePictures() error {
	session := orm.NewSession()
	defer session.Close()
	_, err := session.Exec("delete from pad_picture where pad_id = ?", self.Id)
	for key, _ := range self.Pictures {
		session.Insert(PadPicture{PadId: self.Id, PictureId: self.Pictures[key].Id})
	}
	Log.Info("Pad", self.Id, " pictures updated!")
	return err
}

func (self *Pad) Delete() error {
	session := orm.NewSession()
	defer session.Close()
	_, err := session.Delete(self)
	_, err = session.Exec("delete from pad_picture where pad_id = ?", self.Id)
	Log.Info("Pad ", self.Name, " deleted")
	return err
}

func (self *Pad) DeletePads(array []int) error {
	session := orm.NewSession()
	defer session.Close()
	_, err := session.In("id", array).Delete(&Pad{})
	_, err = session.In("pad_id", array).Delete(&PadPicture{})
	Log.Info("Pads: ", array, " deleted")
	return err
}

func (self *Pad) SelectAll() ([]Pad, error) {
	var pads []Pad
	err := orm.Find(&pads)
	return pads, err
}

func BatchLoadPictures(pads []Pad) {
	for key, _ := range pads {
		pads[key].LoadPictures()
	}
}

func (self *Pad) LoadPictures() {
	var pictures []Picture
	err := orm.Join("LEFT", "pad_picture", "picture.id=pad_picture.picture_id").Where("pad_picture.pad_id=?", self.Id).Find(&pictures, &Picture{})
	PanicIf(err)
	self.Pictures = pictures
}

func (self *Pad) SearchByPage() ([]Pad, int64, error) {
	total, err := orm.Count(self)
	var pads []Pad
	err = orm.OrderBy(self.GetSortProperties()[0].Column+" "+self.GetSortProperties()[0].Direction).Limit(self.GetPageSize(), self.GetDisplayStart()).Find(&pads, self)
	BatchLoadPictures(pads)
	return pads, total, err
}

func (pad Pad) Validate(errors *binding.Errors, r *http.Request) {
	if len(pad.Name) > 20 {
		errors.Fields["name"] = "Length of name should be less than 20."
	}
}
