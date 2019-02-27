package main

import (
	"fmt";
	"os";
	"getopt";
	"strings";
)

func main() {
	var nargv []string = os.Args[1:]
	var optmap = map[string][]string{}
	var invalid_optmap = []string{}
	var args []string
	var forward []string

	var options = []getopt.Option {
		{Help: "Options group1:"},
		{Names: []string{"h", "H", "help"}, Argtype: getopt.N, Help: "out put the usage info"},
		{Names: []string{"f", "F", "file"}, Argtype: getopt.M, Help: "file to be parse"},
		{Names: []string{"wenj", "wenjian"}, Link: "f", Hide: true, Help: "deprecated use f instead"},
		{Names: []string{"v"}, Argtype: getopt.N, Help: "verbose output, -vvv means verbose level 3"},
		{Names: []string{"x"}, Argtype: getopt.Y, Help: "dump binary file to text"},
		{Names: []string{"s"}, Argtype: getopt.Y, Help: "enable smart mode", Hide: false},
		{Names: []string{"S"}, Link: "s", Hide: true},

		{Help: "\nOptions group2:"},
		{Names: []string{"e"}, Argtype: getopt.M, Help: "sed -e option, will forward to child sed process", Forward: true},
		{Names: []string{"r"}, Argtype: getopt.N, Help: "sed -r option, will forward to child sed process", Forward: true},
		{Names: []string{"n"}, Argtype: getopt.N, Help: "sed -n option, will forward to child sed process", Forward: true},
	}

	// debug info
	fmt.Println("Debug info:")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Println(nargv)
	getopt.GetUsage(options)

	optmap, invalid_optmap, args, forward = getopt.GetOptions(options, nargv)
	for k, v := range optmap {
		fmt.Printf("opt(%s) = %v\n", k, v)
	}
	fmt.Println()

	for _, v := range invalid_optmap {
		fmt.Println(v)
	}
	fmt.Println()

	fmt.Println("params:", args)
	fmt.Println("forward:", forward)
	fmt.Println(strings.Repeat("-", 80))

	// start your code
	if _, ok := optmap["h"]; ok {
		fmt.Println("Usage:")
		//getopt.GetUsage(options)
	}
	if val, ok := optmap["f"]; ok {
		filelist := val
		fmt.Println("file list:", filelist)
	}
	if val, ok := optmap["v"]; ok {
		verboselevel := len(val)
		fmt.Println("verbose level:", verboselevel)
	}
	if val, ok := optmap["s"]; ok {
		smartmode := val[len(val)-1]
		fmt.Println("smart mode:", smartmode)
	}
}
