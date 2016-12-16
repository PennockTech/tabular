// Copyright © 2016 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

/*
The tabular package provides for a model of a two-dimensional grid of cells, or
a "table".  Sub-packages provide for rendering such a table in convenient ways.
For many uses, you only need to import the sub-package.

Sub-packages provide for rendering a TextTable which is designed for a
fixed-cell grid layout of character cells, such as a Unix TTY; for rendering as
an HTML table; for rendering as a CSV.

The core package provides a Table interface; each sub-package embeds a Table in
their own core object, so all methods in the Table interface can be directly
used upon such objects.  These provide for your basic addition of cells, etc.

The API is designed to let you add data quickly, without error-checking every
addition.  Instead, error containers are used, and errors accumulate therein.
The table is an error-container.  A row which is not in a table is also an
error container, but any such errors get moved into the table's collection when
the row is added, and the row thenceforth uses the table's container.


In addition to tables of rows of cells, with virtual columns, we also have
"properties" and "property callbacks".  People using tables with pre-canned
rendering should not need to care about these, but people writing renderers
will need to.

Properties are held by cells, rows, columns and tables.  They are metadata,
designed to be used by sub-packages or anything doing non-trivial table
embedding, to register extra data "about" something.  They're akin to stdlib's
"context", in that the various users can maintained their own namespaced
properties, using a similar API, but can be automatically updated by the table.

Examples might include terminal text properties, color, calculation attributes.

Properties are held by anything which satisfies the PropertyOwner interface.
Such objects then have GetProperty and SetProperty methods.

There is an API balancing act for callers to decide between registering a
property upon the table itself, or wrapping the table: if you want the value to
be automatically updated as a table is mutated, the property approach might
make the most sense, but then with a simple API to present it via a method of
your wrapper object.

A Property Callback is "anything with the UpdateProperties method", which can
be registered with something in the table to be called to automatically update
properties.  These can be registered on a table, row, column or cell; most are
applied to cells, but some are applied to larger objects.  The registration needs
to know what object is to hold the callback, which target the callback is to be
applied to, when it should be invoked, and the new callback being registered.
*/
package tabular // import "go.pennock.tech/tabular"
