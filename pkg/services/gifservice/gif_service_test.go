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
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	GiphyURI = "fake.com.br"

	httpmock.RegisterResponder("GET", fmt.Sprintf(GiphyURI, "tag", os.Getenv("GIPHY_API_KEY")), httpmock.NewStringResponder(http.StatusInternalServerError, ""))

	result, err := gs.GetRandomByTag("tag")

	assert.NotNil(t, err)
	assert.Empty(t, result, "mounting request should got error")
}

func TestGetRandomByTagRestClientErrorDoingRequest(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", fmt.Sprintf(GiphyURI, "tag2", os.Getenv("GIPHY_API_KEY")), httpmock.NewBytesResponder(http.StatusInternalServerError, []byte("")))

	result, err := GifService.GetRandomByTag("tag")
	expectedError := errortools.NewInternalServerError("error doing request. gifservice.GetRandomByTag")

	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, http.StatusInternalServerError, err.Status)
	assert.IsType(t, expectedError, err, "doing request should got same error type")
}
