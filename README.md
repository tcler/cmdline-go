# getopt.go
A getopt implementation in golang. Very similar to [getopt.tcl](https://github.com/tcler/getopt.tcl)

features in my getopt.tcl:
1. generate usage/help info from option list.
2. support GNU style option and more flexible: -a --along --b -c carg -d=darg -ooptionalarg -- --notoption
2. not just support a short and a long option, you can define a *List* {h help Help ? 帮助}
3. hide attribute of option object， used to hide some option in usage/help info
4. forward option

# Example code
see here: https://github.com/tcler/getopt.go/blob/master/main.go
