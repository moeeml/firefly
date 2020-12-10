package controller

import "github.com/kataras/iris/v12/mvc"

type Index struct {
	Basic
}

func (c *Index) Get() mvc.Result {
	return mvc.View{Name: "index.html"}
}
