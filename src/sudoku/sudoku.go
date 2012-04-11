package sudoku

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

func (s *Solver) Solve() {
    for {
        delta := false
        for i := 0; i < 9; i++ {
            for j := 0; j < 9; j++ {
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
                for k1 := i/3*3; k1 < i/3*3+3; k1++ {
                    for k2 := j/3*3; k2 < j/3*3+3; k2++ {
                        if s.matrix[k1][k2].final {
                            v &= ^(1<<(s.matrix[k1][k2].value-1))
                        }
                    }
                }
                if v != v_s {
                    delta = true
                }
                if count(v) == 1 {
                    v, f = fastlog2(v), true
                }
                s.matrix[i][j].value, s.matrix[i][j].final = v, f
            }
        }
        if !delta {
            break
        }
    }
}
