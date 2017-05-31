Hello,

The following package has been {{.Action}}:

Id:		{{.Package.BuildDate.Year}}-{{.Package.Id}}
Name:	{{.Package.Name}}/{{.Package.Architecture}}
For:		{{.Package.Platform}}/{{.Package.Repo}}
Type:	{{.Package.Type}}
Built:	{{.Package.BuildDate}}

{{if ne .Package.Status "testing"}}Comments on the package:
{{range .Package.Karma}}
* {{.User.Email}} voted {{.Vote | convertKarma}} and commented: "{{.Comment}}".
{{end}}{{end}}

More information available at the Kahinah website:
{{.KahinahUrl}}/builds/{{.Package.Id}}
