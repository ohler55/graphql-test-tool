package main_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/graphql-test-tool/gtt"
)

func gttTest(t *testing.T, filepath string) {
	uc, err := gtt.NewUseCase(filepath)
	if err != nil {
		t.Fatal(err.Error())
	}
	r := gtt.Runner{
		Server:   fmt.Sprintf("http://localhost:%d", testPort),
		Base:     "/graphql",
		Indent:   2,
		UseCases: []*gtt.UseCase{uc},
	}
	if testing.Verbose() {
		r.ShowComments = true
		r.ShowResponses = true
		r.ShowRequests = true
	}

	if err = r.Run(); err != nil {
		t.Fatal(err.Error())
	}
}

func TestTypes(t *testing.T) {
	gttTest(t, "gtt/types.json")
}
