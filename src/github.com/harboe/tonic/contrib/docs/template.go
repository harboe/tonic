package docs

import (
	"io"
	"log"
	"text/template"
)

const text_template = `
	<h1>{{.Name}}</h1>
	<p>{{.Description}}</p>

	<h2>Index:</h2>
	<dl>{{ range $i,$r := .Routes }}
		<dd><a href="#{{$r.Name}}">{{$r.Name}}</a></dd>
	{{ end }}</dl>
	{{ range $r := .Routes }}
	<h2 id="{{$r.Name}}">{{$r.Name}}</h2>
	<pre><code>{{$r.Method}} {{$r.Path}}</code></pre>
	<p>{{$r.Description}}</p>

	<strong>Parameters:</strong>
	<table style="width: 100%">
	<tr>
		<th>Name</th>
		<th>Type</th>
		<th>Description</th>
	</tr>{{ range $p := $r.Parameters }}
	<tr>
		<td>{{$p.Name}}</td>
		<td>{{$p.Type}}</td>
		<td>
			{{if not $p.Optional}}<strong>Required.</strong>{{end}}
			{{if $p.Multiple}}<i>Multiple.</i>{{end}}
			{{$p.Description}}</td>
	</tr>{{end}}
	</table>
	{{end}}
`

func ToText(w io.Writer, pkg Package) {
	// Create a new template and parse the letter into it.
	t := template.Must(
		template.New("doc").Parse(text_template))

	if err := t.Execute(w, pkg); err != nil {
		log.Fatal("executing template:", err)
	}
}
