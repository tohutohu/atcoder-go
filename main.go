package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	client "github.com/tohutohu/atcoder-go/client/atcoder"
	"github.com/urfave/cli"
)

//go:generate go-assets-builder -s="/data" -o bindata.go data

var (
	contestRe = regexp.MustCompile(`^a(r|b|g)c[0-9]{3}$`)
	taskRe    = regexp.MustCompile(`^(a|b|c|d|e|f|g|h)$`)
	tasks     = map[string][]string{
		"abc": []string{"a", "b", "c", "d"},
		"arc": []string{"a", "b", "c", "d"},
		"agc": []string{"a", "b", "c", "d", "e", "f"},
	}
)

func main() {
	app := cli.NewApp()
	app.Name = "atcoder go"
	app.Usage = "useful atcoder support commands"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:    "start",
			Aliases: []string{"s"},
			Usage:   "start atcoder",
			Action: func(c *cli.Context) error {
				if len(c.Args()) != 1 {
					return cli.NewExitError("invalid argument", 1)
				}
				contestName := c.Args().First()
				if contestRe.MatchString(contestName) {
					err := mkdir(contestName)
					if err != nil {
						return err
					}
					contestType := contestName[:3]
					fmt.Println(contestType)
					var wg sync.WaitGroup
					c, err := client.New()
					if err != nil {
						return err
					}
					sem := make(chan struct{}, 10)
					for _, dirName := range tasks[contestType] {
						wg.Add(1)
						sem <- struct{}{}
						problemName := dirName
						go func(dirPath string) {
							defer wg.Done()
							samples, state, _ := c.GetTaskInfo(contestName, problemName)
							setUpCPPDir(dirPath, samples, state)
							<-sem
						}(contestName + "/" + dirName)
					}
					wg.Wait()
					return nil
				}
				return cli.NewExitError("invalid contest name", 1)
			},
		},
		{
			Name: "submit",
			Action: func(ctx *cli.Context) error {
				if len(ctx.Args()) > 2 {
					return cli.NewExitError("invalid arguments", 1)
				}
				fileName := ctx.Args().First()
				contest, task, err := getContestAndTask(fileName)
				if err != nil {
					return err
				}
				if !contestRe.MatchString(contest) || !taskRe.MatchString(task) {
					return cli.NewExitError("invalid file", 1)
				}
				c, err := client.New()
				if err != nil {
					return err
				}
				fmt.Printf("Start submit code, contest:%s task:%s", contest, task)
				ch := make(chan struct{})

				go func() {
					t := time.NewTicker(300 * time.Millisecond)
					for {
						select {
						case <-t.C:
							fmt.Print(".")
						case <-ch:
							fmt.Println("Done")
							t.Stop()
							ch <- struct{}{}
							return
						}
					}
				}()

				file, err := os.Open(fileName)
				if err != nil {
					return err
				}
				body, err := ioutil.ReadAll(file)
				if err != nil {
					return err
				}

				if err := c.Submit(contest, task, string(body)); err != nil {
					fmt.Println("Submit failed")
					return err
				}
				ch <- struct{}{}
				<-ch
				fmt.Println("Submit complete")
				return nil
			},
		},
		{
			Name: "login",
			Action: func(ctx *cli.Context) error {
				_, err := client.New()
				if err != nil {
					return err
				}
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func setUpCPPDir(path string, samples []client.Sample, state string) error {
	if err := mkdir(path); err != nil {
		return err
	}
	mainFile, err := os.Create(path + "/main.go")
	defer mainFile.Close()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	mainFile.Write(Assets.Files["/main.go"].Data)
	testFile, err := os.Create(path + "/main_test.go")
	defer testFile.Close()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	testBody := string(Assets.Files["/main_test.go"].Data)
	testInputs := ""
	testOutputs := ""
	for _, sample := range samples {
		testInputs += fmt.Sprintf("`%s`, ", sample.Input)
		testOutputs += fmt.Sprintf("`%s`, ", sample.Output)
	}
	testBody = strings.Replace(testBody, "\"sampleInput-placeholder\"", testInputs, 1)
	testBody = strings.Replace(testBody, "\"sampleOutput-placeholder\"", testOutputs, 1)
	testFile.Write([]byte(testBody))

	stateFile, err := os.Create(path + "/state.txt")
	defer stateFile.Close()
	if err != nil {
		return err
	}
	rep := regexp.MustCompile(`\n{2,}`)
	stateFile.Write([]byte(rep.ReplaceAllString(state, "\n\n")))
	return nil
}

func mkdir(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		if err = os.MkdirAll(path, 0777); err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
	} else {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}

func getContestAndTask(fileName string) (contest, task string, err error) {
	filePath, err := filepath.Abs(fileName)
	if err != nil {
		return
	}

	dirPath := filepath.Dir(filePath)

	dirs := strings.Split(dirPath, "/")
	contest = dirs[len(dirs)-2]
	task = dirs[len(dirs)-1]
	return
}
