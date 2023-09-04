package src

import (
	"bytes"
	"fmt"
	"testing"
)

func TestFoo(t *testing.T) {
	equal := bytes.Equal([]byte("3bbbb90c752898472ceb30317b538172"), []byte("3bbbb90c752898472ceb30317b538172"))
	fmt.Println(equal)
	b := "3bbbb90c752898472ceb30317b538172" == "3bbbb90c752898472ceb30317b538172"
	fmt.Println(b)
}
