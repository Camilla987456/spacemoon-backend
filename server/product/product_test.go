package product

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"spacemoon/product"
	"testing"
)

func TestHandler_Get(t *testing.T) {
	var fakePersistence Persistence = stubPersistence{}
	testHandler := MakeHandler(fakePersistence)
	fakeRequest := httptest.NewRequest(http.MethodGet, "/product", http.NoBody)
	spy := spyWriter{}
	testHandler.ServeHTTP(&spy, fakeRequest)
	fmt.Printf("%s", spy.written)
}

type spyWriter struct {
	written string
}

func (s *spyWriter) Header() http.Header {
	//TODO implement me
	panic("implement me")
}

func (s *spyWriter) Write(bytes []byte) (int, error) {
	s.written = s.written + fmt.Sprintf("%s", bytes)
	return len(bytes), nil
}

func (s *spyWriter) WriteHeader(_ int) {
	//TODO implement me
	panic("implement me")
}

type stubPersistence struct {
}

func (s stubPersistence) GetProducts() product.Products {
	return expectedProducts
}

var expectedProducts = make(product.Products)

func init() {
	expectedProducts["product1-id"] = product.Dto{
		Name:        "product1",
		Price:       1,
		Description: "",
	}
	expectedProducts["product2-id"] = product.Dto{
		Name:        "product2",
		Price:       10,
		Description: "",
	}
}
