package movie

import (
	"context"
	"errors"
	modelmetadata "microgomovies/metadata/pkg"
	"microgomovies/movie/internal/gateway"
	modelmovie "microgomovies/movie/pkg"
	model "microgomovies/rating/pkg"
)

var ErrNotFound = errors.New("movie metadata not found")

type ratingGateway interface {
	GetAggregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error)
	PutRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error
}

type metadataGateway interface {
	Get(ctx context.Context, id string) (*modelmetadata.Metadata, error)
}

type Controller struct {
	metadataGateway metadataGateway
	ratingGateway   ratingGateway
}

func New(ratingGateway ratingGateway, metadataGateway metadataGateway) *Controller {
	return &Controller{ratingGateway: ratingGateway, metadataGateway: metadataGateway}
}

// Returns the movie details including aggregated rating & movie metadata
func (c *Controller) Get(ctx context.Context, id string) (*modelmovie.MovieDetails, error) {
	meta, err := c.metadataGateway.Get(ctx, id)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	details := &modelmovie.MovieDetails{Metadata: *meta}
	recordID := model.RecordID(id)
	rating, err := c.ratingGateway.GetAggregatedRating(ctx, recordID, model.RecordTypeMovie)
	if err != nil && !errors.Is(err, gateway.ErrNotFound) {
		// Non error if no rating
	} else if err != nil {
		return nil, err
	}

	details.Rating = &rating
	return details, nil
}
