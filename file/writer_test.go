package file

import (
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {

	var (
		assert  = assert.New(t)
		require = require.New(t)
	)

	id := make([]byte, 16)
	io.ReadFull(rand.Reader, id)

	dir := fmt.Sprintf("/tmp/alm-%x", id)

	err := os.MkdirAll(dir, 0777)
	require.NoError(err)

	for i := 0; i < 2; i++ {

		w := New(dir)

		wc, err := w.Add("foo", "bar", "wee")
		assert.NoError(err)

		fmt.Fprintln(wc, "line 1")
		fmt.Fprintln(wc, "line 2")

		wc.Close()

		out, err := ioutil.ReadFile(dir + "/foo-bar-wee.log")
		assert.NoError(err)
		assert.Equal("line 1\nline 2\n", string(out))

	}

}
