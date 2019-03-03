/*
 *
 * getopt -- go implementation of https://github.com/tcler/getopt.tcl
 * 
 * (C) 2019 Jianhong Yin <yin-jianhong@163.com>
 *
 */

package getopt

import (
	"fmt";
	"strings";
	"regexp";
	"os";
)

type ParseStat int
const (
        NOTOPT ParseStat = iota
        KNOWN
        NEEDARG
        UNKNOWN
	END
	AGAIN
)

type ArgType int
const (
	U ArgType = iota
	N
	O
	Y
	M
)

type Option struct {
   Names []string
   Argtype ArgType
   Help string
   Link string
   Hide bool
   Forward bool
}

func getOptObj(options []Option, optname string, followlink bool, nesting *int) *Option {
	if followlink {
		*nesting += 1
		if *nesting < 0 {
			*nesting = 1
		}

		if (*nesting > 4) {
			return nil
		}
	}

	for _, opt := range options {
		for _, n := range opt.Names {
			if n == optname {
				if followlink && len(opt.Link) > 0 {
					return getOptObj(options, opt.Link, followlink, nesting)
				} else {
					return &opt
				}
			}
		}
	}
	return nil
}

func argparse (options []Option, argv []string) (ParseStat, []string, string, string) {
	var result ParseStat = UNKNOWN
	var nargv []string = argv
	var optname string
	var optarg string
	var hasval bool
	var val string
	var opt *Option
	var argtype ArgType

	if len(nargv) == 0 {
		return END, nargv, optname, optarg
	}

	rarg := nargv[0]
	nargv = nargv[1:]

	if rarg == "-" {
		optarg = rarg
		return NOTOPT, nargv, optname, optarg
	}

	if rarg == "--" {
		return END, nargv, optname, optarg
	}

	r, _ := regexp.Compile("^-.*")
	switch {
	case r.MatchString(rarg):
		opttype := "long"
		optname = rarg[1:]
		if optname[0:1] == "-" {
			optname = optname[1:]
		} else {
			opttype = "short"
		}

		idx := strings.Index(optname, "=")
		if idx > -1 {
			ooptname := optname
			toptname := ooptname[:idx]
			if nil != getOptObj(options, toptname, false, nil) {
				optname = toptname
				val = ooptname[idx+1:]
				hasval = true
			}
		}

		opt = getOptObj(options, optname, false, nil)
		if opt != nil {
			argtype = opt.Argtype
			if len(opt.Link) > 0 {
				optlink2 := getOptObj(options, opt.Link, false, nil)
				if optlink2 != nil {
					argtype = optlink2.Argtype
				}
			}

			optname = opt.Names[0]
			result = KNOWN
			switch argtype {
			case O:
				if hasval {
					optarg = val
				}
			case Y, M:
				if hasval {
					optarg = val
				} else if (len(nargv) > 0 && nargv[0] != "--") {
					optarg = nargv[0]
					nargv = nargv[1:]
				} else {
					result = NEEDARG
				}
			}
		} else if hasval == false && opttype == "short" && len(optname) > 1 {
			var argv2 []string
			ShortLoop:
			for (len(optname) > 0) {
				s := optname[0:1]
				optname = optname[1:]

				if s == "=" || s == "-" || s == "\\" || s == "'" || s == "\"" {
					break
				}

				opt = getOptObj(options, s, false, nil)
				if opt == nil {
					argv2 = append(argv2, "-" + s)
					continue
				} else {
					argtype = opt.Argtype
					if len(opt.Link) > 0 {
						optlink2 := getOptObj(options, opt.Link, false, nil)
						if optlink2 != nil {
							argtype = optlink2.Argtype
						}
					}

					switch argtype {
					case O:
						argv2 = append(argv2, "-" + s + "=" + optname)
						break ShortLoop
					default:
						argv2 = append(argv2, "-" + s)
					}
				}
			}
			nargv = append(argv2, nargv...)
			return AGAIN, nargv, optname, optarg
		} else {
			result = UNKNOWN
		}
	default:
		optarg = rarg
		result = NOTOPT
	}

	return result, nargv, optname, optarg
}

