package ui


/*
#include <gtk/gtk.h>
#include <stdlib.h>

static inline void free_string(char* s) { free(s); }
*/
// #cgo pkg-config: gtk+-2.0
import "C"
import (
    "os"
    "unsafe"
    "solver"
    "unicode"
    "strconv"
    "strings"
    "github.com/mattn/go-gtk/gtk"
    "github.com/mattn/go-gtk/gdk"
    "github.com/mattn/go-gtk/glib"
)

type (
    cancel_stack_item struct {
        matrix [9][9]uint8
        next *cancel_stack_item
    }
)

var (
    Term chan bool
    s *solver.Solver
    entries [9][9]*gtk.GtkEntry
    examples *gtk.GtkComboBoxText
    prev_y, prev_x uint
    desc_normal, desc_bold *[0]byte
    cancel_stack *cancel_stack_item
    f_size int
)

func cs_push(m [9][9]uint) {
    cs := new(cancel_stack_item)
    for i := uint(0); i < s.Size; i++ {
        for j := uint(0); j < s.Size; j++ {
            cs.matrix[i][j] = uint8(m[i][j])
        }
    }
    cs.next = cancel_stack
    cancel_stack = cs
}

func cs_pop() [9][9]uint {
    if cancel_stack == nil {
        return [9][9]uint{}
    }
    m := [9][9]uint{}
    for i := uint(0); i < s.Size; i++ {
        for j := uint(0); j < s.Size; j++ {
            m[i][j] = uint(cancel_stack.matrix[i][j])
        }
    }
    cancel_stack = cancel_stack.next
    return m
}

func load_sudoku(path string) bool {
    buf := make([]byte, 1024)
    f, err := os.Open(path)
    if err != nil {
        return false
    }
    defer f.Close()

    clear()
    x := uint(0)
    for {
        n, _ := f.Read(buf)
        if n == 0 { break }
        for i := 0; i < n; i++ {
            if x >= s.Size*s.Size { return true }
            if buf[i] >= 49 && buf[i] <= byte(48+s.Size) {
                entries[x/s.Size][x%s.Size].SetText(strconv.Itoa(int(buf[i]-48)))
                modify_font(x/s.Size, x%s.Size, desc_bold)
                x++
            } else if buf[i] == 32 || buf[i] == 46 { // ' ' or '.'
                entries[x/s.Size][x%s.Size].SetText("")
                x++
            }
        }
    }
    return true
}

func clear() {
    for i := uint(0); i < s.Size; i++ {
        for j := uint(0); j < s.Size; j++ {
            entries[i][j].SetText("")
            entries[i][j].SetTooltipText("")
            modify_font(i, j, desc_bold)
        }
    }
    entries[0][0].GrabFocus()
}

func check_field_error(f bool, y1 uint, x1 uint, y2 uint, x2 uint) bool {
    if f { entries[y1][x1].GrabFocus() }
    modify_base(unsafe.Pointer(entries[y2][x2].Widget), gdk.Color("red"))
    return false
}

func check_field(m *[9][9]uint) bool {
    for i := uint(0); i < s.Size; i++ {
        for j := uint(0); j < s.Size; j++ {
            if m[i][j] == 0 { continue }
            flag := true
            b_y, b_x := i/(s.Size/3)*(s.Size/3), j/3*3
            // check row
            for k2 := uint(0); k2 < s.Size; k2++ {
                if k2 != j && m[i][k2] == m[i][j] {
                    flag = check_field_error(flag, i, j, i, k2)
                }
            }
            // check column
            for k1 := uint(0); k1 < s.Size; k1++ {
                if k1 != i && m[k1][j] == m[i][j] {
                    flag = check_field_error(flag, i, j, k1, j)
                }
            }
            // check box
            for k1 := b_y; k1 < b_y+s.Size/3; k1++ {
                for k2 := b_x; k2 < b_x+3; k2++ {
                    if (k1 != i || k2 != j) && m[k1][k2] == m[i][j] {
                        flag = check_field_error(flag, i, j, k1, k2)
                    }
                }
            }
            if !flag { return false }
        }
    }
    return true
}

// stub
func modify_base(v unsafe.Pointer, color *gdk.GdkColor) {
    C.gtk_widget_modify_base((*_Ctype_GtkWidget)(v), C.GtkStateType(gtk.GTK_STATE_NORMAL), (*C.GdkColor)(unsafe.Pointer(&color.Color)))
}

