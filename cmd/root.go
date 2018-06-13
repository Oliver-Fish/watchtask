package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Oliver-Fish/watchtask/execmanager"
	"github.com/Oliver-Fish/watchtask/watcher"
	"github.com/spf13/cobra"
)

var paths []string
var cmds []string

func init() {
	rootCmd.PersistentFlags().StringSliceVarP(&paths, "path", "p", []string{}, "Paths/Files to watch for changes on")
	rootCmd.PersistentFlags().StringSliceVarP(&cmds, "cmd", "c", []string{}, "Commands to run on changes")
	rootCmd.MarkFlagRequired("cmd")
	if len(paths) == 0 {
		p, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		paths = append(paths, p)
	}
}

var tasks []*execmanager.Task
var rootCmd = &cobra.Command{
	Use:   "watchtask",
	Short: "WatchTask is a tool to run commands on changes in specified directories and/or files",
	Long: `WatchTask is a quick way to run tasks on changes on watched directories and/or files
Examples:
Run Command on changes in current directory
watchtask -c "go run main.go"
Run Multiple Commands on change
watchtask -c "go run main.go","echo TestTest
Run Command on specified path 
watchtask -c "go run main.go" -p "~/webapp"
"
	`, Run: func(cmd *cobra.Command, args []string) {
		//Ensure we have atleast one command arg, otherwise we have no command to run on change detected
		if len(cmds) == 0 {
			log.Fatal("No cmd argument specified, run watchtask -h for usage info")
		}
		//Build our tasks struct with the commands we will need to run on change
		for _, v := range cmds {
			var cPath string
			var cArgs []string
			sCmd := strings.Split(v, " ")
			for k, v := range sCmd {
				if k == 0 {
					cPath = v
				} else {
					cArgs = append(cArgs, v)
				}
			}
			task := execmanager.Task{
				CmdPath: cPath,
				Args:    cArgs,
				Running: false,
			}
			tasks = append(tasks, &task) //Append a pointer of our task to a package scope slice of task, this lets us cleanup and restart processes easily
		}
		//Ensure We cleanup after a Interrupt signal i.e Ctrl+C in terminal
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c        //Block until we get our signal
			shutdown() //Cleanup
			os.Exit(0) //Exit program
		}()
		//Cleanup incase we panic below
		defer shutdown()
		//Loop forever, the watchers will block until they detect changes then the tasks will restart when they unblock
		for {
			log.Print("Changes detected")
			for _, v := range tasks {
				log.Printf("Running command - %v %v", v.CmdPath, strings.Join(v.Args, " "))
				if v.Running {
					err := v.Kill()
					if err != nil {
						panic(err)
					}
				}
				err := v.Start()
				if err != nil {
					panic(err)
				}
				scanner := bufio.NewScanner(v.StdoutPipe) //stdout
				go func() {                               //
					for scanner.Scan() {
						log.Print(scanner.Text())
					}
				}()
			}
			w, err := watcher.Create(paths)
			if err != nil {
				panic(err)
			}
			err = w.Run()
			if err != nil {
				panic(err)
			}
		}

	},
}

//Shutdown attempts to cleans up any processes that we have spawned
func shutdown() {
	for _, v := range tasks {
		if v.Running {
			err := v.Kill()
			if err != nil {
				//The Program is shutting down all we can do is inform that we were unable to cleanup, any panicing or fatal here might stop us killing tasks we can cleanup
				log.Printf("Error: Unable to shutdown %v %v failed with error - %v", v.CmdPath, v.Args, err)
			}
		}
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
