Tabular Overview
================

At the core, tabular knows nothing about rendering.  The base-level library,
in the repo's top-level, can be imported and used to create tables and add
content.  Rendering requires using another layer to convert the output for
display.

A `Table` is an interface.  There is one core public type which implements the
interface.  This allows the core type to be embedded in the wrapper/display
objects and for those to satisfy the table interface, thus being tables
themselves.  Rows and cells are not interfaces.  The core public type for a
table is, imaginatively, `*ATable`.

Callers doing simple table usage should import the `auto` sub-package and use
the `RenderTable` interface type, or use tabular's code `Table` interface
(which is part of `RenderTable`).

A table consists of rows of cells and some metadata.  The metadata includes
virtual columns, allowing for addressing by column too.  Columns are
identified by the header name.  There is only one (or zero) header row per
table.

Errors in adding data are usually not reported immediately, to let data stream
in.  Instead, errors accumulate in an error holder.  Rows hold errors, but
once a row is part of a table, its errors become the tables' errors (and the
error container is diverted to be the table's).  An error container can be
interrogated for its list of current errors.

The errors are either a list of non-nil errors, or nil.  An empty list should
never be returned.  If a nil is returned in the list of errors then that is a
bug in tabular.

A row either contains cells or is a "special" row.  The only type of special
row is a "separator" row.  The tabular layer itself doesn't know what a
separator row is, beyond that it exists and a row can be one.  The
`AddSeparator()` table method adds one, the `IsSeparator()` row method asks a
row if it is one.  The `Cells()` method, which returns an array of cells, can
return nil if and only if the row is special (ie, at present, a separator).  A
real row is always a splice of cells, even if that splice is empty.

A `Cell` contains "an object".  That object can be a string, something which
satisfies `Stringer` or `GoStringer`, a rune, or another `Cell`.  Cells can
contain cells and this is intended to allow for dynamic update, based upon
evaluation.  Defining a `MarshalText` method may be advisable and a future API
bump might choose to prefer `MarshalText` to `Stringer` or `GoStringer`.

If a `Cell` contains an object then various rendering layers may make use of
other interfaces satisfied by that object to determine how to display it;
loosely, think of "width" and "height", but this will be covered in more
detail below.

Cells, Rows, Columns and Tables can have "properties" set upon them.
Properties are namespaced objects, very similar to Golang's net contexts.
Clients of the tabular package are free to decorate items with whatever
properties they want.

The tabular package supports automatically updating properties at "addition"
time and at "render" time.  This is done by setting callbacks.  Callbacks can
be on a table or a row.  Within the table, they can be registered for use on a
table or a column, for when a row is added, or when a cell is added.

The cell's location in the grid is not a property, but is available via a
method call upon the cell.

The callbacks and properties should not be exposed to end-users.

All child objects have links back to their containers.  This is used, eg, to
be able to get column information for a given cell.  This does mean that there
are ownership loops.


Sub-package Commonalities
-------------------------

In all cases:

* There is a `Wrap()` function which takes any `tabular.Table` and returns
  a wrapper object for this sub-package.
* There is a `New()` function which generates an empty table for this
  sub-package.
* The wrapper objects have `Render()` and `RenderTo()` methods, and the
  packages have top-level functions which create a wrapper object, with
  default options, and calls the object methods.
* The `Render()` method will return a string of the rendered text, together
  with an error.
* The `RenderTo(io.Writer)` method uses a stream-based approach and only
  returns an error.
* The wrapper object's type is named with a `FooTable` naming style, accepting
  that this causes some stuttering.  This is acceptable because most callers
  should never need to specify the type, but instead be using the
  `tabular.Table` interface if they care at all beyond letting the type be
  inferred.  With so many variants, I went for clarity over smoothness in
  reading out the fully-qualified type name.

Because the sub-package `New()` returns an object which satisfies the
`tabular.Table` interface, it should be capable of being populated like any
other, and most callers with simple use-cases should be able to only import
the sub-package, not `tabular` itself.

It is possible to use table properties to store data, but that's more
complexity than is usually warranted.  If you want static attributes, put them
in your wrapper object.  If you want dynamic attributes, updated based upon
content, _then_ use properties, and consider how to hide this from your users.


Auto
----

This is the `auto` sub-package of `tabular`.

This provides a `RenderTable` interface, which adds `Render()` and
`RenderTo()` methods to the core `tabular.Table`.

Unlike the more-specific sub-packages, the `New()` and `Wrap()` methods take a
string argument.  The string is a style.  This string is taken to be a
dot-joined sequence of sections, where the first section is a sub-package or a
`texttable/decoration.Decoration`.

So `auto.New("csv")` returns an `auto.RenderTable` which is satisfied by a
`*csv.CSVTable`.  `auto.New("texttable.utf8-light")` is equivalent to
`auto.New("utf8-light")` and returns an `auto.RenderTable` which is satisfied
by a `*texttable.TextTable` with the decoration set to `utf8-light`.

`auto.New("csv.foo")` is currently equivalent to `auto.New("csv")` but future
extensions might pass the `foo` onto some appropriate initialization of the
`csv` package.

The `auto` sub-package does not provide specific implementations of the
`Render` or `RenderTo` sub-packages, but does provide the usual package-level
wrappers.

In addition, `auto.ListStyles()` returns a sorted list of strings, each of
which is a valid input for `auto.New(style)`.  The list is not guaranteed to
be exhaustive, but should cover the common cases.  It is guaranteed to be
exhaustive of all top-level style names (the bit before the first `.`).


