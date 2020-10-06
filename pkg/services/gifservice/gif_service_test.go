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

func TestMain(m *testing.M) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
}
func TestGetRandomByTagRestClientNoError(t *testing.T) {
	httpmock.RegisterResponder("GET", fmt.Sprintf(GiphyURI, "tag", os.Getenv("GIPHY_API_KEY")), httpmock.NewStringResponder(200,
		`{
			"title": "tag",
			"images": {
				"original": {
					"url": "https://delivery-much-challange-tag.gif"
				}	
			}
		}`,
	))
	result, err := gs.GetRandomByTag("tag")

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
}

func TestGetRandomByTagRestClientErrorUnmarshaling(t *testing.T) {
	httpmock.RegisterResponder("GET", fmt.Sprintf(GiphyURI, "tag", os.Getenv("GIPHY_API_KEY")), httpmock.NewStringResponder(200,
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
}

func TestGetRandomByTagRestClientErrorMountingRequest(t *testing.T) {
	GiphyURI = "fake.com.br"

	httpmock.RegisterResponder("GET", fmt.Sprintf(GiphyURI, "tag", os.Getenv("GIPHY_API_KEY")), httpmock.NewStringResponder(http.StatusInternalServerError, ""))

	result, err := gs.GetRandomByTag("")

	assert.NotNil(t, err)
	assert.Empty(t, result, "mounting request should got error")
}
