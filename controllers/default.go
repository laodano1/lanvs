package controllers

import (
	"github.com/astaxie/beego"
	"os"
	"log"
	"bufio"
	"io"
)


type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	//c.Data["Website"] = "beego.me"
	//c.Data["Email"] = "astaxie@gmail.com"
	//c.TplName = "index.tpl"

	//c.ServeJSON()
	//c.Ctx.Output.Body()
	vdo := "/Users/sophia/my_files/Media/[阳光电影www.ygdy8.com].破·局.HD.720p.国语中字.mkv"
	beego.Info("play video:", vdo)

	f, err := os.Open(vdo)
	if err != nil {
		log.Println("open file error:", err)
	}
	rd := bufio.NewReader(f)

	//by, err := ioutil.ReadAll(f)
	//rd := bytes.NewReader(by)
	//http.ServeContent(c.Ctx.ResponseWriter, c.Ctx.Request, vdo, time.Now(), rd)

	//http.ServeFile(c.Ctx.ResponseWriter, c.Ctx.Request, vdo)
	//c.Ctx.ResponseWriter.Header().Set("Content-Type", "video/mp4")
	//c.Ctx.ResponseWriter.Header().Set("Content-Type", "video/mp4")
	//c.Ctx.ResponseWriter.Header().Set("Content-Type", "video/x-matroska")
	c.Ctx.ResponseWriter.Header().Set("Content-Type", "video/webm")

	io.Copy(c.Ctx.ResponseWriter, rd)
}
