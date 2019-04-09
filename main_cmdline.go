package main

import (
	"fmt";
	"os";
	"strings";
	"github.com/tcler/cmdline-go/cmdline";
)

func main() {
	var nargv []string = os.Args[1:]
	var cli cmdline.Cmdline

	var options = []cmdline.Option {
		{Help: "Options group1:"},
		{Names: "h H help", Argtype: cmdline.N, Help: "out put the usage info"},
		{Names: "f F file", Argtype: cmdline.M, Help: "file to be parse"},
		{Names: "wenj wenjian", Link: "f", Hide: true, Help: "deprecated use f instead"},
		{Names: "o", Argtype: cmdline.O, Help: "mount option"},
		{Names: "v", Argtype: cmdline.N, Help: "verbose output, -vvv means verbose level 3"},
		{Names: "x", Argtype: cmdline.Y, Help: "dump binary file to text"},
		{Names: "s", Argtype: cmdline.Y, Help: "enable smart mode", Hide: false},
		{Names: "S", Link: "s", Hide: true},

		{Help: "\nOptions group2:"},
		{Names: "e", Argtype: cmdline.M, Help: "sed -e option, will forward to child sed process", Forward: true},
		{Names: "r", Argtype: cmdline.N, Help: "sed -r option, will forward to child sed process", Forward: true},
		{Names: "n", Argtype: cmdline.N, Help: "sed -n option, will forward to child sed process", Forward: true},
	}

	// debug info
	fmt.Println("Debug info:")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Println(nargv)
	cmdline.GetUsage(options)

	cli = cmdline.Parse(options, nargv)
	for k, v := range cli.OptionMap {
		fmt.Printf("opt(%s) = %#v\n", k, v)
	}
	fmt.Println()

	for _, v := range cli.InvalidOptions {
		fmt.Println(v)
	}
	fmt.Println()

	fmt.Printf("params: %#v\n", cli.Args)
	fmt.Printf("forward: %#v\n", cli.ForwardOptions)
	fmt.Println(strings.Repeat("-", 80))

	// start your code
	if cli.HasOption("help") {
		fmt.Println("Usage: ...")
		//cli.GetUsage(options)
	}
	if cli.HasOption("file") {
		fmt.Printf("file list: %#v\n", cli.GetOptionArgList("file"))
	}
	if cli.HasOption("v") {
		verboselevel := cli.GetOptionNumber("v")
		fmt.Println("verbose level:", verboselevel)
	}
	if cli.HasOption("s") {
		smartmode := cli.GetOptionArgString("s")
		fmt.Printf("smart mode: \"%s\"\n", smartmode)
	}
}
