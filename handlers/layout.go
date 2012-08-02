package handlers

import (
	. "github.com/paulbellamy/mango"
	"github.com/sunfmin/mangotemplate"
	"html/template"
)

type provider struct {
}

type Header struct {
}

func (p *provider) LayoutData(env Env) interface{} {
	return &Header{}
}

func LayoutAndRenderer() (l Middleware, r Middleware) {
	tpl, err := template.ParseGlob("templates/*/*.html")
	if err != nil {
		panic(err)
	}
	l = mangotemplate.MakeLayout(tpl, "main", &provider{})
	r = mangotemplate.MakeRenderer(tpl)
	return
}
