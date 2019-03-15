package routers

import (
	"lanvs/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/",    &controllers.HomePageController{})
    beego.Router("/abc", &controllers.MainController{})
}
