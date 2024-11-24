package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOpenPQ(t *testing.T) {
	db := OpenPQ()
	assert.NotNil(t, db)

}
