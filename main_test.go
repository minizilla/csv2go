package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/minizilla/testr"
)

func TestCSV2Go(t *testing.T) {
	fsrc := testr.MustV(os.Open("./testdata/model.go"))
	fdst := testr.MustV(os.Open("./testdata/model.csv.go"))
	var want, got bytes.Buffer
	_ = testr.MustV(want.ReadFrom(fdst))

	testr.Must(run(fsrc, &got))

	assert := testr.New(t)
	res := bytes.Compare(want.Bytes(), got.Bytes())
	assert.Equal(res, 0)
}
