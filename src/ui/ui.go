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
    var (
        files *gtk.GtkHBox
        examples_cnt int
        newfile_flag bool
    )

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
                examples_cnt++
            }
        }
    }
    dir.Close()
    examples.Connect("changed", func() {
        load_sudoku("examples/"+examples.GetActiveText())
    })

    newfile := gtk.Entry()
    newfile.Connect("activate", func() {
        filename := newfile.GetText()
        if filename != "" {
            f, err := os.Create("examples/"+filename)
            if err == nil {
                for i := 0; i < 9; i++ {
                    for j := 0; j < 9; j++ {
                        v := []byte(entries[i][j].GetText())
                        if len(v) == 0 || v[0] < 49 || v[0] > 57 {
                            v = []byte{' '}
                        }
                        f.Write(v)
                        if j == 2 || j == 5 {
                            f.WriteString("*")
                        }
                    }
                    f.WriteString("\n")
                    if i == 2 || i == 5 {
                        f.WriteString("***********\n")
                    }
                }
                f.Close()
            }
            examples.AppendText(filename)
            examples.SetActive(examples_cnt)
            examples_cnt++
        }
        files.ShowAll()
        newfile.SetText("")
        newfile.Hide()
    })

    icon := gtk.Image()
    icon.SetFromStock(gtk.GTK_STOCK_SAVE_AS, gtk.GTK_ICON_SIZE_BUTTON)
    export := gtk.Button()
    export.SetImage(icon)
    export.Clicked(func() {
        if !newfile_flag {
            files.Add(newfile)
            newfile_flag = true
        }
        files.ShowAll()
        examples.Hide()
        export.Hide()
        newfile.GrabFocus()
    })

    files = gtk.HBox(false, 0)
    files.Add(export)
    files.Add(examples)

    buttons := gtk.HBox(true, 5)
    buttons.Add(solve_btn)
    buttons.Add(clear_btn)

    vbox.Add(table)
    vbox.Add(files)
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
