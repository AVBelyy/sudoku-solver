package ui

import (
    "os"
    "unsafe"
    "solver"
    "unicode"
    "strconv"
    "github.com/mattn/go-gtk/gtk"
    "github.com/mattn/go-gtk/gdk"
    "github.com/mattn/go-gtk/glib"
)

var (
    Term chan bool
    s *solver.Solver
    entries [9][9]*gtk.GtkEntry
)

func load_sudoku(path string) bool {
    buf := make([]byte, 1024)
    f, err := os.Open(path)
    if err != nil {
        return false
    }
    defer f.Close()

    clear()
    x := 0
    for {
        n, _ := f.Read(buf)
        if n == 0 || x == 81 { break }
        for i := 0; i < n; i++ {
            if buf[i] >= 49 && buf[i] <= 57 {
                entries[x/9][x%9].SetText(strconv.Itoa(int(buf[i]-48)))
                x++
            } else if buf[i] == 32 {
                entries[x/9][x%9].SetText("")
                x++
            }
        }
    }
    return true
}

func clear() {
    for i := 0; i < 9; i++ {
        for j := 0; j < 9; j++ {
            entries[i][j].SetText("")
        }
    }
}

func Init() {
    s = new(solver.Solver)

    gtk.Init(&os.Args)

    window := gtk.Window(gtk.GTK_WINDOW_TOPLEVEL)
    window.SetResizable(false)
    window.SetTitle("Sudoku solver")
    window.Connect("destroy", func() {
        gtk.MainQuit()
    })

    vbox := gtk.VBox(false, 10)

    table := gtk.Table(3, 3, false)
    for y := uint(0); y < 3; y++ {
        for x := uint(0); x < 3; x++ {
            subtable := gtk.Table(3, 3, false)
            for sy := uint(0); sy < 3; sy++ {
                for sx := uint(0); sx < 3; sx++ {
                    w := gtk.Entry()
                    w.SetWidthChars(1)
                    w.SetMaxLength(1)
                    w.Connect("key-press-event", func(ctx *glib.CallbackContext) bool {
                        data := ctx.Data().([]uint)
                        y, x := data[0], (data[1]+1)%9
                        arg := ctx.Args(0)
                        kev := *(**gdk.EventKey)(unsafe.Pointer(&arg))
                        r := rune(kev.Keyval)
                        if unicode.IsLetter(r) {
                            return true
                        }
                        if x == 0 {
                            y++
                        }
                        if y == 9 {
                            return false
                        }
                        if unicode.IsNumber(r) || r&0xFF == 9 {
                            entries[y][x].GrabFocus()
                        }
                        return r&0xFF == 9
                    }, []uint{3*y+sy, 3*x+sx})
                    subtable.Attach(w, sx, sx+1, sy, sy+1, gtk.GTK_FILL, gtk.GTK_FILL, 0, 0)
                    entries[3*y+sy][3*x+sx] = w
                }
            }
        table.Attach(subtable, x, x+1, y, y+1, gtk.GTK_FILL, gtk.GTK_FILL, 3, 3)
        }
    }

    solve_btn := gtk.ButtonWithLabel("Solve")
    solve_btn.Clicked(func() {
        var matrix [9][9]uint

        for i := 0; i < 9; i++ {
            for j := 0; j < 9; j++ {
                v, _ := strconv.Atoi(entries[i][j].GetText())
                matrix[i][j] = uint(v)
            }
        }
        s.Load(matrix)
        s.Solve()
        for i := uint(0); i < 9; i++ {
            for j := uint(0); j < 9; j++ {
                v := int(s.Get(i, j))
                if v != 0 {
                    entries[i][j].SetText(strconv.Itoa(v))
                }
            }
        }
    })
    clear_btn := gtk.ButtonWithLabel("Clear")
    clear_btn.Clicked(clear)

    examples := gtk.ComboBoxNewText()
    // scan `examples` folder
    dir, err := os.Open("examples")
    if err == nil {
        names, err := dir.Readdirnames(0)
        if err == nil {
            for _, v := range names {
                examples.AppendText(v)
            }
        }
    }
    dir.Close()
    examples.Connect("changed", func() {
        load_sudoku("examples/"+examples.GetActiveText())
    })

    buttons := gtk.HBox(true, 5)
    buttons.Add(solve_btn)
    buttons.Add(clear_btn)

    vbox.Add(table)
    vbox.Add(examples)
    vbox.Add(buttons)

    window.Add(vbox)
    window.ShowAll()
}

func Event_loop() {
    defer func() {
        Term <- true
    }()
    gtk.Main()
}
