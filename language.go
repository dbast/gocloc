package main

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Language struct {
	name           string
	line_comment   string
	multi_line     string
	multi_line_end string
	files          []string
	code           int32
	comments       int32
	blanks         int32
	lines          int32
	total          int32
	printed        bool
}
type Languages []Language

var reShebangEnv *regexp.Regexp = regexp.MustCompile("^#! *(\\S+/env) ([a-zA-Z]+)")
var reShebangLang *regexp.Regexp = regexp.MustCompile("^#! *[.a-zA-Z/]+/([a-zA-Z]+)")

func (ls Languages) Len() int {
	return len(ls)
}
func (ls Languages) Swap(i, j int) {
	ls[i], ls[j] = ls[j], ls[i]
}
func (ls Languages) Less(i, j int) bool {
	return ls[i].code > ls[j].code
}

var Exts map[string]string = map[string]string{
	"as":       "as",
	"asm":      "s",
	"S":        "s",
	"s":        "s",
	"awk":      "awk",
	"bat":      "bat",
	"btm":      "bat",
	"cmd":      "bat",
	"bash":     "bash",
	"sh":       "sh",
	"c":        "c",
	"csh":      "csh",
	"ec":       "c",
	"pgc":      "c",
	"cs":       "cs",
	"clj":      "clj",
	"coffee":   "coffee",
	"cfm":      "cfm",
	"cfc":      "cfc",
	"cmake":    "cmake",
	"cc":       "cpp",
	"cpp":      "cpp",
	"cxx":      "cpp",
	"pcc":      "cpp",
	"c++":      "cpp",
	"css":      "css",
	"d":        "d",
	"dart":     "dart",
	"dts":      "dts",
	"dtsi":     "dts",
	"el":       "lisp",
	"lisp":     "lisp",
	"lsp":      "lisp",
	"lua":      "lua",
	"sc":       "lisp",
	"f":        "f77",
	"f77":      "f77",
	"for":      "f77",
	"ftn":      "f77",
	"pfo":      "f77",
	"f90":      "f90",
	"f95":      "f90",
	"f03":      "f90",
	"f08":      "f90",
	"go":       "go",
	"h":        "h",
	"hs":       "hs",
	"hpp":      "hpp",
	"hh":       "hpp",
	"html":     "html",
	"hxx":      "hpp",
	"jai":      "jai",
	"java":     "java",
	"js":       "js",
	"jl":       "jl",
	"json":     "json",
	"jsx":      "jsx",
	"lds":      "lds",
	"less":     "less",
	"m":        "m",
	"md":       "md",
	"markdown": "md",
	"ml":       "ml",
	"mli":      "ml",
	"mm":       "mm",
	"makefile": "makefile",
	"mustache": "mustache",
	"m4":       "m4",
	"l":        "lex",
	"php":      "php",
	"pas":      "pas",
	"PL":       "pl",
	"pl":       "pl",
	"pm":       "pl",
	"plan9sh":  "plan9sh",
	"text":     "text",
	"txt":      "text",
	"polly":    "polly",
	"proto":    "proto",
	"py":       "py",
	"r":        "r",
	"rake":     "rb",
	"rb":       "rb",
	"rhtml":    "rhtml",
	"rs":       "rs",
	"sass":     "sass",
	"scss":     "sass",
	"sml":      "sml",
	"sql":      "sql",
	"swift":    "swift",
	"tex":      "tex",
	"sty":      "tex",
	"toml":     "toml",
	"ts":       "ts",
	"vim":      "vim",
	"xml":      "xml",
	"xsl":      "xsl",
	"xslt":     "xsl",
	"yaml":     "yaml",
	"yml":      "yaml",
	"y":        "y",
	"zsh":      "zsh",
}

var LanguageByScript map[string]string = map[string]string{
	"make":   "make",
	"perl":   "pl",
	"rc":     "plan9sh",
	"python": "py",
	"ruby":   "rb",
}

func getShebang(line string) (shebangLang string, ok bool) {
	if reShebangEnv.MatchString(line) {
		ret := reShebangEnv.FindAllStringSubmatch(line, -1)
		if len(ret[0]) == 3 {
			shebangLang = ret[0][2]
			if sl, ok := LanguageByScript[shebangLang]; ok {
				return sl, ok
			}
			return shebangLang, true
		}
	}

	if reShebangLang.MatchString(line) {
		ret := reShebangLang.FindAllStringSubmatch(line, -1)
		if len(ret[0]) >= 2 {
			shebangLang = ret[0][1]
			if sl, ok := LanguageByScript[shebangLang]; ok {
				return sl, ok
			}
			return shebangLang, true
		}
	}

	return "", false
}

func getFileTypeByShebang(path string) (shebangLang string, ok bool) {
	line := ""
	func() {
		fp, err := os.Open(path)
		if err != nil {
			return // ignore error
		}
		defer fp.Close()

		scanner := bufio.NewScanner(fp)
		for scanner.Scan() {
			l := scanner.Text()
			line = strings.TrimSpace(l)
			break
		}
	}()

	shebangLang, ok = getShebang(line)
	return shebangLang, ok
}

func getFileType(path string) (ext string, ok bool) {
	ext = filepath.Ext(path)
	base := filepath.Base(path)

	switch base {
	case "CMakeLists.txt":
		return "cmake", true
	case "configure.ac":
		return "m4", true
	case "Makefile.am":
		return "makefile", true
	}

	switch strings.ToLower(base) {
	case "makefile":
		return "makefile", true
	}

	shebangLang, ok := getFileTypeByShebang(path)
	if ok {
		return shebangLang, true
	}

	if len(ext) >= 2 {
		return ext[1:], true
	}
	return ext, ok
}

func NewLanguage(name, line_comment, multi_line, multi_line_end string) *Language {
	return &Language{
		name:           name,
		line_comment:   line_comment,
		multi_line:     multi_line,
		multi_line_end: multi_line_end,
		files:          []string{},
	}
}
