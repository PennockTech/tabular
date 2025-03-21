// Copyright © 2016,2025 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

/*
The html wrapper provides a means for generating HTML from a tabular table.
Like everything else, the interface is still very much in flux.
Where texttable used properties on the table to control rendering, HTMLTable
uses a wrapper object which methods can be set upon.
I want to see which approach is "better".
*/
package html // import "go.pennock.tech/tabular/html"

import (
	"bytes"
	"html/template"
	"io"
	"strings"

	"go.pennock.tech/tabular"
	"go.pennock.tech/tabular/color"
	"go.pennock.tech/tabular/properties"
)

// HTMLTable wraps a tabular Table to provide some extra information used
// in rendering to HTML.  Id and Class are properties of the top-level table.
// Caption will be inserted if present.
// TemplateName can be used if you are managing general html/template namespaces,
// else the template will be unnamed.
type HTMLTable struct {
	tabular.Table
	Id           string
	Class        string
	Caption      string
	TemplateName string

	rowClassGenerator func(rowNum int, ctx any) template.HTMLAttr
	rowClassCtx       any

	template *template.Template

	cachedOmitColumns []bool
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
	callable func(rowNum int, ctx any) template.HTMLAttr,
	userCtx any,
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
<table {{- with .Class}} class="{{.}}"{{end}} {{- with .Id}} id="{{.}}"{{end}} {{- with (BGColor Table) }} style="background-color: {{.}}"{{end}}>
{{- with .Caption}}
  <caption>{{.}}</caption>
{{- end}}
  <colgroup>
{{- range $n, $hdr := Headers}}<col class="{{ColumnClass $hdr}}" {{- with (BGColor (Column $n)) }} style="background-color: {{.}}"{{end}} />{{end -}}
  </colgroup>
  <thead>
    <tr {{- if .HaveRowClass}} class="{{RowClass 0}}"{{end}}>
{{- range Headers}}<th>{{.}}</th>{{end -}}
    </tr>
  </thead>
  <tbody>
{{- range $i, $row := Rows}}{{if OmitRow $row | not}}{{if $row.IsSeparator | not}}
    <tr {{- if $.HaveRowClass}} class="{{RowClass (OnePlus $i)}}"{{end}} {{- with (BGColor .) }} style="background-color: {{.}}"{{end}}>
{{- range $cell := CellsOf $row }}<td {{- with (BGColor $cell) }} style="background-color: {{.}}"{{end}}>{{$cell}}</td>{{end -}}
    </tr>
{{- end}}{{end}}{{end}}
  </tbody>
</table>
`

func cellsNotOmitted(cells []tabular.Cell, omit []bool) []*tabular.Cell {
	r := make([]*tabular.Cell, 0, len(cells))
	for i := range cells {
		if !omit[i] {
			r = append(r, &cells[i])
		}
	}
	return r
}

func cellToColumnClass(cell *tabular.Cell) template.HTMLAttr {
	s := "col-" + strings.Replace(cell.String(), " ", "-", -1)
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r > 127 {
			continue
		}
		switch {
		case 'A' <= r && r <= 'Z':
		case 'a' <= r && r <= 'z':
		case '0' <= r && r <= '9':
		case r == '-', r == '_':
		case r == ' ' || r == '\t':
			b.WriteRune('-')
			continue
		default:
			// skip whatever's not allow-listed above
			continue
		}
		b.WriteRune(r)
	}

	return template.HTMLAttr(b.String())
}

func lookupColor(item any, prop any) string {
	var (
		colRaw any
		col    color.Color
		ok     bool
	)
	switch container := item.(type) {
	case HTMLTable:
		colRaw = container.GetProperty(prop)
	case tabular.Table:
		colRaw = container.GetProperty(prop)
	case *tabular.ATable:
		colRaw = container.GetProperty(prop)
	case *tabular.Row:
		colRaw = container.GetProperty(prop)
	case tabular.Cell:
		colRaw = container.GetProperty(prop)
	case *tabular.Cell:
		colRaw = container.GetProperty(prop)
	case *tabular.Column:
		colRaw = container.GetProperty(prop)
	default:
		return ""
	}
	if colRaw == nil {
		return ""
	}
	if col, ok = colRaw.(color.Color); ok {
		return col.HTML()
	}
	return ""
}

func (ht *HTMLTable) getFuncs() template.FuncMap {
	return template.FuncMap{
		"Table":       func() tabular.Table { return ht.Table },
		"Headers":     func() []*tabular.Cell { return cellsNotOmitted(ht.Table.Headers(), ht.cachedOmitColumns) },
		"RowClass":    func(i int) template.HTMLAttr { return ht.rowClassGenerator(i, ht.rowClassCtx) },
		"Column":      func(i int) *tabular.Column { return ht.Table.Column(i + 1) },
		"ColumnClass": cellToColumnClass,
		"CellsOf":     func(r *tabular.Row) []*tabular.Cell { return cellsNotOmitted(r.Cells(), ht.cachedOmitColumns) },
		"OnePlus":     func(i int) int { return i + 1 },
		"Rows":        func() []*tabular.Row { return ht.Table.AllRows() },
		"OmitRow": func(r *tabular.Row) (bool, error) {
			return properties.ExpectBoolPropertyOrNil(properties.Omit, r.GetProperty(properties.Omit), "html:OmitRow", "row", 0)
		},
		// want: "FGColor": func[T tabular.Table|*tabular.Cell|*tabular.Row]() string {}
		"FGColor": func(item any) string { return lookupColor(item, properties.FGColor) },
		"BGColor": func(item any) string { return lookupColor(item, properties.BGColor) },
	}
}

// RenderTo writes the table to the provided writer, stopping if it should encounter an error.
func (ht *HTMLTable) RenderTo(w io.Writer) (err error) {
	ht.InvokeRenderCallbacks()

	if err = ht.populateOmitCache(); err != nil {
		return
	}
	defer func() {
		ht.cachedOmitColumns = nil
	}()

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

func (ht *HTMLTable) populateOmitCache() error {
	ht.cachedOmitColumns = make([]bool, ht.NColumns())

	var (
		defaultOmit bool
		err         error
	)

	if defaultOmit, err = properties.ExpectBoolPropertyOrNil(
		properties.Omit,
		ht.Column(0).GetProperty(properties.Omit),
		"html:RenderTo", "default column", 0); err != nil {
		return err
	}

	for i := range ht.NColumns() {
		omit := ht.Column(i + 1).GetProperty(properties.Omit)
		if omit != nil {
			if ht.cachedOmitColumns[i], err = properties.ExpectBoolPropertyOrNil(properties.Omit, omit, "html:RenderTo", "column", i+1); err != nil {
				return err
			}
		} else {
			ht.cachedOmitColumns[i] = defaultOmit
		}
	}

	return nil
}
