package xenditv3

import (
	"encoding/json"
	"testing"

	"github.com/cheekybits/is"
)

func TestCustomerCreateRequestJSON(t *testing.T) {
	is := is.New(t)

	req := CustomerCreateRequest{
		ReferenceID:  "1234567890123",
		Type:         "INDIVIDUAL",
		Email:        "buyer@example.com",
		MobileNumber: "+6281234567890",
		GivenNames:   "Reyvin",
		Surname:      "Test",
	}
	b, err := json.Marshal(req)
	is.NoErr(err)

	var m map[string]interface{}
	is.NoErr(json.Unmarshal(b, &m))
	is.Equal("Reyvin", m["given_names"])
	is.Equal("Test", m["surname"])
	_, hasNested := m["individual_detail"]
	is.True(!hasNested)
}

func TestParseCustomerList(t *testing.T) {
	is := is.New(t)

	wrapped := []byte(`{"data":[{"id":"abc","reference_id":"buyer-1","email":"a@b.c"}]}`)
	list, err := parseCustomerList(wrapped)
	is.NoErr(err)
	is.Equal(1, len(list))
	is.Equal("abc", list[0].ID)

	plain := []byte(`[{"id":"def","reference_id":"buyer-2"}]`)
	list, err = parseCustomerList(plain)
	is.NoErr(err)
	is.Equal(1, len(list))
	is.Equal("def", list[0].ID)

	single := []byte(`{"id":"ghi","reference_id":"buyer-3"}`)
	list, err = parseCustomerList(single)
	is.NoErr(err)
	is.Equal(1, len(list))
	is.Equal("ghi", list[0].ID)
}

func TestCustomerURLs(t *testing.T) {
	is := is.New(t)
	client := Client{BaseURL: DefaultBaseURL}

	is.Equal(
		"https://api.xendit.co/customers?reference_id=1234567890123",
		getCustomersByReferenceIDURL(client, "1234567890123"),
	)
	is.Equal("https://api.xendit.co/customers", createCustomerURL(client))
}

func TestErrorResponse(t *testing.T) {
	is := is.New(t)

	e := &ErrorResponse{ErrorCode: "INVALID", ErrorMessage: "bad request"}
	e.markHTTPError(400)
	is.True(e.ErrorStatus)
	is.Equal("[INVALID] bad request", e.Error())

	e.markHTTPError(201)
	is.True(!e.ErrorStatus)
}
