package main

import (
	"github.com/stikypiston/spyglass/lens"
	"github.com/stikypiston/spyglass/lenses/applications"

	"github.com/stikypiston/spyglass/lenses/nerdfont"
	"github.com/stikypiston/spyglass/lenses/power"

	"github.com/stikypiston/spyglass/lenses/files"
)

var Lenses = []lens.Lens{
	applications.New(),
	power.New(),
	nerdfont.New(),
	files.New(),
}
