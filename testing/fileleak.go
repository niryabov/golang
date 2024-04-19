//go:build !solution

package fileleak

import (
	"log"
	"os"
)

type testingT interface {
	Errorf(msg string, args ...interface{})
	Cleanup(func())
}

func VerifyNone(t testingT) {
	files, err := os.ReadDir("/proc/self/fd")
	if err != nil {
		log.Fatal(err)
	}
	opened_files := make(map[string]bool)
	descriptions := make(map[string]string)
	for _, file := range files {
		str, _ := os.Readlink("/proc/self/fd/" + file.Name())
		opened_files[file.Name()] = true
		descriptions[file.Name()] = str

	}
	t.Cleanup(func() {
		files, err := os.ReadDir("/proc/self/fd")
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			str, _ := os.Readlink("/proc/self/fd/" + file.Name())
			val, ok := opened_files[file.Name()]
			if !ok || (str != descriptions[file.Name()]) {
				t.Errorf("Opened file with fd %v", val)
			}
		}
	})

}
