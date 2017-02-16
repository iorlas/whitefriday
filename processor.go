package whitefriday

import (
    "strings"
    "golang.org/x/net/html"
    "log"
    "bytes"
    "golang.org/x/net/html/atom"
    "fmt"
)

type State struct {
    IsBold bool
    IsItalic bool
    ListDepth int
}

func Convert(text string) string {
    parsed, err := html.Parse(strings.NewReader(strings.TrimSpace(text)))
    if err != nil {
        log.Fatal(err)
    }

    result, err := parse(State{}, parsed.LastChild.LastChild)
    if err != nil {
        panic(err)
    }
    return result
}

func parse(state State, source *html.Node) (s string, err error) {
    result := bytes.NewBufferString("")
    for node := source.FirstChild; node != nil; node = node.NextSibling {
        var foundParser Parser
        var hasParser bool
        parsersCycle: for _, parser := range parsers {
            for _, atom := range parser.Atoms {
                if atom == node.DataAtom {
                    foundParser = parser
                    hasParser = true
                    break parsersCycle
                }
            }
        }
        if !hasParser {
            err = html.Render(result, node)
        } else if line, err := foundParser.Process(state, node, func(childState State) (string, error) {
            return parse(childState, node)
        }); err == nil {
            _, err = result.WriteString(line)
        }
        if err != nil {
            return s, err
        }
    }
    return result.String(), nil
}

type Parser struct {
    Atoms []atom.Atom
    Process func(State, *html.Node, func(State) (string, error)) (string, error)
}

var parsers = []Parser{{
    Atoms: []atom.Atom{0},
    Process: func(state State, source *html.Node, processChildren func(State) (string, error)) (s string, err error) {
        return source.Data, nil
    },
}, {
    Atoms: []atom.Atom{atom.B, atom.Strong},
    Process: func(state State, source *html.Node, processChildren func(State) (string, error)) (s string, err error) {
        state.IsBold = true
        inner, err := processChildren(state)
        state.IsBold = false
        if err != nil {
            return "", err
        }
        return fmt.Sprintf("**%s**", strings.TrimSpace(inner)), nil
    },
}, {
    Atoms: []atom.Atom{atom.I, atom.Em},
    Process: func(state State, source *html.Node, processChildren func(State) (string, error)) (s string, err error) {
        state.IsItalic = true
        inner, err := processChildren(state)
        state.IsItalic = false
        if err != nil {
            return "", err
        }
        return fmt.Sprintf("*%s*", strings.TrimSpace(inner)), nil
    },
}, {
    Atoms: []atom.Atom{atom.Br},
    Process: func(state State, source *html.Node, processChildren func(State) (string, error)) (s string, err error) {
        return " \n", nil
    },
}, {
    Atoms: []atom.Atom{atom.P},
    Process: func(state State, source *html.Node, processChildren func(State) (string, error)) (s string, err error) {
        inner, err := processChildren(state)
        if err != nil {
            return "", err
        }
        return "\n\n" + strings.TrimSpace(inner) + "\n\n", nil
    },
}, {
    Atoms: []atom.Atom{atom.Ul, atom.Ol},
    Process: func(state State, source *html.Node, processChildren func(State) (string, error)) (s string, err error) {
        state.ListDepth += 1
        inner, err := processChildren(state)
        state.ListDepth -= 1
        if err != nil {
            return "", err
        }
        if state.ListDepth > 0 {
            return "\n" + inner, nil
        }
        return inner, nil
    },
}, {
    Atoms: []atom.Atom{atom.Li},
    Process: func(state State, source *html.Node, processChildren func(State) (string, error)) (s string, err error) {
        buf := bytes.NewBufferString(strings.Repeat("\t", state.ListDepth - 1))
        inner, err := processChildren(state)
        if err != nil {
            return s, err
        }
        if source.Parent.DataAtom == atom.Ol {
            buf.WriteString("1. ")
        } else {
            buf.WriteString("* ")
        }
        buf.WriteString(strings.TrimSpace(inner))
        if source.NextSibling != nil {
            buf.WriteString("\n")
        }
        return buf.String(), nil
    },
}, {
    Atoms: []atom.Atom{atom.Blockquote},
    Process: func(state State, source *html.Node, processChildren func(State) (string, error)) (s string, err error) {
        inner, err := processChildren(state)
        if err != nil {
            return "", err
        }
        buf := bytes.NewBufferString("")
        for i, line := range strings.Split(strings.TrimSpace(inner), "\n") {
            if i != 0 {
                buf.WriteString("\n")
            }
            buf.WriteString("> ")
            buf.WriteString(line)
        }
        return buf.String(), nil
    },
}, {
    Atoms: []atom.Atom{atom.A},
    Process: func(state State, source *html.Node, processChildren func(State) (string, error)) (s string, err error) {
        inner, err := processChildren(state)
        if err != nil {
            return "", err
        }

        var href, title string
        for _, attr := range source.Attr {
            if attr.Key == "href" {
                href = attr.Val
            } else if attr.Key == "title" {
                title = attr.Val
            }
        }
        result := bytes.NewBufferString("[")
        result.WriteString(strings.TrimSpace(inner))
        result.WriteString("](")
        result.WriteString(strings.TrimSpace(href))
        if title != "" {
            result.WriteString(" \"")
            result.WriteString(strings.TrimSpace(title))
            result.WriteString("\"")
        }
        result.WriteString(")")
        return result.String(), nil
    },
}, {
    Atoms: []atom.Atom{atom.Img},
    Process: func(state State, source *html.Node, processChildren func(State) (string, error)) (s string, err error) {
        var src, alt, title string
        for _, attr := range source.Attr {
            switch attr.Key {
            case "src":
                src = attr.Val
            case "title":
                title = attr.Val
            case "alt":
                alt = attr.Val
            }
        }
        result := bytes.NewBufferString("![")
        result.WriteString(strings.TrimSpace(alt))
        result.WriteString("](")
        result.WriteString(strings.TrimSpace(src))
        if title != "" {
            result.WriteString(" \"")
            result.WriteString(strings.TrimSpace(title))
            result.WriteString("\"")
        }
        result.WriteString(")")
        return result.String(), nil
    },
}}
