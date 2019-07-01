package xenditgo_test

import (
	"testing"

	"github.com/cheekybits/is"
	xenditgo "github.com/grosenia/xendit-go-client"
)

func TestDefaultEnvironmentType(t *testing.T) {
	is := is.New(t)

	client := xenditgo.NewClient()
	is.Equal(xenditgo.Sandbox, client.APIEnvType)

}
