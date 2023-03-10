package main

import (
	"context"
	"errors"
	"log"
	"os"
	"spacemoon/login"
	"spacemoon/product"
	"spacemoon/product/category"
	"spacemoon/product/ratings"
	"spacemoon/server/firestore"
	"strings"
	"time"
)

func getLoginPersistence() login.Persistence {
	creds := os.Getenv(googleCredentials)
	if strings.TrimSpace(creds) == "" {
		log.Default().Print("using temporary (dev) login persistence - no google credentials file set")
		return makeTemporaryLoginPersistence()
	}
	per, err := firestore.GetPersistence(context.Background())
	if err != nil {
		log.Default().Printf("using temporary (dev) login persistence - error getting firestore persistence: %s", err.Error())
		return makeTemporaryLoginPersistence()
	}
	return per
}

func makeTemporaryLoginPersistence() login.Persistence {
	return &temporaryLoginPersistence{
		users:  Credentials{},
		tokens: login.Tokens{},
	}
}

func getProductPersistence() product.Persistence {
	creds := os.Getenv(googleCredentials)
	if strings.TrimSpace(creds) == "" {
		log.Default().Print("using temporary (dev) product persistence - no google credentials file set")
		return &temporaryProductPersistence{savedProducts: map[product.Id]product.Dto{}}
	}
	per, err := firestore.GetPersistence(context.Background())
	if err != nil {
		log.Default().Printf("using temporary (dev) product persistence - error getting firestore persistence: %s", err.Error())
		return &temporaryProductPersistence{savedProducts: map[product.Id]product.Dto{}}
	}
	return per
}

func getProductRatingsPersistence() ratings.Persistence {
	return &temporaryRatingsPersistence{r: map[product.Id]ratings.Rating{}}
}

func getCategoryPersistence() category.Persistence {
	return &temporaryCategoryPersistence{categories: map[category.Name]category.DTO{}}
}

type temporaryProductPersistence struct {
	savedProducts product.Products
}

func (t *temporaryProductPersistence) DeleteProduct(id product.Id) error {
	delete(t.savedProducts, id)
	return nil
}

func (t *temporaryProductPersistence) GetProducts() (product.Products, error) {
	return t.savedProducts, nil
}

func (t *temporaryProductPersistence) SaveProduct(p product.Product) error {
	if t.savedProducts == nil {
		t.savedProducts = make(product.Products)
	}
	t.savedProducts[p.GetId()] = p.DTO()
	return nil
}

type temporaryCategoryPersistence struct {
	categories category.Categories
}

func (t *temporaryCategoryPersistence) DeleteCategory(name category.Name) {
	delete(t.categories, name)
}

func (t *temporaryCategoryPersistence) SaveCategory(dto category.DTO) {
	if t.categories == nil {
		t.categories = make(category.Categories)
	}
	t.categories[dto.Name] = dto
}

func (t *temporaryCategoryPersistence) GetCategories() category.Categories {
	return t.categories
}

type temporaryLoginPersistence struct {
	users  Credentials
	tokens login.Tokens
}

func (t *temporaryLoginPersistence) DeleteUser(name login.UserName) error {
	//TODO implement me
	panic("implement me")
}

func (t *temporaryLoginPersistence) ValidateCredentials(usr login.UserName, p login.Password) bool {
	if t.users[usr] == p {
		return true
	}
	return false
}

func (t *temporaryLoginPersistence) GetUser(token login.Token) (login.UserName, error) {
	tokenInfo, exists := t.tokens[token]
	if !exists {
		return "", errors.New("token not found")
	}
	if tokenInfo.Expiration.Before(time.Now()) {
		delete(t.tokens, token)
		return "", errors.New("token expired, deleted")
	}
	return tokenInfo.User, nil
}

func (t *temporaryLoginPersistence) SetUserToken(user login.UserName, token login.Token, tokenDuration time.Duration) error {
	if t.tokens == nil {
		t.tokens = make(login.Tokens)
	}
	t.tokens[token] = login.TokenDetails{
		User:       user,
		Expiration: time.Now().Add(tokenDuration),
	}
	return nil
}

func (t *temporaryLoginPersistence) SignUpUser(u login.UserName, p login.Password) error {
	if t.users == nil {
		t.users = make(Credentials)
	}
	t.users[u] = p
	return nil
}

type temporaryRatingsPersistence struct {
	r ratings.Ratings
}

func (t *temporaryRatingsPersistence) ReadRating(id product.Id) ratings.Rating {
	return t.r[id]
}

func (t *temporaryRatingsPersistence) SaveRating(id product.Id, rating ratings.Rating) {
	if t.r == nil {
		t.r = make(ratings.Ratings)
	}
	t.r[id] = rating
}

type Credentials map[login.UserName]login.Password

const googleCredentials = "GOOGLE_APPLICATION_CREDENTIALS"
