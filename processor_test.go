package whitefriday

import (
    "testing"
    "runtime"
    "fmt"
    "path/filepath"
    "reflect"
)

func TestSimple(t *testing.T) {
}

func TestBold(t *testing.T) {
    equals(t, "**Derpy** **Herpy**", Convert("<b>Derpy</b> <strong>Herpy</strong>"))
}

func TestItalic(t *testing.T) {
    equals(t, "*Derpy* *Herpy*", Convert("<i>Derpy</i> <em>Herpy</em>"))
}

func TestBoldItalic(t *testing.T) {
    equals(t, "*I am **Derpy**, bitch!*", Convert("<i>I am <b>Derpy</b>, bitch!</i>"))
}

func TestUnstyledList(t *testing.T) {
    equals(t, "* Wow\n* Dat **az**\n* Hey\nYou", Convert("<ul><li>Wow</li><li>Dat <b>az</b></li><li>Hey\nYou</li></ul>"))
}

func TestOrderedList(t *testing.T) {
    equals(t, "1. Wow\n1. Dat **az**\n1. Hey\nYou", Convert("<ol><li>Wow</li><li>Dat <b>az</b></li><li>Hey\nYou</li></ol>"))
}

func TestBlockquoteSimple(t *testing.T) {
    text := "<blockquote>Hey\nYou!</blockquote>"
    equals(t, "> Hey\n> You!", Convert(text))
}

func TestNewline(t *testing.T) {
    equals(t, "Derpy\n\nHer\n\npy", Convert("Derpy<br/>Her<br>py"))
}

func TestAnchor(t *testing.T) {
    equals(t, "[Linkue](https://ya.ru/a.jpg)", Convert("<a href=\"https://ya.ru/a.jpg\">Linkue</a>"))
}

func TestAnchorWithTitle(t *testing.T) {
    equals(t, "[Linkue](https://ya.ru/a.jpg \"Sasi\")", Convert("<a href=\"https://ya.ru/a.jpg\" title=\"Sasi\">Linkue</a>"))
}

func TestImage(t *testing.T) {
    equals(t, "![Derp](https://ya.ru/a.jpg)", Convert("<img alt=\"Derp\" src=\"https://ya.ru/a.jpg\"></img>"))
}

func TestImageWithTitle(t *testing.T) {
    equals(t, "![Derp](https://ya.ru/a.jpg \"Sasi\")", Convert("<img alt=\"Derp\" src=\"https://ya.ru/a.jpg\" title=\"Sasi\"></img>"))
}

func TestDeepList(t *testing.T) {
    equals(t, "1. Wow\n\t* Dat\n\t* Fat\n1. Hey\nYou", Convert("<ol><li>Wow</li><li>Dat <b>az</b></li><li>Hey\nYou</li></ol>"))
}


// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
    if !condition {
        _, file, line, _ := runtime.Caller(1)
        fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
        tb.FailNow()
    }
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
    if err != nil {
        _, file, line, _ := runtime.Caller(1)
        fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
        tb.FailNow()
    }
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
    if !reflect.DeepEqual(exp, act) {
        _, file, line, _ := runtime.Caller(1)
        fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
        tb.FailNow()
    }
}