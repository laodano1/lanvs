package controllers

import "github.com/astaxie/beego"

type (
	HomePageController struct {
		beego.Controller
	}
)

func (c *HomePageController) Get() {
	beego.Info("in HomePageController method Get")

	c.TplName = "bootstrap.html"
	//c.TplName = "bootstrap.tpl"
}
