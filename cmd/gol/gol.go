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
	mesh := NewMesh(*rowsPtr, *colsPtr)
	mesh.SetValue(10, 10, 1)
	mesh.Chaos()
	mesh.Swap()
	fmt.Println("Starting")
	for i := 0; ; i++ {
		stat := mesh.Update()
		mesh.Swap()
		fmt.Println("loop", i, stat)
	}
}