func modify_font(y uint, x uint, desc *[0]byte) {
    C.gtk_widget_modify_font((*_Ctype_GtkWidget)(unsafe.Pointer(entries[y][x].Widget)), desc)
}

func Init(size uint) {
    var (
        files *gtk.GtkHBox
        examples_cnt int
        newfile_flag bool
        cs_desc_normal *_Ctype_char
        cs_desc_bold *_Ctype_char
    )

    s = new(solver.Solver)
    s.Size = size

    if s.Size == 9 {
        cs_desc_normal = C.CString("Sans 14")
        cs_desc_bold = C.CString("Sans Bold 14")
    } else {
        cs_desc_normal = C.CString("Sans 16")
        cs_desc_bold = C.CString("Sans Bold 16")
    }
    desc_normal = C.pango_font_description_from_string(cs_desc_normal)
    desc_bold = C.pango_font_description_from_string(cs_desc_bold)
    C.free_string(cs_desc_normal)
    C.free_string(cs_desc_bold)

    gtk.Init(&os.Args)

    window := gtk.Window(gtk.GTK_WINDOW_TOPLEVEL)
    window.SetResizable(false)
    window.SetTitle("Sudoku solver")
    window.Connect("destroy", func() {
        gtk.MainQuit()
    })
    window.Connect("key-press-event", func(ctx *glib.CallbackContext) {
        arg   := ctx.Args(0)
        kev   := *(**gdk.EventKey)(unsafe.Pointer(&arg))
        r, st := rune(kev.Keyval), gdk.GdkModifierType(kev.State)
        if st & gdk.GDK_CONTROL_MASK != 0 {
            if r == 122 || r == 90 { // Ctrl-Z
                m := cs_pop()
                clear()
                for i := uint(0); i < s.Size; i++ {
                    for j := uint(0); j < s.Size; j++ {
                        v := int(m[i][j])
                        if v != 0 {
                            entries[i][j].SetText(strconv.Itoa(v))
                        }
                    }
                }
            }
        }
    })

    vbox := gtk.VBox(false, 10)

    table := gtk.Table(3, s.Size/3, false)
    bg := [2]*gdk.GdkColor{gdk.Color("white"), gdk.Color("#e9f2ea")}
    for y := uint(0); y < 3; y++ {
        for x := uint(0); x < s.Size/3; x++ {
            subtable := gtk.Table(s.Size/3, 3, false)
            for sy := uint(0); sy < s.Size/3; sy++ {
                for sx := uint(0); sx < 3; sx++ {
                    w := gtk.Entry()
                    w.SetWidthChars(1)
                    w.SetMaxLength(1)
                    if s.Size == 9 {
                        w.SetSizeRequest(23, 25)
                    } else {
                        w.SetSizeRequest(25, 27)
                    }
                    w.Connect("key-press-event", func(ctx *glib.CallbackContext) bool {
                        data := ctx.Data().([]uint)
                        y, x := data[0], data[1]
                        arg := ctx.Args(0)
                        kev := *(**gdk.EventKey)(unsafe.Pointer(&arg))
                        r := rune(kev.Keyval)
                        switch r&0xFF {
                            case 81:
                                if x != 0 || y != 0 {
                                    if x == 0 {
                                        x = s.Size-1
                                        y--
                                    } else {
                                        x--
                                    }
                                }
                            case 82:
                                if y != 0 { y-- }
                            case 83:
                                if x != s.Size-1 || y != s.Size-1 {
                                    if x == s.Size-1 {
                                        x = 0
                                        y++
                                    } else {
                                        x++
                                    }
                                }
                            case 84:
                                if y != s.Size-1 { y++ }
                        }
                        if y != data[0] || x != data[1] {
                            entries[y][x].GrabFocus()
                        }
                        if unicode.IsOneOf([]*unicode.RangeTable{unicode.L, unicode.Z}, r) {
                            return true
                        }
                        return false
                    }, []uint{(s.Size/3)*y+sy, 3*x+sx})
                    w.Connect("grab-focus", func(ctx *glib.CallbackContext) {
                        data := ctx.Data().([]uint)
                        y, x := data[0], data[1]
                        for k := 0; k < 2; k++ {
                            for i := uint(0); i < s.Size; i++ {
                                modify_base(unsafe.Pointer(entries[i][prev_x].Widget), bg[k])
                            }
                            for j := uint(0); j < s.Size; j++ {
                                modify_base(unsafe.Pointer(entries[prev_y][j].Widget), bg[k])
                            }
                            prev_y, prev_x = y, x
                        }
                    }, []uint{(s.Size/3)*y+sy, 3*x+sx})
                    subtable.Attach(w, sx, sx+1, sy, sy+1, gtk.GTK_FILL, gtk.GTK_FILL, 0, 0)
                    entries[(s.Size/3)*y+sy][3*x+sx] = w
                    modify_font((s.Size/3)*y+sy, 3*x+sx, desc_bold)
                }
            }
        table.Attach(subtable, x, x+1, y, y+1, gtk.GTK_FILL, gtk.GTK_FILL, 3, 3)
        }
    }

    solve_btn := gtk.ButtonWithLabel("Solve")
    solve_btn.Clicked(func() {
        var m1, m2 [9][9]uint

        for i := uint(0); i < s.Size; i++ {
            for j := uint(0); j < s.Size; j++ {
                v, _ := strconv.Atoi(entries[i][j].GetText())
                m1[i][j] = uint(v)
            }
        }
        if !check_field(&m1) { return }
        cs_push(m1)
        s.Load(m1)
        s.Solve()
        if s.Finals != s.Size*s.Size {
            // let's try some tough algorithms :)
            s.ToughSolve()
        }
        for i := uint(0); i < s.Size; i++ {
            for j := uint(0); j < s.Size; j++ {
                v := int(s.Get(i, j))
                m2[i][j] = uint(v)
                if v != 0 {
                    entries[i][j].SetText(strconv.Itoa(v))
                } else {
                    var c_string [9]string
                    // get list of candidates
                    c_uint, l := s.GetCandidates(i, j)
                    for k := uint(0); k < l; k++ {
                        c_string[k] = strconv.Itoa(int(c_uint[k]))
                    }
                    // make a tooltip with them
                    entries[i][j].SetTooltipText(strings.Join(c_string[:l], " "))
                }
            }
        }
        // check for differences
        for i := uint(0); i < s.Size; i++ {
            for j := uint(0); j < s.Size; j++ {
                if m1[i][j] == m2[i][j] {
                    modify_font(i, j, desc_bold)
                } else {
                    modify_font(i, j, desc_normal)
                }
            }
        }
    })
    clear_btn := gtk.ButtonWithLabel("Clear")
    clear_btn.Clicked(func() {
        m := [9][9]uint{}
        for i := uint(0); i < s.Size; i++ {
            for j := uint(0); j < s.Size; j++ {
                m[i][j] = s.Get(i, j)
            }
        }
        cs_push(m)
        clear()
    })

    examples = gtk.ComboBoxText()
    // scan `examples` folder
    sz := strconv.Itoa(int(s.Size))
    dir, err := os.Open("examples/"+sz+"x"+sz)
    if err == nil {
        names, err := dir.Readdirnames(0)
        if err == nil {
            for _, v := range names {
                examples.AppendText(v)
                examples_cnt++
            }
        }
        dir.Close()
    }
    examples.Connect("changed", func() {
        sz := strconv.Itoa(int(s.Size))
        load_sudoku("examples/"+sz+"x"+sz+"/"+examples.GetActiveText())
    })

    newfile := gtk.Entry()
    newfile.Connect("activate", func() {
        filename := newfile.GetText()
        if filename != "" {
            sz := strconv.Itoa(int(s.Size))
            f, err := os.Create("examples/"+sz+"x"+sz+"/"+filename)
            if err == nil {
                for i := uint(0); i < s.Size; i++ {
                    for j := uint(0); j < s.Size; j++ {
                        v := []byte(entries[i][j].GetText())
                        if len(v) == 0 || v[0] < 49 || v[0] > byte(48+s.Size) {
                            v = []byte{' '}
                        }
                        f.Write(v)
                        if (j+1) % 3 == 0 && j+1 != s.Size {
                            f.WriteString("*")
                        }
                    }
                    f.WriteString("\n")
                    if (i+1) % (s.Size/3) == 0 && i+1 != s.Size {
                        if s.Size == 9 {
                            f.WriteString("***********\n")
                        } else {
                            f.WriteString("*******\n")
                        }
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
