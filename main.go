package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/urfave/cli"
)

//go:generate go-assets-builder -s="/data" -o bindata.go data

type Sample struct {
	Input  string
	Output string
}

func main() {
	app := cli.NewApp()
	app.Name = "atc"
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
				if strings.Contains(contestName, "abc") {
					err := mkdir(contestName)
					if err != nil {
						return err
					}
					var wg sync.WaitGroup
					sem := make(chan struct{}, 10)
					for _, dirName := range []string{"/a", "/b", "/c", "/d"} {
						wg.Add(1)
						sem <- struct{}{}
						problemName := dirName[1:]
						go func(dirPath string) {
							defer wg.Done()
							samples, err := getSample(contestName, problemName)
							if err != nil {
								panic(err)
							}
							setUpCPPDir(dirPath, samples)
							<-sem
						}(contestName + dirName)
					}
					wg.Wait()
					return nil
				} else if strings.Contains(contestName, "arc") {

					return nil
				}
				return cli.NewExitError("invalid contest name", 1)
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func setUpCPPDir(path string, samples []Sample) error {
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

func getSample(contest, problem string) ([]Sample, error) {
	samples := []Sample{}
	url := fmt.Sprintf("https://beta.atcoder.jp/contests/%s/tasks/%s_%s", contest, contest, problem)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return samples, err
	}
	sample := Sample{}
	doc.Find("div.part>section>pre").Each(func(_ int, s *goquery.Selection) {
		if s.Parent().Parent().Parent().HasClass("io-style") {
			return
		}
		if sample.Input == "" {
			sample.Input = s.Text()
		} else {
			sample.Output = s.Text()
			samples = append(samples, sample)
			sample = Sample{}
		}
	})
	return samples, err
}
