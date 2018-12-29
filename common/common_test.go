package common_test

import (
	"errors"
	"testing"

	"github.com/n0trace/photon/common"
)

func panicked(f func()) (ret bool) {
	defer func() {
		if x := recover(); x != nil {
			ret = true
		}
	}()
	f()
	return
}
func TestMust(t *testing.T) {
	f := func() error {
		return errors.New("test error")
	}

	var wantRet = true
	ret := panicked(func() {
		common.Must(f())
	})

	if ret != wantRet {
		t.Errorf("Must() ret = %v, wantRet %v", ret, wantRet)
		return
	}

	wantRet = false

	ret = panicked(func() {
		common.Must(nil)
	})

	if ret != wantRet {
		t.Errorf("Must() ret = %v, wantRet %v", ret, wantRet)
		return
	}

}

func TestMust2(t *testing.T) {
	f := func() error {
		return errors.New("test error")
	}

	var wantRet = true
	ret := panicked(func() {
		common.Must2(nil, f())
	})

	if ret != wantRet {
		t.Errorf("Must2() ret = %v, wantRet %v", ret, wantRet)
		return
	}

	wantRet = false

	ret = panicked(func() {
		common.Must2(nil, func() error { return nil }())
	})

	if ret != wantRet {
		t.Errorf("Must2() ret = %v, wantRet %v", ret, wantRet)
		return
	}
}
