package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ardanlabs/service/app/sales-api/handlers"
	"github.com/ardanlabs/service/business/data/product"
	"github.com/ardanlabs/service/business/tests"
	"github.com/google/go-cmp/cmp"
)

// ProductTests holds methods for each product subtest. This type allows
// passing dependencies for tests while still providing a convenient syntax
// when subtests are registered.
type ProductTests struct {
	app       http.Handler
	userToken string
}

func TestProducts(t *testing.T) {
	test := tests.NewIntegration(t)
	t.Cleanup(test.Teardown)

	shutdown := make(chan os.Signal, 1)
	tests := ProductTests{
		app:       handlers.API("develop", shutdown, test.Log, test.Auth, test.DB),
		userToken: test.Token(),
	}

	t.Run("crudProducts", tests.crudProduct)
}

// crudProduct performs a complete test of CRUD against the api.
func (pt *ProductTests) crudProduct(t *testing.T) {
	pt.postProduct201(t)
	// defer pt.deleteProduct204(t, p.ID)

	// pt.getProduct200(t, p.ID)
	// pt.putProduct204(t, p.ID)
}

// postProduct201 validates a product can be created with the endpoint.
func (pt *ProductTests) postProduct201(t *testing.T) product.Product {
	np := product.NewProduct{
		Name:     "Comic Books",
		Cost:     25,
		Quantity: 60,
	}

	body, err := json.Marshal(&np)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	var got product.Product
	t.Log("Given the need to create a new product with the products endpoint.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the declared product value.", testID)
		{
			if w.Code != http.StatusCreated {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 201 for the response : %v", tests.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 201 for the response.", tests.Success, testID)

			if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to unmarshal the response : %v", tests.Failed, testID, err)
			}

			// Define what we wanted to receive. We will just trust the generated
			// fields like ID and Dates so we copy p.
			exp := got
			exp.Name = "Comic Books"
			exp.Cost = 25
			exp.Quantity = 60

			if diff := cmp.Diff(got, exp); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get the expected result. Diff:\n%s", tests.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get the expected result.", tests.Success, testID)
		}
	}

	return got
}
