package handler

import (
	"github.com/go-martini/martini"

	. "github.com/easykoo/go-display/common"
	"github.com/easykoo/go-display/middleware"
	"github.com/easykoo/go-display/model"

	"encoding/json"
	"strings"
)

func Pad(ctx *middleware.Context) {
	switch ctx.R.Method {
	case "POST":
		pad := new(model.Pad)
		pad.SetPageActive(true)
		pad.SetPageSize(ParseInt(ctx.R.FormValue("iDisplayLength")))
		pad.SetDisplayStart(ParseInt(ctx.R.FormValue("iDisplayStart")))
		columnNum := ctx.R.FormValue("iSortCol_0")
		sortColumn := ctx.R.FormValue("mDataProp_" + columnNum)
		pad.AddSortProperty(sortColumn, ctx.R.FormValue("sSortDir_0"))
		pads, total, err := pad.SearchByPage()
		PanicIf(err)
		ctx.Set("aaData", pads)
		ctx.Set("iTotalDisplayRecords", total)
		ctx.Set("iTotalRecords", total)
		ctx.JSON(200, ctx.Response)
	default:
		ctx.HTML(200, "pad/pad", ctx)
	}
}

func NewPad(ctx *middleware.Context, pad model.Pad) {
	switch ctx.R.Method {
	case "POST":
		err := pad.Insert()
		PanicIf(err)
		ctx.Redirect("/pad")
	default:
		ctx.HTML(200, "pad/edit", ctx)
	}
}

func EditPadView(ctx *middleware.Context, params martini.Params) {
	id := params["id"]
	pad := new(model.Pad)
	pad.Id = ParseInt(id)
	pic, err := pad.GetById()
	PanicIf(err)
	ctx.Set("Pad", pic)
	ctx.HTML(200, "pad/edit", ctx)
}

func EditPad(ctx *middleware.Context, pad model.Pad) {
	err := pad.Update()
	PanicIf(err)
	ctx.Redirect("/pad")
}

func DeletePad(ctx *middleware.Context, params martini.Params) {
	id := params["id"]
	pad := new(model.Pad)
	pad.Id = ParseInt(id)
	err := pad.Delete()
	PanicIf(err)
	ctx.Set("success", true)
	ctx.JSON(200, ctx.Response)
}

func DeletePads(ctx *middleware.Context) {
	pads := ctx.R.FormValue("padArray")
	var res []int
	json.Unmarshal([]byte(pads), &res)
	pad := new(model.Pad)
	err := pad.DeletePads(res)
	PanicIf(err)
	ctx.Set("success", true)
	ctx.JSON(200, ctx.Response)
}

func ChoosePicture(ctx *middleware.Context, pad model.Pad) {
	if pad.PictureStr != "" {
		var pictureIds []string
		pictureIds = strings.Split(pad.PictureStr, ",")
		for _, pictureId := range pictureIds {
			pad.Pictures = append(pad.Pictures, model.Picture{Id: ParseInt(pictureId)})
		}
	}
	err := pad.UpdatePictures()
	PanicIf(err)
	ctx.Redirect("/pad")
}

func ChoosePictureView(ctx *middleware.Context, params martini.Params) {
	id := params["id"]
	pad := new(model.Pad)
	pad.Id = ParseInt(id)
	pic, err := pad.GetById()
	PanicIf(err)
	ctx.Set("Pad", pic)
	ctx.HTML(200, "pad/choosePicture", ctx)
}

func Info(ctx *middleware.Context, params martini.Params) {
	pad := new(model.Pad)
	name := params["name"]
	pad.Name = name
	pad, err := pad.Get()
	PanicIf(err)
	if pad == nil {
		ctx.Set("success", false)
		ctx.Set("message", "Pad not exists!")
		ctx.JSON(200, ctx.Response)
	} else {
		ctx.Set("success", true)
		ctx.Set("Pad", pad)
		ctx.JSON(200, ctx.Response)
	}
}
