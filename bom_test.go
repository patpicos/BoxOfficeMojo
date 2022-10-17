package boxofficemojo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchTopGunMaverick(t *testing.T) {
	t.Parallel()
	bom, err := Search("tt1745960")
	assert.Nil(t, err)
	fmt.Println(bom)
	assert.Equal(t, 4751, bom.UsaScreens)
	// fmt.Printf("%+v", bom)
}

// https://www.imdb.com/title/tt1498870
// Straight to Video (STV) and should fail returning screens
func TestSearch40yrOldVirgin(t *testing.T) {
	t.Parallel()
	bom, err := Search("tt1498870")
	assert.NotNil(t, err)
	assert.Equal(t, 0, bom.UsaScreens)
	// fmt.Printf("%+v", bom)
}
