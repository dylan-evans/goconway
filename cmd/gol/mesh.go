package main

import (
	"math/rand"
)

const BITS = 64

type Cell struct {
	value uint8
}

type Mesh struct {
	Rows, Columns    uint64
	current, staging []uint64
}

func NewMesh(rows, columns uint64) Mesh {
	mesh := Mesh{Rows: rows,
		Columns: columns,
	}
	mesh.current = make([]uint64, mesh.Size(), mesh.Size())
	mesh.staging = make([]uint64, mesh.Size(), mesh.Size())
	return mesh
}

func (mesh *Mesh) Size() uint64 {
	return mesh.Rows * mesh.Columns
}

func (mesh *Mesh) getOffset(row, column uint64) uint64 {
	return (row*mesh.Columns + column) / BITS
}

func (mesh *Mesh) getMask(row, column uint64) (uint64, uint64) {
    offset := row*mesh.Columns + column
    return offset / BITS, offset % BITS
}

func (mesh *Mesh) GetValue(row, column uint64) bool {
    offset, mask := mesh.getMask(row, column)
    return bool(mesh.current[offset] & mask > 0)
}

func (mesh *Mesh) SetValue(row, column uint64, value bool) {
    offset, mask := mesh.getMask(row, column)
    if value {
        mesh.staging[offset] |= (1 << mask)
    } else {
        mesh.staging[offset] &^= (1 << mask)
    }
}

func (mesh *Mesh) UpdateRow(row uint64) Stat {
	var stat Stat
	for col := uint64(0); col < mesh.Columns; col++ {
		count := mesh.Counter(row, col)
		if mesh.GetValue(row, col) {
			if count < 2 || count > 3 {
				mesh.SetValue(row, col, false)
				stat.Changed++
				stat.Dead++
			} else {
				mesh.SetValue(row, col, true)
				stat.Alive++
			}
		} else {
			if count == 3 {
				mesh.SetValue(row, col, true)
				stat.Changed++
				stat.Alive++
			} else {
				mesh.SetValue(row, col, false)
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
	total := Stat{}
	stat_ch := make(chan Stat, mesh.Rows*mesh.Columns)
	for row := uint64(0); row < mesh.Rows; row++ {
		go func(irow uint64) {
			stat_ch <- mesh.UpdateRow(irow)
		}(row)
	}
	for i := uint64(0); i < mesh.Rows; i++ {
		total.Add(<-stat_ch)
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
					if mesh.GetValue(uint64(n_row), uint64(n_col)) {
                        total++
                    }
				}
			}
		}
	}
	return total
}

func (mesh *Mesh) Chaos() {
	for i := uint64(0); int(i) < len(mesh.current); i++ {
		mesh.current[i] = rand.Uint64()
	}
}
