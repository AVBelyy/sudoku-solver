package solver

type m_item struct {
    value uint
    final bool
}

type point struct {
    y uint
    x uint
}

type chain struct {
    path [81]point
    length uint
}

type Solver struct {
    matrix [9][9]m_item
    Size uint
    Finals uint
}

var (
    DiabolicDelta bool
    XYChains_flag bool
)

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

func (ch *chain) in(p point) bool {
    for i := uint(0); i < ch.length; i++ {
        if ch.path[i] == p {
            return true
        }
    }
    return false
}

func (ch *chain) add(p point) {
    ch.path[ch.length] = p
    ch.length += 1
}

func (s *Solver) Load(src [9][9]uint) {
    s.Finals = 0
    // transform input matrix to convenient format
    for i := uint(0); i < s.Size; i++ {
        for j := uint(0); j < s.Size; j++ {
            v, f := src[i][j], true
            switch {
                case v == 0:
                    v = (1<<s.Size)-1
                    f = false
                case 1 <= v && v <= s.Size:
                    s.Finals += 1
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
    return res, cnt
}

func (s *Solver) Solve() {
    var pair [9]uint
    
    for {
        delta := false
        for i := uint(0); i < s.Size; i++ {
            for j := uint(0); j < s.Size; j++ {
                b_y, b_x := i/(s.Size/3)*(s.Size/3), j/3*3
                v, f := s.matrix[i][j].value, s.matrix[i][j].final
                if f {
                    continue
                }
                v_s := v
                // check cell's row
                for k := uint(0); k < s.Size; k++ {
                    if s.matrix[i][k].final {
                        v &= ^(1<<(s.matrix[i][k].value-1))
                    }
                }
                // check cell's column
                for k := uint(0); k < s.Size; k++ {
                    if s.matrix[k][j].final {
                        v &= ^(1<<(s.matrix[k][j].value-1))
                    }
                }
                // check cell's box
                for k1 := b_y; k1 < b_y+s.Size/3; k1++ {
                    for k2 := b_x; k2 < b_x+3; k2++ {
                        if s.matrix[k1][k2].final {
                            v &= ^(1<<(s.matrix[k1][k2].value-1))
                        }
                    }
                }
                // check for hidden singles in row
                for x := uint(0); x < s.Size; x++ {
                    if v&(1<<x) == 0 { continue }
                    for k := uint(0); k < s.Size; k++ {
                        if !s.matrix[i][k].final && k != j {
                            if s.matrix[i][k].value&(1<<x) != 0 {
                                goto hidden_singles_row_next
                            }
                        }
                    }
                    v, f, delta = x+1, true, true
                    s.Finals += 1
                    goto final
                    hidden_singles_row_next:
                }
                // check for hidden singles in column
                for x := uint(0); x < s.Size; x++ {
                    if v&(1<<x) == 0 { continue }
                    for k := uint(0); k < s.Size; k++ {
                        if !s.matrix[k][j].final && k != i {
                            if s.matrix[k][j].value&(1<<x) != 0 {
                                goto hidden_singles_column_next
                            }
                        }
                    }
                    v, f, delta = x+1, true, true
                    s.Finals += 1
                    goto final
                    hidden_singles_column_next:
                }
                // check for hidden singles in box
                for x := uint(0); x < s.Size; x++ {
                    if v&(1<<x) == 0 { continue }
                    for k1 := b_y; k1 < b_y+s.Size/3; k1++ {
                        for k2 := b_x; k2 < b_x+3; k2++ {
                            if !s.matrix[k1][k2].final && (k1 != i || k2 != j) {
                                if s.matrix[k1][k2].value&(1<<x) != 0 {
                                    goto hidden_singles_box_next
                                }
                            }
                        }
                    }
                    v, f, delta = x+1, true, true
                    s.Finals += 1
                    goto final
                    hidden_singles_box_next:
                }
                // check for hidden pairs in row
                for x := uint(0); x < s.Size; x++ {
                    occurs, where := 0, uint(0)
                    pair[x] = 10
                    if v&(1<<x) == 0 { continue }
                    for k := uint(0); k < s.Size; k++ {
                        if !s.matrix[i][k].final && k != j {
                            if s.matrix[i][k].value&(1<<x) != 0 {
                                occurs, where = occurs+1, k
                            }
                        }
                    }
                    if occurs == 1 {
                        pair[x] = where
                    }
                }
                for k1 := uint(0);  k1 < s.Size; k1++ {
                    if pair[k1] == 10 { continue }
                    for k2 := k1+1; k2 < s.Size; k2++ {
                        if pair[k1] == pair[k2] {
                            v = (1<<k1)|(1<<k2)
                            s.matrix[i][pair[k1]].value = v
                            break
                        }
                    }
                }
                // check for hidden pairs in column
                for x := uint(0); x < s.Size; x++ {
                    occurs, where := 0, uint(0)
                    pair[x] = 10
                    if v&(1<<x) == 0 { continue }
                    for k := uint(0); k < s.Size; k++ {
                        if !s.matrix[k][j].final && k != i {
                            if s.matrix[k][j].value&(1<<x) != 0 {
                                occurs, where = occurs+1, k
                            }
                        }
                    }
                    if occurs == 1 {
                        pair[x] = where
                    }
                }
                for k1 := uint(0);  k1 < s.Size; k1++ {
                    if pair[k1] == 10 { continue }
                    for k2 := k1+1; k2 < s.Size; k2++ {
                        if pair[k1] == pair[k2] {
                            v = (1<<k1)|(1<<k2)
                            s.matrix[pair[k1]][j].value = v
                            break
                        }
                    }
                }
                // check for hidden pairs in box
                for x := uint(0); x < s.Size; x++ {
                    occurs, where := 0, uint(0)
                    pair[x] = 10
                    if v&(1<<x) == 0 { continue }
                    for k1 := b_y; k1 < b_y+s.Size/3; k1++ {
                        for k2 := b_x; k2 < b_x+3; k2++ {
                            if !s.matrix[k1][k2].final && (k1 != i || k2 != j) {
                                if s.matrix[k1][k2].value&(1<<x) != 0 {
                                    occurs, where = occurs+1, k1<<4|k2
                                }
                            }
                        }
                    }
                    if occurs == 1 {
                        pair[x] = where
                    }
                }
                for k1 := uint(0);  k1 < s.Size; k1++ {
                    if pair[k1] == 10 { continue }
                    for k2 := k1+1; k2 < s.Size; k2++ {
                        if pair[k1] == pair[k2] {
                            v = (1<<k1)|(1<<k2)
                            s.matrix[pair[k1]>>4][pair[k1]&0xF].value = v
                            break
                        }
                    }
                }
                // check for naked pairs in row
                for k := uint(0); k < s.Size; k++ {
                    if v == s.matrix[i][k].value && !s.matrix[i][k].final && k != j && count(v) == 2 {
                        for x := uint(0); x < s.Size; x++ {
                            if v&(1<<x) == 0 { continue }
                            for k_ := uint(0); k_ < s.Size; k_++ {
                                if k_ != k && k_ != j && !s.matrix[i][k_].final {
                                    s.matrix[i][k_].value &= ^(1<<x)
                                }
                            }
                        }
                    }
                }
                // check for naked pairs in column
                for k := uint(0); k < s.Size; k++ {
                    if v == s.matrix[k][j].value && !s.matrix[k][j].final && k != i && count(v) == 2 {
                        for x := uint(0); x < s.Size; x++ {
                            if v&(1<<x) == 0 { continue }
                            for k_ := uint(0); k_ < s.Size; k_++ {
                                if k_ != k && k_ != i && !s.matrix[k_][j].final {
                                    s.matrix[k_][j].value &= ^(1<<x)
                                }
                            }
                        }
                    }
                }
                // check for naked pairs in box
                for k1 := b_y; k1 < b_y+s.Size/3; k1++ {
                    for k2 := b_x; k2 < b_x+3; k2++ {
                        if v == s.matrix[k1][k2].value && !s.matrix[k1][k2].final && (k1 != i || k2 != j) && count(v) == 2 {
                            for x := uint(0); x < s.Size; x++ {
                                if v&(1<<x) == 0 { continue }
                                for k1_ := b_y; k1_ < b_y+s.Size/3; k1_++ {
                                    for k2_ := b_x; k2_ < b_x+3; k2_++ {
                                        if (k1_ != k1 || k2_ != k2) && (k1_ != i || k2_ != j) && !s.matrix[k1_][k2_].final {
                                            s.matrix[k1_][k2_].value &= ^(1<<x)
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
                    s.Finals += 1
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

func (s *Solver) XYChains_links(p1 point, p2 point, common uint) bool {
    return count(s.matrix[p1.y][p1.x].value & s.matrix[p2.y][p2.x].value) == common
}

func (s *Solver) XYChains_recur(ch chain, p point) {
    if XYChains_flag { return }
    
    b_y, b_x := p.y/(s.Size/3)*(s.Size/3), p.x/3*3
    // now recursively find the rest of path

    defer func() {
        if ch.length >= 3 && s.XYChains_links(ch.path[0], ch.path[ch.length-1], 1) {
            p1, p2 := ch.path[0], ch.path[ch.length-1]
            mask := [2][9][9]bool{}
            // calculate number of bit that will be unset
            x := fastlog2(s.matrix[p1.y][p1.x].value & s.matrix[p2.y][p2.x].value)-1
            for i, v := range []point{p1, p2} {
                _b_y, _b_x := v.y/(s.Size/3)*(s.Size/3), v.x/3*3
                // row
                for k := uint(0); k < s.Size; k++ {
                    if k != v.x {
                        mask[i][v.y][k] = true
                    }
                }
                // column
                for k := uint(0); k < s.Size; k++ {
                    if k != v.y {
                        mask[i][k][v.x] = true
                    }
                }
                // box
                for k1 := _b_y; k1 < _b_y+3; k1++ {
                    for k2 := _b_x; k2 < _b_x+3; k2++ {
                        if k1 != v.y || k2 != v.x {
                            mask[i][k1][k2] = true
                        }
                    }
                }
            }
            for i := uint(0); i < s.Size; i++ {
                for j := uint(0); j < s.Size; j++ {
                    if mask[0][i][j] && mask[1][i][j] && !s.matrix[i][j].final {
                        s.matrix[i][j].value &= ^(1<<x)
                        DiabolicDelta = true
                    }
                }
            }
            XYChains_flag = true
        }
    }()

    // scan row
    for k := uint(0); k < s.Size; k++ {
        if k != p.x && !s.matrix[p.y][k].final && count(s.matrix[p.y][k].value) == 2 {
            p2 := point{p.y, k}
            if !ch.in(p2) && s.XYChains_links(p, p2, 1) && (ch.length == 1 || s.XYChains_links(ch.path[ch.length-2], p2, 0)) {
                ch_new := ch
                ch_new.add(p2)
                s.XYChains_recur(ch_new, p2)
            }
        }
    }
    // scan column
    for k := uint(0); k < s.Size; k++ {
        if k != p.y && !s.matrix[k][p.x].final && count(s.matrix[k][p.x].value) == 2 {
            p2 := point{k, p.x}
            if !ch.in(p2) && s.XYChains_links(p, p2, 1) && (ch.length == 1 || s.XYChains_links(ch.path[ch.length-2], p2, 0)) {
                ch_new := ch
                ch_new.add(p2)
                s.XYChains_recur(ch_new, p2)
            }
        }
    }
    // scan box
    for k1 := b_y; k1 < b_y+s.Size/3; k1++ {
        for k2 := b_x; k2 < b_x+3; k2++ {
            if (k1 != p.y || k2 != p.x) && !s.matrix[k1][k2].final && count(s.matrix[k1][k2].value) == 2 {
                p2 := point{k1, k2}
                if !ch.in(p2) && s.XYChains_links(p, p2, 1) && (ch.length == 1 || s.XYChains_links(ch.path[ch.length-2], p2, 0)) {
                    ch_new := ch
                    ch_new.add(p2)
                    s.XYChains_recur(ch_new, p2)
                }
            }
        }
    }
}

func (s *Solver) XYChains() {
    XYChains_flag = false
    // let's find a start point - an unsolved cell with exactly 2 candidates
   
    for i := uint(0); i < s.Size; i++ {
        for j := uint(0); j < s.Size; j++ {
            if !s.matrix[i][j].final && count(s.matrix[i][j].value) == 2 {
                ch := chain{}
                ch.add(point{i, j})
                s.XYChains_recur(ch, point{i, j})
            }
        }
    }
}

func (s *Solver) DiabolicSolve() {
    DiabolicDelta = false
    
    // more complicated and slower algorithms are presented here
    s.XYChains()
    if DiabolicDelta {
        s.Solve()
    }
}
