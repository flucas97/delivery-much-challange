package gifservice

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/flucas97/delivery-much-challange/internal/domain/gif"
	"github.com/flucas97/delivery-much-challange/tools/errortools"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

var gs gifService

func TestGetRandomByTagRestClientNoError(t *testing.T) {
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
}

func TestGetRandomByTagRestClientErrorUnmarshaling(t *testing.T) {
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
}

func TestGetRandomByTagRestClientErrorMountingRequest(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	GiphyURI = "fake.com.br"

	httpmock.RegisterResponder("GET", fmt.Sprintf(GiphyURI, "tag", os.Getenv("GIPHY_API_KEY")), httpmock.NewStringResponder(http.StatusInternalServerError, ""))

	result, err := gs.GetRandomByTag("tag")

	assert.Equal(t, "error mounting request. gifservice.GetRandomByTag", err.Message)
	assert.NotNil(t, err)
	assert.Empty(t, result, "mounting request should got error")
}

func TestGetRandomByTagRestClientErrorDoingRequest(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	e := errortools.APIErrorInterface.NewInternalServerError("error doing request. gifservice.GetRandomByTag")

	httpmock.RegisterResponder("GET", fmt.Sprintf(GiphyURI, " ", os.Getenv("GIPHY_API_KEY")), httpmock.NewErrorResponder(e))

	result, err := gs.GetRandomByTag(" ")

	assert.Equal(t, e.Message, err.Message)
	assert.NotNil(t, err)
	assert.Empty(t, result, "doing request should got error")
}
