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

func Init() {
    s = new(solver.Solver)

    gtk.Init(&os.Args)

    window := gtk.Window(gtk.GTK_WINDOW_TOPLEVEL)
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
    clear_btn := gtk.ButtonWithLabel("Cancel")

    buttons := gtk.HBox(true, 5)
    buttons.Add(solve_btn)
    buttons.Add(clear_btn)

    vbox.Add(table)
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
