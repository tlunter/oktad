package main

import "fmt"

import "github.com/jessevdk/go-flags"
import "github.com/tj/go-debug"
import "github.com/peterh/liner"

func main() {
	var opts struct {
		ConfigFile string `short:"c" long:"config" description:"Path to config file"`
	}

	debug := debug.Debug("oktad:main")
	args, err := flags.Parse(&opts)
	if err != nil {
		return
	}

	debug("loading configuration data")
	// try to load configuration
	oktaCfg, err := parseConfig(opts.ConfigFile)

	if err != nil {
		fmt.Println("Error reading config file!")
		debug("cfg read err: %s", err)
		return
	}

	if len(args) <= 0 {
		fmt.Println("You must supply a profile name, sorry.")
		return
	}

	destArn, err := readAwsProfile(
		fmt.Sprintf("profile %s", args[0]),
	)

	if err != nil {
		fmt.Println("Error reading AWS configuration!")
		return
	}

	user, pass, err := readUserPass()
	if err != nil {
		// if we got an error here, the user bailed on us
		debug("control-c caught in liner, probably")
		return
	}

	if user == "" || pass == "" {
		fmt.Println("Must supply a username and password!")
		return
	}

	err = login(oktaCfg, user, pass, destArn)
	if err != nil {
		fmt.Println("Error grabbing temporary credentials!")
		debug("login err %s", err)
		return
	}
}

// reads the username and password from the command line
// returns user, then pass, then an error
func readUserPass() (user string, pass string, err error) {
	li := liner.NewLiner()

	// remember to close or weird stuff happens
	defer li.Close()

	li.SetCtrlCAborts(true)
	user, err = li.Prompt("Username: ")
	if err != nil {
		return
	}

	pass, err = li.PasswordPrompt("Password: ")
	if err != nil {
		return
	}

	return
}