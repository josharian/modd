package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/cortesi/modd"
	"github.com/cortesi/modd/conf"
	"github.com/cortesi/termlog"
	"gopkg.in/alecthomas/kingpin.v2"
)

const modfile = "./modd.conf"
const lullTime = time.Millisecond * 300

func main() {
	file := kingpin.Flag(
		"file",
		fmt.Sprintf("Path to modfile (%s)", modfile),
	).
		Default(modfile).
		PlaceHolder("PATH").
		Short('f').
		String()
	//
	// paths := kingpin.Arg(
	// 	"path",
	// 	"Paths to monitor for changes.",
	// ).Required().Strings()
	//
	// beep := kingpin.Flag("beep", "Beep if any command returned an error").
	// 	Short('b').
	// 	Bool()
	//
	// nocommon := kingpin.Flag("nocommon", "Don't exclude commonly ignored files").
	// 	Short('c').
	// 	Bool()
	//
	// daemons := kingpin.Flag("daemon", "Daemon to keep running").
	// 	PlaceHolder("CMD").
	// 	Short('d').
	// 	Strings()
	//
	// prep := kingpin.Flag("prep", "Prep command to run before daemons are restarted").
	// 	PlaceHolder("CMD").
	// 	Short('p').
	// 	Strings()

	// excludes := kingpin.Flag("exclude", "Glob pattern for files to exclude from monitoring").
	// 	PlaceHolder("PATTERN").
	// 	Short('x').
	// 	Strings()

	cmdstats := kingpin.Flag("cmdstats", "Show stats on command execution").
		Short('s').
		Default("false").
		Bool()

	debug := kingpin.Flag("debug", "Debugging for devd development").
		Default("false").
		Bool()

	kingpin.Version(modd.Version)
	kingpin.Parse()
	log := termlog.NewLog()

	if *debug {
		log.Enable("debug")
		modd.Logger = log
	}
	if *cmdstats {
		log.Enable("cmdstats")
	}

	ret, err := ioutil.ReadFile(*file)
	if err != nil {
		kingpin.Fatalf("%s", err)
	}
	cnf, err := conf.Parse(*file, string(ret))
	if err != nil {
		kingpin.Fatalf("%s", err)
	}
	fmt.Println(cnf)

	// modchan := make(chan modd.Mod)
	// exc := *excludes
	// if !*nocommon {
	// 	exc = append(*excludes, modd.CommonExcludes...)
	// }
	// err := modd.Watch(*paths, exc, lullTime, modchan)
	// if err != nil {
	// 	kingpin.Fatalf("Fatal error: %s", err)
	// }
	// err = modd.RunProcs(*prep, log)
	// if err != nil {
	// 	if *beep {
	// 		fmt.Print("\a")
	// 	}
	// }
	// d := modd.DaemonPen{}
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt, os.Kill)
	// go func() {
	// 	d.Shutdown(<-c)
	// 	os.Exit(0)
	// }()
	// d.Start(*daemons, log)
	// for mod := range modchan {
	// 	log.SayAs("debug", "Delta: \n%s", mod.String())
	// 	err := modd.RunProcs(*prep, log)
	// 	if err != nil {
	// 		if *beep {
	// 			fmt.Print("\a")
	// 		}
	// 		continue
	// 	}
	// 	d.Restart()
	// }
}
