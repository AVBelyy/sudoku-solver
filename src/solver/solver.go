package solver

type m_item struct {
    value uint
    final bool
}

type Solver struct {
    matrix [9][9]m_item
}

func count(bm uint) uint {
    var cnt uint
    for ; bm != 0; cnt++ {
        bm &= bm-1
    }
    return cnt
}

func fastlog2(bm uint) uint {
    var log uint
    for ; bm != 0; log++ {
        bm >>= 1
    }
    return log
}

func (s *Solver) Load(src [9][9]uint) {
    // transform input matrix to convenient format
    for i := 0; i < 9; i++ {
        for j := 0; j < 9; j++ {
            v, f := src[i][j], true
            switch {
                case v == 0:
                    v = (1<<9)-1
                    f = false
                case 1 <= v && v <= 9:
                    // do nothing
                default:
                    // TODO: throw an exception: wrong src matrix format
            }
            s.matrix[i][j].value, s.matrix[i][j].final = v, f
        }
    }
}

func (s *Solver) Get(y, x uint) uint {
    if s.matrix[y][x].final {
        return s.matrix[y][x].value
    }
    return 0
}

func (s *Solver) GetCandidates(y, x uint) ([9]uint, uint) {
    var (
        i, cnt uint
        res [9]uint
    )
    if s.matrix[y][x].final {
        res[0] = s.matrix[y][x].value
    } else {
        for bm := s.matrix[y][x].value; bm != 0; i++ {
            if bm & 1 != 0 {
                res[cnt] = i+1
                cnt += 1
            }
            bm >>= 1
        }
    }
    return res, cnt+1
}

func (s *Solver) Solve() {
    for {
        delta := false
        for i := uint(0); i < 9; i++ {
            for j := uint(0); j < 9; j++ {
                b_y, b_x := i/3*3, j/3*3
                v, f := s.matrix[i][j].value, s.matrix[i][j].final
                if f {
                    continue
                }
                v_s := v
                // check cell's row
                for k := 0; k < 9; k++ {
                    if s.matrix[i][k].final {
                        v &= ^(1<<(s.matrix[i][k].value-1))
                    }
                }
                // check cell's column
                for k := 0; k < 9; k++ {
                    if s.matrix[k][j].final {
                        v &= ^(1<<(s.matrix[k][j].value-1))
                    }
                }
                // check cell's box
                for k1 := b_y; k1 < b_y+3; k1++ {
                    for k2 := b_x; k2 < b_x+3; k2++ {
                        if s.matrix[k1][k2].final {
                            v &= ^(1<<(s.matrix[k1][k2].value-1))
                        }
                    }
                }
                // check for hidden singles in row
                for x := uint(0); x < 9; x++ {
                    if v&(1<<x) == 0 { continue }
                    for k := uint(0); k < 9; k++ {
                        if !s.matrix[i][k].final && k != j {
                            if s.matrix[i][k].value&(1<<x) != 0 {
                                goto hidden_singles_row_next
                            }
                        }
                    }
                    v, f, delta = x+1, true, true
                    goto final
                    hidden_singles_row_next:
                }
                // check for hidden singles in column
                for x := uint(0); x < 9; x++ {
                    if v&(1<<x) == 0 { continue }
                    for k := uint(0); k < 9; k++ {
                        if !s.matrix[k][j].final && k != i {
                            if s.matrix[k][j].value&(1<<x) != 0 {
                                goto hidden_singles_column_next
                            }
                        }
                    }
                    v, f, delta = x+1, true, true
                    goto final
                    hidden_singles_column_next:
                }
                // check for hidden singles in box
                for x := uint(0); x < 9; x++ {
                    if v&(1<<x) == 0 { continue }
                    for k1 := b_y; k1 < b_y+3; k1++ {
                        for k2 := b_x; k2 < b_x+3; k2++ {
                            if !s.matrix[k1][k2].final && (k1 != i || k2 != j) {
                                if s.matrix[k1][k2].value&(1<<x) != 0 {
                                    goto hidden_singles_box_next
                                }
                            }
                        }
                    }
                    v, f, delta = x+1, true, true
                    goto final
                    hidden_singles_box_next:
                }
                // check for naked pairs in column
                for k := uint(0); k < 9; k++ {
                    if v == s.matrix[k][j].value && !s.matrix[k][j].final && k != i && count(v) == 2 {
                        for x := uint(0); x < 9; x++ {
                            if v&(1<<x) == 0 { continue }
                            for k_ := uint(0); k_ < 9; k_++ {
                                if k_ != k && k_ != i && !s.matrix[k_][j].final {
                                    s.matrix[k_][j].value &= ^(1<<x)
                                    delta = true
                                }
                            }
                        }
                    }
                }
                // check for naked pairs in box
                for k1 := b_y; k1 < b_y+3; k1++ {
                    for k2 := b_x; k2 < b_x+3; k2++ {
                        if v == s.matrix[k1][k2].value && !s.matrix[k1][k2].final && (k1 != i || k2 != j) && count(v) == 2 {
                            for x := uint(0); x < 9; x++ {
                                if v&(1<<x) == 0 { continue }
                                for k1_ := b_y; k1_ < b_y+3; k1_++ {
                                    for k2_ := b_x; k2_ < b_x+3; k2_++ {
                                        if (k1_ != k1 || k2_ != k2) && (k1_ != i || k2_ != j) && !s.matrix[k1_][k2_].final {
                                            s.matrix[k1_][k2_].value &= ^(1<<x)
                                            delta = true
                                        }
                                    }
                                }
                            }
                        }
                    }
                }
                if v != v_s {
                    delta = true
                }
                if count(v) == 1 {
                    v, f = fastlog2(v), true
                }
                final:
                s.matrix[i][j] = m_item{v, f}
            }
        }
        if !delta {
            break
        }
    }
}
