package util

import (
	"bytes"
	"html/template"

	beego "github.com/beego/beego/v2/adapter"
)

var (
	outwardloc = beego.AppConfig.String("outwardloc")
	PREFIX     = beego.AppConfig.String("urlprefix")
)

func GetPrefixStringWithData(dest string, data interface{}) string {
	// no need to prefix if the dest has no / before it
	temp := template.Must(template.New("prefixTemplate").Parse(dest))
	var b bytes.Buffer

	err := temp.Execute(&b, data)
	if err != nil {
		panic(err)
	}

	result := b.String()
	return GetPrefixString(result)
}

func GetPrefixString(dest string) string {
	if PREFIX == "" {
		return dest
	}

	return "/" + PREFIX + dest
}

func GetFullUrlStringWithData(dest string, data interface{}) string {
	return outwardloc + "/" + GetPrefixStringWithData(dest, data)
}

func GetFullUrlString(dest string) string {
	return outwardloc + "/" + GetPrefixString(dest)
}
