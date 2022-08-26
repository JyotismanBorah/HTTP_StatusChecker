package handler

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReqHandler(t *testing.T) {
	//create dummy request
	req, err := http.NewRequest("GET", "/websites", nil)
	if err != nil {
		log.Fatal("Error in creating http request: ", err)
	}

	//create ResponseRecorder to satisfy ResponseWriter
	rr := httptest.NewRecorder()
	handlr := http.HandlerFunc(ReqHandler)

	handlr.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler failed in GET method: got %v, want %v\n", status, http.StatusOK)
	}
	expectedBody := `Website list is empty, please add URLs using POST method`
	if rr.Body.String() != expectedBody {
		t.Errorf("Handler returned unexpected body: got %v, want %v\n", rr.Body.String(), expectedBody)
	}

	fmt.Println("GET method tested")

	jsonBody := []byte(`{
		"websites": [
			"http://www.google.com",
			"http://www.fakewebsite1.com"
		]
	}`)
	bodyReader := bytes.NewReader(jsonBody)

	req, err = http.NewRequest("POST", "/websites", bodyReader)
	if err != nil {
		log.Fatal("Error in creating http request: ", err)
	}

	//create ResponseRecorder to satisfy ResponseWriter
	rr = httptest.NewRecorder()
	handlr = http.HandlerFunc(ReqHandler)

	handlr.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler failed in POST method: got %v, want %v\n", status, http.StatusOK)
	}

	fmt.Println("POST method tested")
}

func TestUpdateStatusUtil(t *testing.T) {
	H := httpChecker{}
	key1 := "http://www.google.com"
	key2 := "http://www.fakewebsite1.com"
	status1 := H.Check(context.Background(), key1)
	status2 := H.Check(context.Background(), key2)
	assert.Equal(t, status1, true)
	assert.Equal(t, status2, false)
	fmt.Println("UpdateStatusUtil function tested")
}
