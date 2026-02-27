package main

import (
	"github.com/stikypiston/spyglass/lens"
	"github.com/stikypiston/spyglass/lenses/applications"
	"github.com/stikypiston/spyglass/lenses/files"
)

var Lenses = []lens.Lens{
	applications.New(),
	files.New(),
}
