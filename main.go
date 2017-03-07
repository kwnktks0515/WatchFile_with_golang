package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

//FileInfo is ABC
type FileInfo struct {
	Path string
	Time time.Time
}

func main() {
	Watch(func() {
		fmt.Println("Hello")
	})
}

//Watch is main function
func Watch(callbacl func()) {
	var list []FileInfo
	if deep, arg, err := getarg(); err == nil {
		if list, err = getpath(deep, arg); err != nil {
			fmt.Println(err)
			return
		}
	} else {
		fmt.Println(err)
		return
	}
	for _ = range time.Tick(1 * time.Second) {
		if ischange(&list) {
			callbacl()
		}
	}

}

//Init is get file path
func getpath(deep int, args []string) ([]FileInfo, error) {
	var list []FileInfo
	for _, arg := range args {
		max := deep - 1
		flag := true
		for _, str := range strings.Split(arg, "/") {
			max++
			if flag && str == ".." {
				max++
				flag = false
			}
		}
		filepath.Walk(arg, func(path string, info os.FileInfo, err error) error {
			if max < len(strings.Split(path, "/")) {
				return filepath.SkipDir
			}
			if !info.IsDir() {
				list = append(list, FileInfo{path, info.ModTime()})
			}
			//fmt.Println(deep, path, info.IsDir(), info.Name(), info.Mode(), info.ModTime())
			return nil
		})
	}
	return list, nil
}

func getarg() (int, []string, error) {
	var array []string
	deep, err := strconv.Atoi(os.Args[1])
	if err != nil {
		return 0, nil, err
	}
	for _, arg := range os.Args[2:] {
		if _, err := os.Stat(arg); err != nil {
			return 0, nil, err
		}
		if arg[len(arg)-1] == '/' {
			arg = arg[:len(arg)-1]
		}
		array = append(array, arg)
	}
	return deep, array, nil
}

func ischange(list *[]FileInfo) bool {
	flag := false
	for i := range *list {
		if time, _ := os.Stat((*list)[i].Path); time.ModTime().After((*list)[i].Time) {
			(*list)[i].Time = time.ModTime()
			flag = true
		}
	}
	return flag
}
