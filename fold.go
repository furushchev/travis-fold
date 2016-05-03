package main

import (
	"os"

	"fmt"
	"time"

	"github.com/codegangsta/cli"
	"io"
)

func input(fd io.Reader) <-chan []byte {
	ch := make(chan []byte)
	go func(){
		defer close(ch)
		for {
			data := make([]byte, 1)
			n, err := fd.Read(data)
			if n > 0 {
				ch <- data
			} else {
				if err == io.EOF {
					return
				}
			}
		}
	}()
	return ch
}

func CommandExec(c *cli.Context) error {

	// check stdin is piped
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return fmt.Errorf("%s: no stdin", c.App.Name)
	}

	name := c.String("name")
	if name == "" {
		prefix := c.String("prefix")
		name = prefix + "." + time.Now().Format(c.String("layout"))
	}

	fmt.Fprintln(os.Stdout, "travis_fold:start:"+name)
	defer fmt.Fprintln(os.Stdout, "travis_fold:end:"+name)

	for out := range input(os.Stdin) {
		_, err := os.Stdout.Write(out)
		if err != nil {
			return fmt.Errorf("%s: %s", c.App.Name, err)
		}
	}
	return nil
}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: [Error] '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(1)
}

func main() {
	app := cli.NewApp()
	app.Name = Name
	app.Version = Version
	app.Author = "furushchev"
	app.Email = "furushchev@mail.ru"
	app.Copyright = "MIT"
	app.Usage = ""

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "name, n",
			Value: "",
			Usage: "name to be shown on Travis-CI web console",
		},
		cli.StringFlag{
			Name:   "prefix, p",
			Value:  "command",
			Usage:  "prefix for tag on Travis-CI",
			EnvVar: "TRAVIS_FOLD_PREFIX",
		},
		cli.StringFlag{
			Name:  "layout, l",
			Value: "15.04.05",
			Usage: "datetime layout used for tag on Travis-CI",
		},
	}
	app.Action = CommandExec
	app.CommandNotFound = CommandNotFound

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, "%s: [Error] %s", app.Name, err)
		os.Exit(1)
	}
}
