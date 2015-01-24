package model

import (
	. "github.com/easykoo/go-display/common"

	"testing"
)

func Init() {
	SetConfig()
	SetLog()
	SetEngine()
}

func Test_GetHotBlog(t *testing.T) {
	Init()

	Log.Debug(blog)
}