Text Table Display
------------------

This is the `texttable` sub-package of `tabular`.

This system is designed to draw a pretty table on a cell-based display system,
such as a classic Unix terminal emulator.  Every "display cell" (_not_ table
cell) is a fixed pixel width and height, so using box-drawing characters,
everything can be made to line up.

If a cell's object supports the `Height()` method then that overrides a
calculation based on "newlines count + 1".  If a cell's object supports the
`TerminalCellWidth()` method then that overrides a calculation based on text,
figuring out the longest line (multi-line supported) where length is
Unicode-aware and display-width (wide char and combining char) aware.

Calculations upon cells are done at _render_ time, to examine the contents and
determine width and height for text-table purposes.  These are stored as
properties of each cell.  This is a complete table sweep before printing the
first line starts.  Then the rendering uses the properties to size itself and
print the table.

TODO: The maximum widths should become column properties of the table and the
height become a row property.

There is a `decoration` sub-package of `texttable` which has decoration styles
for rendering tables, as ASCII or as a few varieties of Unicode box-drawing.
Decoration objects can be created by callers and set directly upon the table,
or can be set by name.  The names are maintained as a registry within the
`decoration` package.  Each name is a simple string, thus typos are a
potential source of errors.  For the styles native to the `decoration`
package, package constants are exported with the names, permitting
compile-time checks to catch issues.  Eg, use `decoration.D_UTF8_LIGHT_CURVED`
instead of `"utf8-light-curved"` if you are willing to import the `decoration`
package here.  There's a trade-off between provable correctness and importing
more and clients get to choose the level they're happy with.


HTML Table Rendering
--------------------

This is the `html` sub-package of `tabular`.

It does not render "separator" rows.

Table `Id`, `Class` and `Caption` are top-level attributes.

The table is rendered using Golang's `html/template` to handle auto-escaping
of unsafe data.  If the template name matters (you are using templates more
generally) then you can use `TemplateName` on the `HTMLTable` object.

There is an `SetRowClassGenerator()` method to let you register your own
function to be used to emit a class-name on each `<tr>` of the table's body.
See the package docs for more details (you get a row index and your own
context for passing state).


CSV Rendering
-------------

This is the `csv` sub-package of `tabular`.

The rendering is compliant to RFC4180.  All fields are always quoted.  Note in
particular that newlines within strings are not `\n` escaped and double-quotes
are doubled for escaping, thus `""`.  Both of these attributes are
RFC-specified.

The `CSVTable` type is designed to be extensible to change separators,
escaping styles and more.  The _default_ is RFC4180, and that's the only
_current_ style, but we should accept PRs for any sane options, and also
"family" sets, as long as well-specified.  At present, only the
`fieldSeparator` is called out in the struct, and there are no mutators for
it.  That's not a bug, just "not yet implemented, waiting for solid
use-cases".


JSON Rendering
--------------

This is the `json` sub-package of `tabular`.

This emits structured JSON where the table is represented as an array of
objects, one object per row, and the string representations of the column
headers become the object keys.  It's an error to not have headers; it's an
error to have missing, empty or duplicate headers (as rendered to string).

This rendered passes the underlying stored items within the cells to
`encoding/json` for marshalling, so will handle arbitrary types; because we
have only required that cell types support `.String() string`, we have a
fallback of using `String()` if the marshalling returns a sequence of `{}`
corresponding to an empty struct, ie a struct with no exported fields.
If you store a struct with exported fields and have previously relied upon
`String()` being defined, then the output in JSON format will be less than
ideal unless you also define a `MarshalText()` method.

There is no handling for a `MarshalText()` or `MarshalJSON()` method choosing
to return `{}`, on the assumption that if they do so, then a `String()` method
would also choose to return `{}`.  This is currently a brute fallback for a
heuristic to Do The Right Thing, rather than reflection-aware handling.

A future revision might switch away from generic JSON marshalling and simplify
this; if we do so, then we'll support the generic `MarshalText()` but *not*
the `MarshalJSON()` legacy method.


Markdown Rendering
------------------

This is the `markdown` sub-package of `tabular`.

The dialect emitted is "GitHub Flavored Markdown" tables.

For automated production you should avoid this renderer in favor of the `html`
renderer; the markdown table support is extremely limited and there is no
definition of what should happen for a number of otherwise valid inputs.

The expectation is that documentation maintainers can use this package with
their client tool to get output which can be inserted in markdown and be
subject to human review.  If you don't have a human in the loop, then there's
absolutely no reason to not just use the HTML renderer: Markdown passes
through HTML, so the HTML output is more portable and safer.

But if your client tooling supports `auto` and can just be asked "hey, give me
the markdown" then documentation maintainers can use that for grabbing
samples.


Coding Style
------------

1. All code should pass `go fmt` and `go vet`
2. `golint` is too opinionated in ways I care about and is explicitly not a
   fixed goal.  Reducing its complaints for new code is worthwhile, unless
   that contradicts something here.  Changing an API to match golint's style
   is never acceptable, unless there's an API major revision bump.
3. Exported constants should be `ALL_CAPS`
4. Imports should be in batched groups, so that `go fmt` will sort within each
   group but not move between groups; those groups should be:
  1. stdlib packages
  2. testing-related packages, if we're a test
  3. intra-repo/org packages
  4. external third-party packages
