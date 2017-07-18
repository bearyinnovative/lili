package util

import (
	"io/ioutil"
	"log"
	"os/exec"
	"syscall"
)

func init() {
	// log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("init util")
}

func Diff(s1, s2 string) string {
	if s1 == s2 {
		return ""
	}

	p1 := "/tmp/zhihu_old.txt"
	ioutil.WriteFile(p1, []byte(s1), 0644)
	p2 := "/tmp/zhihu_new.txt"
	ioutil.WriteFile(p2, []byte(s2), 0644)

	output, err := ExecBash("diff", "-u", p1, p2)
	if err == nil {
		return output
	}
	eerr, ok := err.(*exec.ExitError)
	if !ok {
		LogIfErr(err)
		return ""
	}
	ws, ok := eerr.Sys().(syscall.WaitStatus)
	if !ok {
		LogIfErr(err)
		return ""
	}
	if ws.ExitStatus() != 1 {
		LogIfErr(err)
		return ""
	}

	// Exit status of 1 means no error, but diffs were found.
	return output
}

func ExecBash(cmd string, arg ...string) (string, error) {
	out, err := exec.Command(cmd, arg...).Output()
	if err != nil {
		return string(out), err
	}

	return string(out), nil
}

func FatalIfErr(err error) {
	if err == nil {
		return
	}

	log.Fatal(err)
}

func LogIfErr(err error) bool {
	if err == nil {
		return false
	}

	log.Println("[ERROR]", err)
	panic(err)
	return true
}

func Log(v ...interface{}) {
	log.Println(v)
}
