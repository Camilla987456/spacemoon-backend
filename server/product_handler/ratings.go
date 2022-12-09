package product_handler

import (
	"fmt"
	"net/http"
	"spacemoon/product"
	"spacemoon/product/ratings"
	"strconv"
)

func MakeRankingsHandler() http.Handler {
	return &rankingsHandler{rater: ratings.NewProductRater()}
}

type rankingsHandler struct {
	rater ratings.ProductRater
}

func (rh *rankingsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	productId := r.FormValue("id")
	switch r.Method {
	case http.MethodGet:
		rating := rh.rater.GetRating(product.Id(productId))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fmt.Sprintf("Rating: %d", rating)))
	case http.MethodPost:
		ratingStr := r.FormValue("rating")
		rating, err := strconv.ParseInt(ratingStr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("could not post rating: " + err.Error()))
		}
		w.WriteHeader(http.StatusOK)
		rh.rater.AddRating(product.Id(productId), ratings.Rating(rating))
		_, _ = w.Write([]byte(fmt.Sprintf("Rating: %d", rh.rater.GetRating(product.Id(productId)))))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
