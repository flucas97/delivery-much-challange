package gifservice

import (
	"os"

	"github.com/flucas97/delivery-much-challange/internal/domain/gif"
	"github.com/flucas97/delivery-much-challange/tools/errortools"
	libgiphy "github.com/sanzaru/go-giphy"
)

var (
	// GiphyURL for Giphy API endpoint
	GiphyURL = "https://api.giphy.com/v1/gifs/random?tag=%s&api_key=%s"
)

var (
	// GifService interface for other layers
	GifService gifServiceInterface = &gifService{}
)

type gifServiceInterface interface {
	GetRandom(tag string) (*gif.Gif, *errortools.APIError)
}

type gifService struct{}

func (gs *gifService) GetRandom(tag string) (*gif.Gif, *errortools.APIError) {
	giphyClient := libgiphy.NewGiphy(os.Getenv("GIPHY_API_KEY"))

	random, err := giphyClient.GetRandom(tag)
	if err != nil {
		return nil, errortools.APIErrorInterface.NewInternalServerError("error getting Giphy. gifservice.GetRandom")
	}

	var result = &gif.Gif{
		URL: random.Data.Image_original_url,
	}

	return result, nil
}
