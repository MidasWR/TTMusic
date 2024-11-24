package server

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetInfoFomAPI(t *testing.T) {
	url, err := getYouTubeVideoURL("Natural", "Imagine Dragons")
	assert.Nil(t, err)
	fmt.Println(url)
}
