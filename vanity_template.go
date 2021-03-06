package main

import "html/template"

// VanityTemplateStr is the HTML template for our responses.
const VanityTemplateStr = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="description" content="goovus - Go Open Vanity URL Server" />
	<meta name="go-import" content="{{.Root}} {{.VCS}} {{.RepoURL}}" />
    <title>goovus - Go Open Vanity URL Server</title>
  </head>
  <body>
    <div id="root">
		<pre>root= {{.Root}}</pre>
		<pre>vcs= {{.VCS}}</pre>
		<pre>repoURL= {{.RepoURL}}</pre>
	</div>
  </body>
</html>
`

// VanityTemplateData is the data to give when executing the template.
type VanityTemplateData struct {
	Root    string
	VCS     string
	RepoURL string
}

// VanityTemplate is the template instance for TemplateStr.
var VanityTemplate = template.Must(template.New("").Parse(VanityTemplateStr))
