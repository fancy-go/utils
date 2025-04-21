package file_util

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"
)

func TestTarGzDir(t *testing.T) {
	if !assert.NoError(t, TarGzDir("/home/peterq/dev/projects/1s/osec-spider-go/scripts", "/home/peterq/dev/projects/1s/osec-spider-go/temp/test.tar.gz")) {
		return
	}
	dir := "/home/peterq/dev/projects/1s/osec-spider-go/temp/untartest_" + time.Now().Format("20060102150405")
	if !assert.NoError(t, os.MkdirAll(dir, os.ModePerm)) {
		return
	}
	rootFiles, err := UnTarGzFile("/home/peterq/dev/projects/1s/osec-spider-go/temp/test.tar.gz", dir)
	if !assert.NoError(t, err) {
		return
	}
	log.Println("root files:", rootFiles)
}
