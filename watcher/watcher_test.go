package watcher

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	setup()
	eCode := m.Run()
	cleanup()
	os.Exit(eCode)
}

const dataDir = "tempTestData"

func setup() {
	cleanup() //Run this first just incase a previous test didn't clean up after itself
	dirs := []string{dataDir + "/dir1", dataDir + "/dir2"}
	for _, v := range dirs {
		err := os.MkdirAll(v, 0700)
		if err != nil {
			log.Fatal(err)
		}
		_, err = os.Create(v + "/data")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func cleanup() {
	err := os.RemoveAll(dataDir)
	if err != nil {
		log.Fatal(err)
	}
}

func TestRun(t *testing.T) {
	tt := map[string][]string{
		"Single-Path": {dataDir + "/dir1"},
		"Multi-Path":  {dataDir + "/dir1", dataDir + "/dir2"},
	}
	for k, v := range tt {
		fmt.Println(k, v)
		t.Run(k, func(t *testing.T) {
			task, err := Create(v)
			if err != nil {
				t.Fatal(err)
			}
			go func(t *testing.T) {
				time.Sleep(time.Second * 1) //Give the watcher 1 second to prepare before editing the file
				dat := "Writing some test data"
				err := ioutil.WriteFile(v[0]+"/data", []byte(dat), 0700)
				if err != nil {
					t.Fatal(err)
				}
			}(t)
			task.Run()
		})
	}
}