func GetOptions (options []Option, argv []string) (map[string][]string, []string, []string, []string) {
	var optmap = make(map[string][]string)
	var invalid_opts = []string{}
	var args []string
	var forward []string

	var opt *Option
	var optname string
	var optarg string
	var nargv []string = argv
	var stat ParseStat
	var nesting int

	Parseloop:
	for len(nargv) > 0 {
		prefix := "-"
		curarg := nargv[0]
		if len(curarg) > 1 && curarg[0:2] == "--" {
			prefix = "--"
		}

		stat, nargv, optname, optarg = argparse(options, nargv)

		switch stat {
		case AGAIN:
			continue
		case NOTOPT:
			args = append(args, optarg)
		case KNOWN:
			nesting = 0
			opt = getOptObj(options, optname, true, &nesting)
			if opt == nil {
				fmt.Fprintf(os.Stderr, "[Warn] (%s) there might be link loop in your option list\n", optname)
				opt = getOptObj(options, optname, false, nil)
			}
			optname = opt.Names[0]

			if opt.Forward {
				switch opt.Argtype {
				case N:
					forward = append(forward, prefix + optname)
				default:
					if prefix == "--" {
						forward = append(forward, prefix + optname + "=" + optarg)
					} else {
						forward = append(forward, prefix + optname + " " + optarg)
					}
				}
				continue
			}
			switch opt.Argtype {
			case M:
				optmap[optname] = append(optmap[optname], optarg)
			case N:
				optmap[optname] = append(optmap[optname], "set")
			default:
				optmap[optname] = append([]string{optarg}, optmap[optname]...)
			}
		case NEEDARG:
			invalid_opts = append(invalid_opts, "option: '" + optname + "' need argument")
		case UNKNOWN:
			invalid_opts = append(invalid_opts, "option: '" + optname + "' undefined")
		case END:
			//end of nargv or get --
			args = append(args, nargv...)
			break Parseloop
		}
	}

	return optmap, invalid_opts, args, forward
}

func genOptdesc(names []string) string {
	var ss string
	var ls string

	for _, n := range names {
		if len(n) == 1 {
			ss += " -" + n
		} else {
			ls += " --" + n
		}
	}

	return ss + ls
}

func GetUsage (options []Option) {
	for _, opt := range options {
		if opt.Hide {
			continue
		}
		if len(opt.Names) == 0 {
			if len(opt.Help) > 0 {
				fmt.Printf("%s\n", opt.Help)
			}
			continue
		}

		pad := 26
		argdesc := ""
		optdesc := genOptdesc(opt.Names)
		switch opt.Argtype {
		case O: argdesc = " [arg]"
		case Y: argdesc = " <arg>"
		case M: argdesc = " {arg}"
		}
		opthelp := opt.Help
		if opt.Help == "" {
			opthelp = "nil #no help found for this options"
		}

		optlen := len(argdesc) + len(optdesc)
		helplen := len(opthelp)

		if optlen > pad-4 && helplen > 8 {
			fmt.Printf("    %-*s\n %*s    %s\n", pad, optdesc+argdesc, pad, "", opthelp)
		} else {
			fmt.Printf("    %-*s %s\n", pad, optdesc+argdesc, opthelp)
		}
	}

	fmt.Println("\nComments:")
	fmt.Println("    *  [arg] means arg is optional, need use --opt=arg to specify an argument")
	fmt.Println("       <arg> means arg is required, and -f a -f b will get the latest 'b'")
	fmt.Println("       {arg} means arg is required, and -f a -f b will get a list 'a b'")
	fmt.Println("")
	fmt.Println("    *  if arg is required, '--opt arg' is same as '--opt=arg'")
	fmt.Println("")
	fmt.Println("    *  '-opt' will be treated as:")
	fmt.Println("           '--opt'    if 'opt' is defined;")
	fmt.Println("           '-o -p -t' if 'opt' is undefined;")
	fmt.Println("           '-o -p=t'  if 'opt' is undefined and '-p' need an argument;")
	fmt.Println("")
}
