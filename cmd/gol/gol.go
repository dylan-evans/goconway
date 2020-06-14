package main

import (
    "fmt"
    "flag"
    "math/rand"
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

type Mesh struct {
    Rows, Columns uint64
    current, staging []Cell
}

func NewMesh(rows, columns uint64) Mesh {
    return Mesh{ Rows: rows,
        Columns: columns,
        current: make([]Cell, rows * columns, rows * columns),
        staging: make([]Cell, rows * columns, rows * columns),
    }
}

func (mesh *Mesh) GetValue(row, column uint64) uint8 {
    return mesh.current[row * mesh.Columns + column].value
}

func (mesh *Mesh) SetValue(row, column uint64, value uint8) {
    mesh.staging[row * mesh.Columns + column].value = value
}

func (mesh *Mesh) UpdateRow(row uint64) Stat {
    var stat Stat
    for col := uint64(0); col < mesh.Columns; col++ {
        count := mesh.Counter(row, col)
        if mesh.GetValue(row, col) > 0 {
            if count < 2 || count > 3 {
                mesh.SetValue(row, col, 0)
                stat.Changed++
                stat.Dead++
            } else {
                stat.Alive++
            }
        } else {
            if count == 3 {
                mesh.SetValue(row, col, 1)
                stat.Changed++
                stat.Alive++
            } else {
                stat.Dead++
            }
        }
    }
    return stat
}

func (mesh *Mesh) Swap() {
    mesh.current, mesh.staging = mesh.staging, mesh.current
}

func (mesh *Mesh) Update() Stat {
    var total Stat
    stat_ch := make(chan Stat, mesh.Rows * mesh.Columns)
    for row := uint64(0); row < mesh.Rows; row++ {
        go func(irow uint64) {
            stat_ch <- mesh.UpdateRow(irow)
        }(row)
    }
    for i := uint64(0); i < mesh.Rows; i++ {
        total.Add(<- stat_ch)
    }
    return total
}


func (mesh *Mesh) Counter(row uint64, column uint64) int {
    var total int
    for _, rmod := range []int{-1, 0, 1} {
        for _, cmod := range []int{-1, 0, 1} {
            if cmod != 0 || rmod != 0 {
                n_row := int(row) + rmod
                n_col := int(column) + cmod
                if n_row > 0 && n_col > 0 && n_row < int(mesh.Rows) && n_col < int(mesh.Columns) {
                    total += int(mesh.GetValue(uint64(n_row), uint64(n_col)))
                }
            }
        }
    }
    return total
}

func (mesh *Mesh) Chaos() {
    for i := uint64(0); i < mesh.Rows * mesh.Columns; i++ {
        mesh.staging[i].value = uint8(rand.Uint64() & 2)
    }
}


func main() {
    rowsPtr := flag.Uint64("rows", 100, "rows")
    colsPtr := flag.Uint64("columns", 100, "columns")
    mesh := NewMesh(*rowsPtr, *colsPtr)
    mesh.SetValue(10, 10, 1)
    fmt.Println("Generating Chaos")
    mesh.Chaos()
    mesh.Swap()
    fmt.Println("Starting")
    for i := 0;; i++ {
        stat := mesh.Update()
        mesh.Swap()
        fmt.Println("loop", i, stat)
    }
}
