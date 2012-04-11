package main

import (
    "fmt"
    "sudoku"
)

func printMatrix (s *sudoku.Solver) {
    var i, j uint
    for i = 0; i < 9; i++ {
        for j = 0; j < 9; j++ {
            v := s.Get(i, j)
            if v == 0 {
                fmt.Printf("  ")
            } else {
                fmt.Printf("%d ", v)
            }
            if j % 3 == 2 && j != 8 {
                fmt.Printf("* ")
            }
        }
        fmt.Printf("\n")
        if i % 3 == 2 && i != 8 {
            fmt.Printf("*********************\n")
        }
    }
}

func main() {
    input := [9][9]uint{{0, 0, 0,/**/1, 0, 5,/**/0, 0, 0},
                        {1, 4, 0,/**/0, 0, 0,/**/6, 7, 0},
                        {0, 8, 0,/**/0, 0, 2,/**/4, 0, 0},
                        /*******************************/
                        /*******************************/
                        {0, 6, 3,/**/0, 7, 0,/**/0, 1, 0},
                        {9, 0, 0,/**/0, 0, 0,/**/0, 0, 3},
                        {0, 1, 0,/**/0, 9, 0,/**/5, 2, 0},
                        /*******************************/
                        /*******************************/
                        {0, 0, 7,/**/2, 0, 0,/**/0, 8, 0},
                        {0, 2, 6,/**/0, 0, 0,/**/0, 3, 5},
                        {0, 0, 0,/**/4, 0, 9,/**/0, 0, 0}}

    s := new(sudoku.Solver)
    s.Load(input)
    fmt.Printf("The input sudoku matrix:\n")
    // out the input matrix
    printMatrix(s)
    s.Solve()
    fmt.Printf("\n\nThe output sudoku matrix:\n")
    // out the output matrix
    printMatrix(s)
}
