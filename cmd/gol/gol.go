package main

import (
	"flag"
	"fmt"
    "time"
)

type Stat struct {
	Alive, Dead, Changed uint64
}

func (stat *Stat) Add(source Stat) {
	stat.Alive += source.Alive
	stat.Dead += source.Dead
	stat.Changed += source.Changed
}

func main() {
	rowsPtr := flag.Uint64("rows", 100, "rows")
	colsPtr := flag.Uint64("columns", 100, "columns")
	flag.Parse()
	mesh := NewMesh(*rowsPtr, *colsPtr)
	fmt.Println("Adding chaos")
	mesh.Chaos()
	fmt.Println("Starting", mesh.Rows, "rows", mesh.Columns, "columns")

    var stat Stat
    var loop int

    go func() {
        // The monitor loop
        delay, _ := time.ParseDuration("1s")
        for {
            fmt.Println("loop", loop, "alive", stat.Alive, "dead", stat.Dead, "changed", stat.Changed)
            time.Sleep(delay)
        }
    }()

	for ; ; loop++ {
		stat = mesh.Update()
		mesh.Swap()
    }
}
