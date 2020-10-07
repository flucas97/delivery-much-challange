package gifservice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/flucas97/delivery-much-challange/internal/domain/gif"
	"github.com/flucas97/delivery-much-challange/tools/errortools"
)

var (
	// GiphyURI for Giphy API endpoint
	GiphyURI = "https://api.giphy.com/v1/gifs/random?tag=%s&api_key=%s&limit=1"
)

var (
	// GifService interface for other layers
	GifService gifServiceInterface = &gifService{}
)

type gifServiceInterface interface {
	GetRandomByTag(tag string) (*gif.Gif, *errortools.APIError)
}

type gifService struct{}

// GetRandomByTag is responsible for getting a Giphy Gif based on a tag
func (gs *gifService) GetRandomByTag(tag string) (*gif.Gif, *errortools.APIError) {
	client := &http.Client{}

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf(GiphyURI, tag, os.Getenv("GIPHY_API_KEY")), nil)
	if err != nil {
		return nil, errortools.APIErrorInterface.NewInternalServerError("error mounting request. gifservice.GetRandomByTag")
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, errortools.APIErrorInterface.NewInternalServerError("error doing request. gifservice.GetRandomByTag")
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errortools.APIErrorInterface.NewInternalServerError("error reading body. gifservice.GetRandomByTag")
	}

	var result gif.Gif

	if err := json.Unmarshal(bytes, &result); err != nil {
		return nil, errortools.APIErrorInterface.NewInternalServerError("error unmarshalling response from client. gifservice.GetRandomByTag")
	}

	return &result, nil
}
