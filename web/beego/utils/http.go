package utils

import (
	"bytes"
	"io"
	"io/ioutil"
	"regexp"
	"strings"
)

func DumpReadCloser(reader io.ReadCloser) ([]byte,io.ReadCloser, error) {
	var buf bytes.Buffer
	tee := io.TeeReader(reader, &buf)
	tmpRead1, tmpRead2 := ioutil.NopCloser(tee), ioutil.NopCloser(&buf)
	reader = tmpRead2
	data,err:=ioutil.ReadAll(tmpRead1)
	return data,reader,err
}

/**
	match restful api url
**/

// KeyMatch3 determines whether key1 matches the pattern of key2 (similar to RESTful path), key2 can contain a *.
// For example, "/foo/bar" matches "/foo/*", "/resource1" matches "/{resource}" ,"/foo" matches "/foo/*"
// key1 target
// key2 src
func KeyMatch3(key1 string, key2 string) bool {
	key2 = strings.Replace(key2, "/*", "/.*", -1)

	re := regexp.MustCompile(`\{[^/]+\}`)
	key2 = re.ReplaceAllString(key2, "$1[^/]+$2")
	if RegexMatch(key1HasVersion(key1)+"/", "^"+key2+"$") {
		return true
	}
	if RegexMatch(key1HasVersion(key1), "^"+key2+"$") {
		return true
	}
	if RegexMatch(key1+"/", "^"+key2+"$") {
		return true
	}
	return RegexMatch(key1, "^"+key2+"$")
}

func key1HasVersion(key1 string) string {
	vTags := []string{"v1", "v2", "v3", "v4"} //,"v5","v6","v7","v8","v9"
	for _, tag := range vTags {
		if index := strings.Index(key1, tag); index >= 0 {
			return key1[index+len(tag):]
		}
	}
	return key1
}

// KeyEqual  case key1='*' or key1=' '   result=true
func ActionEqual(key1 string, key2 string) bool {
	if key1 == "*" {
		return true
	}
	if len(key1) == 0 {
		return true
	}
	key1 = strings.ToLower(strings.TrimSpace(key1))
	key2 = strings.ToLower(strings.TrimSpace(key2))
	return strings.EqualFold(key1, key2)
}

// RegexMatch determines whether key1 matches the pattern of key2 in regular expression.
func RegexMatch(key1 string, key2 string) bool {
	res, err := regexp.MatchString(key2, key1)
	if err != nil {
		panic(err)
	}
	return res
}
