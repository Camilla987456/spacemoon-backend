package product

import (
	"errors"
	"testing"
)

func TestProduct(t *testing.T) {
	var testName Name = "test-name"
	var testPrice Price = 0
	var testDescription Description = ""
	var err error
	var testProduct Product
	testProduct, err = New(testName, testPrice, testDescription)
	if err != nil {
		t.Fatalf("could not instantiate new product")
	}
	expectedDTO := Dto{
		Name:        testName,
		Price:       testPrice,
		Description: testDescription,
		Id:          testProduct.GetId(),
	}
	dto := testProduct.DTO()
	if dto != expectedDTO {
		t.Fatalf("invalid DTO retrieved, expected '%+v' obtained '%+v'\n", expectedDTO, dto)
	}
	if testProduct.GetName() != testName {
		t.Fatalf("invalid name, expected '%s' obtained '%s'\n", testName, testProduct.GetName())
	}
	if testProduct.GetPrice() != testPrice {
		t.Fatalf("invalid price, expected '%d' obtained '%d'\n", testPrice, testProduct.GetPrice())
	}
	if testProduct.GetDescription() != testDescription {
		t.Fatalf("invalid description, expected '%s' obtained '%s'\n", testDescription, testProduct.GetDescription())
	}
}

func TestNew_ThrowsErrorOnEmptyName(t *testing.T) {
	var emptyName Name = ""
	var testPrice Price = 100
	var testDescription Description = "this is a test product"
	_, err := New(emptyName, testPrice, testDescription)
	if err == nil {
		t.Fatalf("empty name not detected")
	}
	if !errors.Is(err, EmptyNameError) {
		t.Fatalf("empty name error not thrown")
	}
}
