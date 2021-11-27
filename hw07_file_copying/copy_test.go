package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("file not found case", func(t *testing.T) {
		err := Copy("testdata/wrong_file", "testdata/dgad.txt", 1000000, 0)
		require.Error(t, ErrUnsupportedFile, err)
		require.Equal(t, ErrUnsupportedFile, err)
	})

	t.Run("oversize offset case", func(t *testing.T) {
		err := Copy("testdata/input.txt", "testdata/dgad.txt", 1000000, 0)
		require.Error(t, ErrOffsetExceedsFileSize, err)
		require.Equal(t, ErrOffsetExceedsFileSize, err)
	})

	t.Run("oversize limit case", func(t *testing.T) {
		err := Copy("testdata/input.txt", "testdata/dgad.txt", 0, 10)
		require.Error(t, ErrLimitExcess, err)
		require.Equal(t, ErrLimitExcess, err)
	})

	t.Run("file copied case", func(t *testing.T) {
		_ = Copy("testdata/out_offset0_limit10.txt", "testdata/dgad.txt", 0, 0)
		require.FileExistsf(t, "testdata/dgad.txt", "File not copied")
	})

	_ = os.Remove("testdata/dgad.txt")
}
