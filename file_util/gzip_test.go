package file_util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTarGzDir(t *testing.T) {
	if !assert.NoError(t, TarGzDir("testdata", "testdata.tar.gz")) {
		return
	}
}
