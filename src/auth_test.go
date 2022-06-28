package core

//var a = Auth{
//	Username: "username",
//	Password: "password",
//}
//
//func TestString(t *testing.T) {
//	assert.Equal(t, a.String(), "username:password")
//}
//
//func TestGetAccessToken(t *testing.T) {
//	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
//	// pass 'nil' as the third parameter.
//	req, err := http.NewRequest("GET", "/auth/access-token", nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	req.Header.Set("Authorization", "Basic something:else")
//
//	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
//	rr := httptest.NewRecorder()
//	handler := http.HandlerFunc(GetAccessToken)
//
//	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
//	// directly and pass in our Request and ResponseRecorder.
//	handler.ServeHTTP(rr, req)
//
//	//// Check the status code is what we expect.
//	if status := rr.Code; status != http.StatusOK {
//		t.Errorf("handler returned wrong status code: got %v want %v",
//			status, http.StatusOK)
//	}
//
//	//// Check the response body is what we expect.
//	var parsedResponse map[string]interface{}
//	err = json.Unmarshal([]byte(rr.Body.String()), &parsedResponse)
//	if err != nil {
//		return
//	}
//
//	if parsedResponse["Scope"] != "*" {
//		t.Errorf("handler returned unexpected body: got %v want %v", parsedResponse["Scope"], "*")
//	}
//}
//
//func TestGetAuthParamsFromHeader(t *testing.T) {
//	req, err := http.NewRequest("GET", "/auth/access-token", nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	req.Header.Set("Authorization", "Basic something:else")
//	auth, err := getAuthFromHeader(req)
//
//	assert.Nil(t, err)
//	assert.Equal(t, auth.String(), "something:else")
//}
