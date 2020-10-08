package gifservice

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

var (
	gs gifService
)

func TestGetRandom(t *testing.T) {
	t.Run("Error receiving a Giphy", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET", fmt.Sprintf("http://api.giphy.com/v1/gifs/random?api_key=%s", os.Getenv("GIPHY_API_KEY")), httpmock.NewStringResponder(http.StatusInternalServerError, `{error}`))

		gif, err := gs.GetRandom("")
		assert.Equal(t, "error getting Giphy. gifservice.GetRandom", err.Message)
		assert.Nil(t, gif)
	})
}
