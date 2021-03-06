package webTypeAdmin

import "github.com/bronze1man/kmg/kmgHtmlTemplate"

type selectTemplateData struct {
	List  []string
	Value string
}

var theTemplate = kmgHtmlTemplate.MustNewSingle(`
{{ define "Main" }}
<!DOCTYPE html>
<html>
<head>
	<title>{{.Title}}</title>
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<!-- Bootstrap -->
	<link href="public/css/bootstrap.min.css" rel="stylesheet">
<style>
.slice{
	background-color: #DFF0D8 !important;
}
.struct{
	background-color: #F2DEDE !important;
}
</style>
</head>
<body>
	<h1>{{.Title}}</h1>
	<div style="display:none" class="alert alert-danger" id="errors-msg-container">
		<button type="button" class="close" data-dismiss="alert" aria-hidden="true" id="errors-msg-close-btn">&times;</button>
		<div id="errors-msg"></div>
	</div>
	{{.InjectHtml}}
	<div class="kmg-type-admin-root" {{if .Path}}data-path="{{.Path}}"{{end}} >
		{{.Html}}
	</div>
	<script src="public/js/jquery-2.0.3.js"></script>
	<script src="public/js/bootstrap.min.js"></script>
	<script src="public/js/kmgTypeAdmin.js"></script>
</body>
</html>
{{ end }}


{{ define "TextInput" }}
<input type="text" class="form-control kmg-single-input input-ms" value="{{.}}"/>
{{ end }}


{{ define "Slice" }}
<table class="table-bordered slice table-condensed table" ><tbody>
	<tr data-path="">
		<td><button type="button" class="btn btn-primary kmg-create-action btn-xs">Create</button></td>
		<td></td>
		<td></td>
	</tr>
	{{range .}}
	<tr data-path="{{.Path}}">
		<td><button type="button" class="btn btn-danger kmg-delete-action btn-xs">Delete</button></td>
		<td>{{.Index}}</td>
		<td>{{.Html}}</td>
	</tr>
	{{end}}
</tbody></table>
{{ end }}


{{ define "Array" }}
<table class="table-bordered slice table-condensed table" ><tbody>
	{{range .}}
	<tr data-path="{{.Path}}">
		<td>{{.Index}}</td>
		<td>{{.Html}}</td>
	</tr>
	{{end}}
</tbody></table>
{{ end }}


{{ define "Struct" }}
<table class="table-bordered struct table-condensed table"><tbody>
	{{range .}}
	<tr data-path="{{.Path}}">
		<td>{{.Name}}</td>
		<td>{{.Html}}</td>
	</tr>
	{{end}}
</tbody></table>
{{ end }}


{{ define "NilPtr" }}
<div data-path="ptr">
<button type="button" class="btn btn-primary kmg-create-action btn-xs">New</button>
</div>
{{ end }}

{{ define "Ptr" }}
<div data-path="ptr">
{{.}}
</div>
{{ end }}

{{ define "Select" }}
<select class="form-control kmg-single-input input-ms">
  {{range .List}}
    <option {{if eq . $.Value}}selected="selected"{{end}}>{{.}}</option>
  {{end}}
</select>
{{ end }}


{{ define "Map" }}
<table class="table-bordered slice table-condensed table" ><tbody>
	<tr class="kmg-map-create-parent" data-path="">
		<td><button type="button" class="btn btn-primary kmg-create-action btn-xs">Create</button></td>
		<td></td>
		<td><input type="text" class="form-control kmg-map-create-key input-ms"/></td>
	</tr>
	{{range .}}
	<tr data-path="{{.Path}}">
		<td><button type="button" class="btn btn-danger kmg-delete-action btn-xs">Delete</button></td>
		<td>{{.Key}}</td>
		<td>{{.Html}}</td>
	</tr>
	{{end}}
</tbody></table>
{{ end }}
`)
