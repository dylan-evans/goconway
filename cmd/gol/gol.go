package main

import (
	"flag"
	"fmt"
)

type Stat struct {
	Alive, Dead, Changed uint64
}

func (stat *Stat) Add(source Stat) {
	stat.Alive += source.Alive
	stat.Dead += source.Dead
	stat.Changed += source.Changed
}

type Cell struct {
	value uint8
}

func main() {
	rowsPtr := flag.Uint64("rows", 100, "rows")
	colsPtr := flag.Uint64("columns", 100, "columns")
    flag.Parse()
	mesh := NewMesh(*rowsPtr, *colsPtr)
    fmt.Println("Adding chaos")
	mesh.Chaos()
	fmt.Println("Starting", mesh.Rows, "rows", mesh.Columns, "columns")
	for i := 0; ; i++ {
		stat := mesh.Update()
		mesh.Swap()
		fmt.Println("loop", i, stat)
	}
}
