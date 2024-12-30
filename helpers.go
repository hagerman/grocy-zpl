package main

import "github.com/OpenPrinting/goipp"

func makeAttrCollection(name string,
	member1 goipp.Attribute, members ...goipp.Attribute) goipp.Attribute {

	col := make(goipp.Collection, len(members)+1)
	col[0] = member1
	copy(col[1:], members)

	return goipp.MakeAttribute(name, goipp.TagBeginCollection, col)
}
