Hello,

The package you have built has been {{.Action}}:

Id:		UPDATE-{{.Package.BuildDate.Year}}-{{.Package.Id}}
Name:	{{.Package.Name}}/{{.Package.Architecture}}
For:		{{.Package.Platform}}/{{.Package.Repo}}
Type:	{{.Package.Type}}
Built:	{{.Package.BuildDate}}

{{if eq .Package.Status "testing"}}Be sure to vote for your own package - and remember that QA
may need more information: the comments section has input that
might prove valuable.{{else}}Comments on the package:
{{range .Package.Karma}}
* {{.User.Email}} voted {{.Vote | convertKarma}} and commented: "{{.Comment}}".
{{end}}{{end}}

More information available at the Kahinah website:
{{.KahinahUrl}}/builds/{{.Package.Id}}
