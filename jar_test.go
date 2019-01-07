package photon

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestJar_Cookies(t *testing.T) {
	cookie := &http.Cookie{Name: "hello", Value: "world"}
	jar := &Jar{}
	u := new(url.URL)
	jar.SetCookies(u, []*http.Cookie{cookie})
	got := jar.Cookies(u)
	if !reflect.DeepEqual(got[0], cookie) {
		t.Errorf("Jar.Cookies()[0] = %v, want %v", got, cookie)
	}
}
