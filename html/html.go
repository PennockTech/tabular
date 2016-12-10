// Copyright © 2016 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

/*
The html wrapper provides a means for generating HTML from a tabular table.

We mostly avoid specifying inline-CSS, instead making it easy to class/id
decorate to let you provide your own styling wrappers.
However, if you specify alignment constraints in the table, those have to
become styles as the `align` attribute is not allowed in HTML5.
If you need to specify alignments but also need to block styles in the HTML
(eg, "Content-Security-Policy" compatibility) then set the .DisableStyle
in the HTMLTable.
*/
package html // import "go.pennock.tech/tabular/html"

import (
	"bytes"
	"html/template"
	"io"

	"go.pennock.tech/tabular"
)

// HTMLTable wraps a tabular Table to provide some extra information used
// in rendering to HTML.  Id and Class are properties of the top-level table.
// Caption will be inserted if present.
// TemplateName can be used if you are managing general html/template namespaces,
// else the template will be unnamed.
// DisableStyle inhibits inline-CSS for the cases where it might arise
// (carefully avoided by default, but triggered by alignments).
type HTMLTable struct {
	tabular.Table
	Id           string
	Class        string
	Caption      string
	TemplateName string
	DisableStyle bool

	rowClassGenerator func(rowNum int, ctx interface{}) template.HTMLAttr
	rowClassCtx       interface{}

	template *template.Template
}

// Wrap returns an HTMLTable rendering object for the given tabular Table.
func Wrap(t tabular.Table) *HTMLTable {
	return &HTMLTable{Table: t}
}

// New returns an HTMLTable with a new Table inside it, access via .Table
// or just use the interface methods on the HTMLTable.
func New() *HTMLTable {
	return Wrap(tabular.New())
}

// SetRowClassGenerator is used to register a user function to be used to emit
// classes for each row.  The HTMLTable object is returned, to permit chaining.
//
// The callable is passed the row number (starting from 0 for the header, 1 for
// the first body row) and whatever object is passed as the context here, which
// may be used for persisting state.
//
// Note that use of a context here means that an HTMLTable can not be concurrent
// rendered from two threads (unless you're doing something very strange and
// handle all locking in the callable, using a mutex in the context, yourself).
// Instead, generate a new HTMLTable wrapper for each table.
//
// Separator rows of the source table are not emitted in HTML tables; they do
// count for indexing of row-numbers, so there may be row-number gaps.  If you
// need to alternate row classes, either keep a flip-flop in the context or
// detect the skipped row-numbers (last-seen in context) and handle specially.
//
// The callable's return should be an html/template.HTMLAttr; this is not
// coerced in this library, to ensure that people writing the callbacks see
// the data-safety type at the time of implementation, to provoke careful thought.
func (ht *HTMLTable) SetRowClassGenerator(
	callable func(rowNum int, ctx interface{}) template.HTMLAttr,
	userCtx interface{},
) *HTMLTable {
	ht.rowClassGenerator = callable
	ht.rowClassCtx = userCtx
	return ht
}

// Render takes an HTMLTable and returns a string representing the fully
// rendered table, or an error.
func (ht *HTMLTable) Render() (string, error) {
	b := &bytes.Buffer{}
	err := ht.RenderTo(b)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

const rawTableTemplateStr = `{{/**/ -}}
<table {{- with .Class}} class="{{.}}"{{end}} {{- with .Id}} id="{{.}}"{{end}}>
{{- with .Caption}}
  <caption>{{.}}</caption>
{{- end}}
  <thead>
    <tr {{- if .HaveRowClass}} class="{{RowClass 0}}"{{end}}>
{{- range Headers}}<th>{{.}}</th>{{end -}}
	</tr>
  </thead>
  <tbody>
{{- range $i, $row := Rows}}{{if $row.IsSeparator | not}}
    <tr {{- if $.HaveRowClass}} class="{{RowClass (OnePlus $i)}}"{{end}}>
{{- range CellsOf $row }}<td>{{.}}</td>{{end -}}
    </tr>
{{- end}}{{end}}
  </tbody>
</table>
`

func cellsToStringArray(cells []tabular.Cell) []string {
	r := make([]string, len(cells))
	for i := range cells {
		r[i] = cells[i].String()
	}
	return r
}

func tabularAlignToHTMLAlign(po tabular.PropertyOwner) (template.HTMLAttr, bool) {
	s, ok := tabular.GetStaticAlignment(po)
	if !ok {
		return template.HTMLAttr("LEFT"), ok // the string shouldn't be used, but let's all-caps it to stand out just-in-case
	}
	switch s {
	case tabular.ALIGN_LEFT:
		return template.HTMLAttr("left"), true
	case tabular.ALIGN_CENTER:
		return template.HTMLAttr("center"), true
	case tabular.ALIGN_RIGHT:
		return template.HTMLAttr("right"), true
	case tabular.ALIGN_PERIOD:
		return template.HTMLAttr("."), true
	case tabular.ALIGN_COMMA:
		return template.HTMLAttr(","), true
	}
	return template.HTMLAttr("LefT"), false
}

func (ht *HTMLTable) getFuncs() template.FuncMap {
	return template.FuncMap{
		"Headers":          func() []string { return cellsToStringArray(ht.Table.Headers()) },
		"RowClass":         func(i int) template.HTMLAttr { return ht.rowClassGenerator(i, ht.rowClassCtx) },
		"CellsOf":          func(r *tabular.Row) []string { return cellsToStringArray(r.Cells()) },
		"OnePlus":          func(i int) int { return i + 1 },
		"Rows":             func() []*tabular.Row { return ht.Table.AllRows() },
		"ColAlignment":     func(c int) string { attr, _ := tabularAlignToHTMLAlign(ht.columns[c]); return attr },
		"HaveColAligmnent": func(c int) bool { _, have := tabularAlignToHTMLAlign(ht.columns[c]); return have },
	}
}

// RenderTo writes the table to the provided writer, stopping if it should encounter an error.
func (ht *HTMLTable) RenderTo(w io.Writer) (err error) {
	ht.InvokeRenderCallbacks()

	if ht.template == nil {
		ht.template, err = template.New(ht.TemplateName).Funcs(ht.getFuncs()).Parse(rawTableTemplateStr)
		if err != nil {
			return
		}
	} else {
		// getFuncs will set up method closures for some functions; overwrite each time
		ht.template.Funcs(ht.getFuncs())
	}

	renderData := struct {
		Id, Class, Caption string
		HaveRowClass       bool
	}{
		Id:           ht.Id,
		Class:        ht.Class,
		Caption:      ht.Caption,
		HaveRowClass: ht.rowClassGenerator != nil,
	}

	return ht.template.Execute(w, renderData)
}
