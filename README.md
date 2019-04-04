# cmdline-go
A getopt/argparser implementation in golang. Very similar to [getopt.tcl](https://github.com/tcler/getopt.tcl)

features in my cmdline-go:
1. generate usage/help info from option list.
2. support GNU style option and more flexible: -a --along --b -c carg -d=darg -ooptionalarg -- --notoption
2. not just support a short and a long option, you can define a *List* {h help Help ? 帮助}
3. hide attribute of option object， used to hide some option in usage/help info
4. forward option

# Example code
see here: https://github.com/tcler/cmdline-go/blob/master/main.go

```
git clone https://github.com/tcler/cmdline-go
cd cmdline-go

go get -u github.com/tcler/cmdline-go/cmdline
go run main_cmdline.go -h -H -f file --file file2 -e 's/abc/xyz/'  -r -n  -s=A -oa=b -S ''  -i -x xfile --wenj=file3 --www -aa -vvv -S DD -- -0 -y

# or:
go get -u github.com/tcler/cmdline-go/getopt
go run main_getopt.go -h -H -f file --file file2 -e 's/abc/xyz/'  -r -n  -s=A -oa=b -S ''  -i -x xfile --wenj=file3 --www -aa -vvv -S DD -- -0 -y
```
