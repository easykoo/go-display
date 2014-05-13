package handler

import (
	"github.com/go-martini/martini"

	. "github.com/easykoo/go-display/common"
	"github.com/easykoo/go-display/middleware"
	"github.com/easykoo/go-display/model"

	"encoding/json"
	"io"
	"os"
	"strings"
)

type Sizer interface {
	Size() int64
}

func Picture(ctx *middleware.Context) {
	switch ctx.R.Method {
	case "POST":
		picture := new(model.Picture)
		picture.SetPageActive(true)
		picture.SetPageSize(ParseInt(ctx.R.FormValue("iDisplayLength")))
		picture.SetDisplayStart(ParseInt(ctx.R.FormValue("iDisplayStart")))
		columnNum := ctx.R.FormValue("iSortCol_0")
		sortColumn := ctx.R.FormValue("mDataProp_" + columnNum)
		picture.AddSortProperty(sortColumn, ctx.R.FormValue("sSortDir_0"))
		pictures, total, err := picture.SearchByPage()
		PanicIf(err)
		ctx.Set("aaData", pictures)
		ctx.Set("iTotalDisplayRecords", total)
		ctx.Set("iTotalRecords", total)
		ctx.JSON(200, ctx.Response)
	default:
		ctx.HTML(200, "picture/picture", ctx)
	}
}

func Upload(ctx *middleware.Context) {
	file, handle, err := ctx.R.FormFile("image")
	PanicIf(err)
	defer file.Close()

	filename := handle.Filename
	ext := SubString(filename, strings.LastIndex(filename, "."), 4)
	filename = GetGuid() + ext
	os.Mkdir("public/pictures", 0777)
	f, err := os.Create("public/pictures/" + filename)
	defer f.Close()
	PanicIf(err)
	io.Copy(f, file)
	Log.Debug("上传文件的大小为: %d", file.(Sizer).Size())
	url := "/pictures/" + filename
	ctx.Set("url", url)
	ctx.JSON(200, ctx.Response)
}

func NewPicture(ctx *middleware.Context, picture model.Picture) {
	switch ctx.R.Method {
	case "POST":
		err := picture.Insert()
		PanicIf(err)
		ctx.Redirect("/picture")
	default:
		ctx.HTML(200, "picture/edit", ctx)
	}
}

func EditPictureView(ctx *middleware.Context, params martini.Params) {
	id := params["id"]
	picture := new(model.Picture)
	picture.Id = ParseInt(id)
	pic, err := picture.Get()
	PanicIf(err)
	ctx.Set("Picture", pic)
	ctx.HTML(200, "picture/edit", ctx)
}

func EditPicture(ctx *middleware.Context, picture model.Picture) {
	err := picture.Update()
	PanicIf(err)
	ctx.Redirect("/picture")
}

func DeletePicture(ctx *middleware.Context, params martini.Params) {
	id := params["id"]
	picture := new(model.Picture)
	picture.Id = ParseInt(id)
	err := picture.Delete()
	PanicIf(err)
	ctx.Set("success", true)
	ctx.JSON(200, ctx.Response)
}

func DeletePictures(ctx *middleware.Context) {
	pictures := ctx.R.FormValue("pictureArray")
	var res []int
	json.Unmarshal([]byte(pictures), &res)
	picture := new(model.Picture)
	err := picture.DeletePictures(res)
	PanicIf(err)
	ctx.Set("success", true)
	ctx.JSON(200, ctx.Response)
}
