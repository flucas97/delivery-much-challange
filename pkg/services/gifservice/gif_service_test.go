package gifservice

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/flucas97/delivery-much-challange/internal/domain/gif"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

var gs gifService

func TestGetRandomByTag(t *testing.T) {
	t.Run("Giphy without original URL", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", fmt.Sprintf(GiphyURI, "mychallange", os.Getenv("GIPHY_API_KEY")), httpmock.NewStringResponder(http.StatusOK,
			`{
				"title": "",
				"images": {
					"original": {
						"url": ""
					}	
				}
			}`,
		))

		result, err := gs.GetRandomByTag("mychallange")
		assert.NotNil(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "no response from client. gifservice.GetRandomByTag", err.Message)
	})

	t.Run("Success receiving a Giphy", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET", fmt.Sprintf(GiphyURI, "tag", os.Getenv("GIPHY_API_KEY")), httpmock.NewStringResponder(http.StatusOK,
			`{
				"title": "tag",
				"images": {
					"original": {
						"url": "https://delivery-much-challange-tag.gif"
					}	
				}
			}`,
		))

		result, err := GifService.GetRandomByTag("tag")

		expected := &gif.Gif{
			Title: "tag",
			Images: gif.Images{
				Original: gif.Original{
					URL: "https://delivery-much-challange-tag.gif",
				},
			},
		}

		assert.Nil(t, err)
		assert.Equal(t, expected, result)
	})
	t.Run("Error unmarshalling response from client", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET", fmt.Sprintf(GiphyURI, "tag", os.Getenv("GIPHY_API_KEY")), httpmock.NewStringResponder(http.StatusOK,
			`{
				"title": "tag",
				"images": {
					"original": {
						"url": "https://delivery-much-challange-tag.gif"
					},
				}
			}`,
		))

		result, err := gs.GetRandomByTag("tag")

		assert.NotNil(t, err)
		assert.Empty(t, result, "unmarshal should got error")
	})

	t.Run("Error mounting request", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		GiphyURI = "fake.com.br"

		httpmock.RegisterResponder("GET", fmt.Sprintf(GiphyURI, "tag", os.Getenv("GIPHY_API_KEY")), httpmock.NewStringResponder(http.StatusInternalServerError, ""))

		result, err := gs.GetRandomByTag("tag")

		assert.Equal(t, "error mounting request. gifservice.GetRandomByTag", err.Message)
		assert.NotNil(t, err)
		assert.Empty(t, result, "mounting request should got error")
	})
}
