//go:generate go-bindata templates
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	sh "github.com/codeskyblue/go-sh"
)

var (
	pathTemplates = "templates"

	execfunc = flag.String("execfunc", "", "path to where the execution function is defined")

	binName = "lambdaserver"
)

func init() {
	flag.Parse()

	if *execfunc == "" {
		panic("execfunc must be non-zeri")
	}
}

func removeTmpDir(dirName string) {
	if rmerr := os.RemoveAll(dirName); rmerr != nil {
		panic(rmerr)
	}
}

func retryCreateTmpDir(dirName string) error {
	err := os.Mkdir(dirName, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func createTmpDir(dirName string) error {
	err := os.Mkdir(dirName, os.ModePerm)
	if _, ok := err.(*os.PathError); ok {
		removeTmpDir(dirName)
		return retryCreateTmpDir(dirName)
	} else if err != nil {
		return err
	}
	return nil
}

func main() {
	dirName := "tmp"

	err := createTmpDir(dirName)
	if err != nil {
		panic(err)
	}

	byt, err := Asset(fmt.Sprintf("%v/endpoint.gotmpl", pathTemplates))
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(fmt.Sprintf("%v/lambdaserver.go", dirName), byt, 0666)
	if err != nil {
		panic(err)
	}

	execbyt, err := ioutil.ReadFile(*execfunc)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(fmt.Sprintf("%v/execfunc.go", dirName), execbyt, 0666)
	if err != nil {
		panic(err)
	}

	err = sh.NewSession().SetDir("./"+dirName).
		Command("go", "build", "-o", "../"+binName, "lambdaserver.go", "execfunc.go").Run()

	if err != nil {
		panic(err)
	}

	removeTmpDir(dirName)

}
