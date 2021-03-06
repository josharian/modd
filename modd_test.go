package modd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/cortesi/modd/conf"
	"github.com/cortesi/modd/utils"
	"github.com/cortesi/modd/watch"
	"github.com/cortesi/termlog"
)

const timeout = 2 * time.Second

func touch(t *testing.T, p string) {
	p = filepath.FromSlash(p)
	err := ioutil.WriteFile(p, []byte("teststring"), 0777)
	if err != nil {
		t.Fatalf("touch: %s", err)
	}
}

func events(p string) []string {
	parts := []string{}
	for _, p := range strings.Split(p, "\n") {
		if strings.HasPrefix(p, ":") {
			p = strings.TrimSpace(p)
			if !strings.HasSuffix(p, ":") {
				parts = append(parts, strings.TrimSpace(p))
			}
		}
	}
	return parts
}

func _testWatch(t *testing.T, modfunc func(), trigger string, expected []string) {
	defer utils.WithTempDir(t)()

	err := os.MkdirAll("a", 0777)
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll("b", 0777)
	if err != nil {
		t.Fatal(err)
	}
	confTxt := `
        ** {
            prep: echo ":all:" @mods
        }
        a/** {
            prep: echo ":a:" @mods
        }
        b/** {
            prep: echo ":b:" @mods
        }
    `
	cnf, err := conf.Parse("test", confTxt)
	if err != nil {
		t.Fatal(err)
	}

	lt := termlog.NewLogTest()

	modchan := make(chan *watch.Mod, 1024)
	cback := func() {
		// There's some race condition in rjeczalik/notify. If we don't wait a
		// bit here, we sometimes don't receive notifications for our changes.
		time.Sleep(200 * time.Millisecond)
		start := time.Now()
		modfunc()
		for {
			if strings.Contains(lt.String(), trigger) {
				break
			}
			if time.Now().Sub(start) > timeout {
				fmt.Println("Timeout!")
				break
			}
			time.Sleep(50 * time.Millisecond)
		}
		modchan <- nil
	}
	_, err = runOnChan(modchan, cback, lt.Log, cnf, "", nil)
	if err != nil {
		t.Fatalf("runOnChan: %s", err)
	}
	ret := events(lt.String())

	if !reflect.DeepEqual(ret, expected) {
		t.Errorf("Expected\n%#v\nGot\n%#v", expected, ret)
	}
}

func TestWatch(t *testing.T) {
	_testWatch(
		t,
		func() { touch(t, "a/touched") },
		"touched",
		[]string{
			":all: ./a/touched",
			":a: ./a/touched",
		},
	)
	_testWatch(
		t,
		func() {
			touch(t, "a/touched")
			touch(t, "b/touched")
		},
		"touched",
		[]string{
			":all: ./a/touched ./b/touched",
			":a: ./a/touched",
			":b: ./b/touched",
		},
	)
}
