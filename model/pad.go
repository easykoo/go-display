package model

import (
	"github.com/easykoo/binding"

	. "github.com/easykoo/go-display/common"

	"net/http"
	"time"
)

type Pad struct {
	Id          int       `form:"padId" xorm:"int(11) pk not null autoincr"`
	Name        string    `form:"name" xorm:"varchar(20) not null"`
	Ip          string    `form:"ip" xorm:"varchar(15) not null"`
	Picture     Picture   `json:"picture_id" xorm:"picture_id int(11) default null"`
	Color       string    `form:"color" xorm:"varchar(20) not null"`
	Description string    `form:"description" xorm:"varchar(100) not null"`
	CreateUser  string    `xorm:"varchar(20) default 'SYSTEM'"`
	CreateDate  time.Time `xorm:"datetime created"`
	UpdateUser  string    `xorm:"varchar(20) default 'SYSTEM'"`
	UpdateDate  time.Time `xorm:"datetime updated"`
	Version     int       `form:"version" xorm:"int(11) version"`
	Page        `xorm:"-"`
}

func (self *Pad) Exist() (bool, error) {
	return orm.Get(self)
}

func (self *Pad) Get() (*Pad, error) {
	pad := &Pad{}
	_, err := orm.Id(self.Id).Get(pad)
	return pad, err
}

func (self *Pad) GetByIp() (*Pad, error) {
	//pad := &Pad{}
	_, err := orm.Get(self)
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

func (self *Pad) Delete() error {
	_, err := orm.Delete(self)
	Log.Info("Pad ", self.Name, " deleted")
	return err
}

func (self *Pad) DeletePads(array []int) error {
	_, err := orm.In("id", array).Delete(&Pad{})
	Log.Info("Pads: ", array, " deleted")
	return err
}

func (self *Pad) SelectAll() ([]Pad, error) {
	var pads []Pad
	err := orm.Find(&pads)
	return pads, err
}

func (self *Pad) SearchByPage() ([]Pad, int64, error) {
	total, err := orm.Count(self)
	var pads []Pad
	err = orm.OrderBy(self.GetSortProperties()[0].Column+" "+self.GetSortProperties()[0].Direction).Limit(self.GetPageSize(), self.GetDisplayStart()).Find(&pads, self)
	return pads, total, err
}

func (pad Pad) Validate(errors *binding.Errors, r *http.Request) {
	if len(pad.Name) > 20 {
		errors.Fields["name"] = "Length of name should be less than 20."
	}
}

func (self *Pad) ChoosePicture(padArray []Pad) error {
	_, err := orm.Update(padArray)
	Log.Info("Pads: ", " updated")
	return err
}

func (self *Pad) ChoosePicture1(array []int, pictureId int) error {
	sql := "update pad set picture_id =" + IntString(pictureId) + " where id in ("
	for i := 0; i < len(array); i++ {
		sql += IntString(array[i])
		if i < len(array)-1 {
			sql += ","
		}
	}
	sql += ")"
	_, err := orm.Exec(sql)
	Log.Info("Pads: ", " updated")
	return err
}
