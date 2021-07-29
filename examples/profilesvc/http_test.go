// Code generated by kok; DO NOT EDIT.
// github.com/RussellLuo/kok

package profilesvc

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

// Ensure that ServiceMock does implement Service.
var _ Service = &ServiceMock{}

type ServiceMock struct {
	DeleteAddressFunc func(ctx context.Context, id string, addressID string) (err error)
	DeleteProfileFunc func(ctx context.Context, id string) (err error)
	GetAddressFunc    func(ctx context.Context, id string, addressID string) (address Address, err error)
	GetAddressesFunc  func(ctx context.Context, id string) (addresses []Address, err error)
	GetProfileFunc    func(ctx context.Context, id string) (profile Profile, err error)
	PatchProfileFunc  func(ctx context.Context, id string, profile Profile) (err error)
	PostAddressFunc   func(ctx context.Context, id string, address Address) (err error)
	PostProfileFunc   func(ctx context.Context, profile Profile) (err error)
	PutProfileFunc    func(ctx context.Context, id string, profile Profile) (err error)
}

func (mock *ServiceMock) DeleteAddress(ctx context.Context, id string, addressID string) (err error) {
	if mock.DeleteAddressFunc == nil {
		panic("ServiceMock.DeleteAddressFunc: not implemented")
	}
	return mock.DeleteAddressFunc(ctx, id, addressID)
}

func (mock *ServiceMock) DeleteProfile(ctx context.Context, id string) (err error) {
	if mock.DeleteProfileFunc == nil {
		panic("ServiceMock.DeleteProfileFunc: not implemented")
	}
	return mock.DeleteProfileFunc(ctx, id)
}

func (mock *ServiceMock) GetAddress(ctx context.Context, id string, addressID string) (address Address, err error) {
	if mock.GetAddressFunc == nil {
		panic("ServiceMock.GetAddressFunc: not implemented")
	}
	return mock.GetAddressFunc(ctx, id, addressID)
}

func (mock *ServiceMock) GetAddresses(ctx context.Context, id string) (addresses []Address, err error) {
	if mock.GetAddressesFunc == nil {
		panic("ServiceMock.GetAddressesFunc: not implemented")
	}
	return mock.GetAddressesFunc(ctx, id)
}

func (mock *ServiceMock) GetProfile(ctx context.Context, id string) (profile Profile, err error) {
	if mock.GetProfileFunc == nil {
		panic("ServiceMock.GetProfileFunc: not implemented")
	}
	return mock.GetProfileFunc(ctx, id)
}

func (mock *ServiceMock) PatchProfile(ctx context.Context, id string, profile Profile) (err error) {
	if mock.PatchProfileFunc == nil {
		panic("ServiceMock.PatchProfileFunc: not implemented")
	}
	return mock.PatchProfileFunc(ctx, id, profile)
}

func (mock *ServiceMock) PostAddress(ctx context.Context, id string, address Address) (err error) {
	if mock.PostAddressFunc == nil {
		panic("ServiceMock.PostAddressFunc: not implemented")
	}
	return mock.PostAddressFunc(ctx, id, address)
}

func (mock *ServiceMock) PostProfile(ctx context.Context, profile Profile) (err error) {
	if mock.PostProfileFunc == nil {
		panic("ServiceMock.PostProfileFunc: not implemented")
	}
	return mock.PostProfileFunc(ctx, profile)
}

func (mock *ServiceMock) PutProfile(ctx context.Context, id string, profile Profile) (err error) {
	if mock.PutProfileFunc == nil {
		panic("ServiceMock.PutProfileFunc: not implemented")
	}
	return mock.PutProfileFunc(ctx, id, profile)
}

type request struct {
	method string
	path   string
	header map[string]string
	body   string
}

func (r request) ServedBy(handler http.Handler) *httptest.ResponseRecorder {
	var req *http.Request
	if r.body != "" {
		reqBody := strings.NewReader(r.body)
		req = httptest.NewRequest(r.method, r.path, reqBody)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	} else {
		req = httptest.NewRequest(r.method, r.path, nil)
	}

	for key, value := range r.header {
		req.Header.Set(key, value)
	}

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	return w
}

type response struct {
	statusCode  int
	contentType string
	body        []byte
}

func (want response) Equal(w *httptest.ResponseRecorder) string {
	resp := w.Result()
	gotBody, _ := ioutil.ReadAll(resp.Body)

	gotStatusCode := resp.StatusCode
	if gotStatusCode != want.statusCode {
		return fmt.Sprintf("StatusCode: got (%d), want (%d)", gotStatusCode, want.statusCode)
	}

	wantContentType := want.contentType
	if wantContentType == "" {
		wantContentType = "application/json; charset=utf-8"
	}

	gotContentType := resp.Header.Get("Content-Type")
	if gotContentType != wantContentType {
		return fmt.Sprintf("ContentType: got (%q), want (%q)", gotContentType, wantContentType)
	}

	if !bytes.Equal(gotBody, want.body) {
		return fmt.Sprintf("Body: got (%q), want (%q)", gotBody, want.body)
	}

	return ""
}

func TestHTTP_PostProfile(t *testing.T) {
	// in contains all the input parameters (except ctx) of PostProfile.
	type in struct {
		profile Profile
	}

	// out contains all the output parameters of PostProfile.
	type out struct {
		err error
	}

	cases := []struct {
		name         string
		request      request
		wantIn       in
		out          out
		wantResponse response
	}{
		{
			name: "ok",
			request: request{
				method: "POST",
				path:   "/profiles",
				body:   `{"profile": {"id": "1234", "name": "kok"}}`,
			},
			wantIn: in{
				profile: Profile{
					ID:   "1234",
					Name: "kok",
				},
			},
			out: out{
				err: nil,
			},
			wantResponse: response{
				statusCode: http.StatusOK,
				body:       []byte(`{}` + "\n"),
			},
		},
		{
			name: "err",
			request: request{
				method: "POST",
				path:   "/profiles",
				body:   `{"profile": {"id": "1234", "name": "kok"}}`,
			},
			wantIn: in{
				profile: Profile{
					ID:   "1234",
					Name: "kok",
				},
			},
			out: out{
				err: ErrAlreadyExists,
			},
			wantResponse: response{
				statusCode: http.StatusBadRequest,
				body:       []byte(`{"error":"already exists"}` + "\n"),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var gotIn in
			w := c.request.ServedBy(NewHTTPRouter(
				&ServiceMock{
					PostProfileFunc: func(ctx context.Context, profile Profile) (err error) {
						gotIn = in{
							profile: profile,
						}
						return c.out.err
					},
				},
				NewCodecs(),
			))

			if !reflect.DeepEqual(gotIn, c.wantIn) {
				t.Fatalf("In: got (%v), want (%v)", gotIn, c.wantIn)
			}

			if errStr := c.wantResponse.Equal(w); errStr != "" {
				t.Fatal(errStr)
			}
		})
	}
}

func TestHTTP_GetProfile(t *testing.T) {
	// in contains all the input parameters (except ctx) of GetProfile.
	type in struct {
		id string
	}

	// out contains all the output parameters of GetProfile.
	type out struct {
		profile Profile
		err     error
	}

	cases := []struct {
		name         string
		request      request
		wantIn       in
		out          out
		wantResponse response
	}{
		{
			name: "ok",
			request: request{
				method: "GET",
				path:   "/profiles/1234",
			},
			wantIn: in{
				id: "1234",
			},
			out: out{
				profile: Profile{
					ID:   "1234",
					Name: "kok",
				},
				err: nil,
			},
			wantResponse: response{
				statusCode: http.StatusOK,
				body:       []byte(`{"profile":{"id":"1234","name":"kok"}}` + "\n"),
			},
		},
		{
			name: "err",
			request: request{
				method: "GET",
				path:   "/profiles/5678",
			},
			wantIn: in{
				id: "5678",
			},
			out: out{
				profile: Profile{},
				err:     ErrNotFound,
			},
			wantResponse: response{
				statusCode: http.StatusNotFound,
				body:       []byte(`{"error":"not found"}` + "\n"),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var gotIn in
			w := c.request.ServedBy(NewHTTPRouter(
				&ServiceMock{
					GetProfileFunc: func(ctx context.Context, id string) (profile Profile, err error) {
						gotIn = in{
							id: id,
						}
						return c.out.profile, c.out.err
					},
				},
				NewCodecs(),
			))

			if !reflect.DeepEqual(gotIn, c.wantIn) {
				t.Fatalf("In: got (%v), want (%v)", gotIn, c.wantIn)
			}

			if errStr := c.wantResponse.Equal(w); errStr != "" {
				t.Fatal(errStr)
			}
		})
	}
}

func TestHTTP_PutProfile(t *testing.T) {
	// in contains all the input parameters (except ctx) of PutProfile.
	type in struct {
		id      string
		profile Profile
	}

	// out contains all the output parameters of PutProfile.
	type out struct {
		err error
	}

	cases := []struct {
		name         string
		request      request
		wantIn       in
		out          out
		wantResponse response
	}{
		{
			name: "ok",
			request: request{
				method: "PUT",
				path:   "/profiles/1234",
				body:   `{"profile": {"id": "1234", "name": "kok", "addresses": [{"id": "0", "location": "here"}]}}`,
			},
			wantIn: in{
				id: "1234",
				profile: Profile{
					ID:   "1234",
					Name: "kok",
					Addresses: []Address{
						{
							ID:       "0",
							Location: "here",
						},
					},
				},
			},
			out: out{
				err: nil,
			},
			wantResponse: response{
				statusCode: http.StatusOK,
				body:       []byte(`{}` + "\n"),
			},
		},
		{
			name: "err",
			request: request{
				method: "PUT",
				path:   "/profiles/1234",
				body:   `{"profile": {"id": "5678", "name": "kok", "addresses": [{"id": "0", "location": "here"}]}}`,
			},
			wantIn: in{
				id: "1234",
				profile: Profile{
					ID:   "5678",
					Name: "kok",
					Addresses: []Address{
						{
							ID:       "0",
							Location: "here",
						},
					},
				},
			},
			out: out{
				err: ErrInconsistentIDs,
			},
			wantResponse: response{
				statusCode: http.StatusBadRequest,
				body:       []byte(`{"error":"inconsistent IDs"}` + "\n"),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var gotIn in
			w := c.request.ServedBy(NewHTTPRouter(
				&ServiceMock{
					PutProfileFunc: func(ctx context.Context, id string, profile Profile) (err error) {
						gotIn = in{
							id:      id,
							profile: profile,
						}
						return c.out.err
					},
				},
				NewCodecs(),
			))

			if !reflect.DeepEqual(gotIn, c.wantIn) {
				t.Fatalf("In: got (%v), want (%v)", gotIn, c.wantIn)
			}

			if errStr := c.wantResponse.Equal(w); errStr != "" {
				t.Fatal(errStr)
			}
		})
	}
}

func TestHTTP_PatchProfile(t *testing.T) {
	// in contains all the input parameters (except ctx) of PatchProfile.
	type in struct {
		id      string
		profile Profile
	}

	// out contains all the output parameters of PatchProfile.
	type out struct {
		err error
	}

	cases := []struct {
		name         string
		request      request
		wantIn       in
		out          out
		wantResponse response
	}{
		{
			name: "ok",
			request: request{
				method: "PATCH",
				path:   "/profiles/1234",
				body:   `{"profile": {"id": "1234", "name": "kok", "addresses": [{"id": "?", "location": "where"}]}}`,
			},
			wantIn: in{
				id: "1234",
				profile: Profile{
					ID:   "1234",
					Name: "kok",
					Addresses: []Address{
						{
							ID:       "?",
							Location: "where",
						},
					},
				},
			},
			out: out{
				err: nil,
			},
			wantResponse: response{
				statusCode: http.StatusOK,
				body:       []byte(`{}` + "\n"),
			},
		},
		{
			name: "err",
			request: request{
				method: "PATCH",
				path:   "/profiles/1234",
				body:   `{"profile": {"id": "5678", "name": "wow"}}`,
			},
			wantIn: in{
				id: "1234",
				profile: Profile{
					ID:   "5678",
					Name: "wow",
				},
			},
			out: out{
				err: ErrInconsistentIDs,
			},
			wantResponse: response{
				statusCode: http.StatusBadRequest,
				body:       []byte(`{"error":"inconsistent IDs"}` + "\n"),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var gotIn in
			w := c.request.ServedBy(NewHTTPRouter(
				&ServiceMock{
					PatchProfileFunc: func(ctx context.Context, id string, profile Profile) (err error) {
						gotIn = in{
							id:      id,
							profile: profile,
						}
						return c.out.err
					},
				},
				NewCodecs(),
			))

			if !reflect.DeepEqual(gotIn, c.wantIn) {
				t.Fatalf("In: got (%v), want (%v)", gotIn, c.wantIn)
			}

			if errStr := c.wantResponse.Equal(w); errStr != "" {
				t.Fatal(errStr)
			}
		})
	}
}

func TestHTTP_DeleteProfile(t *testing.T) {
	// in contains all the input parameters (except ctx) of DeleteProfile.
	type in struct {
		id string
	}

	// out contains all the output parameters of DeleteProfile.
	type out struct {
		err error
	}

	cases := []struct {
		name         string
		request      request
		wantIn       in
		out          out
		wantResponse response
	}{
		{
			name: "ok",
			request: request{
				method: "DELETE",
				path:   "/profiles/1234",
			},
			wantIn: in{
				id: "1234",
			},
			out: out{
				err: nil,
			},
			wantResponse: response{
				statusCode: http.StatusOK,
				body:       []byte(`{}` + "\n"),
			},
		},
		{
			name: "err",
			request: request{
				method: "DELETE",
				path:   "/profiles/5678",
			},
			wantIn: in{
				id: "5678",
			},
			out: out{
				err: ErrNotFound,
			},
			wantResponse: response{
				statusCode: http.StatusNotFound,
				body:       []byte(`{"error":"not found"}` + "\n"),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var gotIn in
			w := c.request.ServedBy(NewHTTPRouter(
				&ServiceMock{
					DeleteProfileFunc: func(ctx context.Context, id string) (err error) {
						gotIn = in{
							id: id,
						}
						return c.out.err
					},
				},
				NewCodecs(),
			))

			if !reflect.DeepEqual(gotIn, c.wantIn) {
				t.Fatalf("In: got (%v), want (%v)", gotIn, c.wantIn)
			}

			if errStr := c.wantResponse.Equal(w); errStr != "" {
				t.Fatal(errStr)
			}
		})
	}
}

func TestHTTP_GetAddresses(t *testing.T) {
	// in contains all the input parameters (except ctx) of GetAddresses.
	type in struct {
		id string
	}

	// out contains all the output parameters of GetAddresses.
	type out struct {
		addresses []Address
		err       error
	}

	cases := []struct {
		name         string
		request      request
		wantIn       in
		out          out
		wantResponse response
	}{
		{
			name: "ok",
			request: request{
				method: "GET",
				path:   "/profiles/1234/addresses",
			},
			wantIn: in{
				id: "1234",
			},
			out: out{
				addresses: []Address{
					{
						ID:       "0",
						Location: "here",
					},
				},
				err: nil,
			},
			wantResponse: response{
				statusCode: http.StatusOK,
				body:       []byte(`{"addresses":[{"id":"0","location":"here"}]}` + "\n"),
			},
		},
		{
			name: "empty",
			request: request{
				method: "GET",
				path:   "/profiles/5678/addresses",
			},
			wantIn: in{
				id: "5678",
			},
			out: out{
				addresses: []Address{},
				err:       nil,
			},
			wantResponse: response{
				statusCode: http.StatusOK,
				body:       []byte(`{"addresses":[]}` + "\n"),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var gotIn in
			w := c.request.ServedBy(NewHTTPRouter(
				&ServiceMock{
					GetAddressesFunc: func(ctx context.Context, id string) (addresses []Address, err error) {
						gotIn = in{
							id: id,
						}
						return c.out.addresses, c.out.err
					},
				},
				NewCodecs(),
			))

			if !reflect.DeepEqual(gotIn, c.wantIn) {
				t.Fatalf("In: got (%v), want (%v)", gotIn, c.wantIn)
			}

			if errStr := c.wantResponse.Equal(w); errStr != "" {
				t.Fatal(errStr)
			}
		})
	}
}

func TestHTTP_GetAddress(t *testing.T) {
	// in contains all the input parameters (except ctx) of GetAddress.
	type in struct {
		id        string
		addressID string
	}

	// out contains all the output parameters of GetAddress.
	type out struct {
		address Address
		err     error
	}

	cases := []struct {
		name         string
		request      request
		wantIn       in
		out          out
		wantResponse response
	}{
		{
			name: "ok",
			request: request{
				method: "GET",
				path:   "/profiles/1234/addresses/0",
			},
			wantIn: in{
				id:        "1234",
				addressID: "0",
			},
			out: out{
				address: Address{
					ID:       "0",
					Location: "here",
				},
				err: nil,
			},
			wantResponse: response{
				statusCode: http.StatusOK,
				body:       []byte(`{"address":{"id":"0","location":"here"}}` + "\n"),
			},
		},
		{
			name: "err",
			request: request{
				method: "GET",
				path:   "/profiles/1234/addresses/9",
			},
			wantIn: in{
				id:        "1234",
				addressID: "9",
			},
			out: out{
				address: Address{},
				err:     ErrNotFound,
			},
			wantResponse: response{
				statusCode: http.StatusNotFound,
				body:       []byte(`{"error":"not found"}` + "\n"),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var gotIn in
			w := c.request.ServedBy(NewHTTPRouter(
				&ServiceMock{
					GetAddressFunc: func(ctx context.Context, id string, addressID string) (address Address, err error) {
						gotIn = in{
							id:        id,
							addressID: addressID,
						}
						return c.out.address, c.out.err
					},
				},
				NewCodecs(),
			))

			if !reflect.DeepEqual(gotIn, c.wantIn) {
				t.Fatalf("In: got (%v), want (%v)", gotIn, c.wantIn)
			}

			if errStr := c.wantResponse.Equal(w); errStr != "" {
				t.Fatal(errStr)
			}
		})
	}
}

func TestHTTP_PostAddress(t *testing.T) {
	// in contains all the input parameters (except ctx) of PostAddress.
	type in struct {
		id      string
		address Address
	}

	// out contains all the output parameters of PostAddress.
	type out struct {
		err error
	}

	cases := []struct {
		name         string
		request      request
		wantIn       in
		out          out
		wantResponse response
	}{
		{
			name: "ok",
			request: request{
				method: "POST",
				path:   "/profiles/1234/addresses",
				body:   `{"address": {"id": "0", "location": "here"}}`,
			},
			wantIn: in{
				id: "1234",
				address: Address{
					ID:       "0",
					Location: "here",
				},
			},
			out: out{
				err: nil,
			},
			wantResponse: response{
				statusCode: http.StatusOK,
				body:       []byte(`{}` + "\n"),
			},
		},
		{
			name: "err",
			request: request{
				method: "POST",
				path:   "/profiles/1234/addresses",
				body:   `{"address": {"id": "0", "location": "here"}}`,
			},
			wantIn: in{
				id: "1234",
				address: Address{
					ID:       "0",
					Location: "here",
				},
			},
			out: out{
				err: ErrAlreadyExists,
			},
			wantResponse: response{
				statusCode: http.StatusBadRequest,
				body:       []byte(`{"error":"already exists"}` + "\n"),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var gotIn in
			w := c.request.ServedBy(NewHTTPRouter(
				&ServiceMock{
					PostAddressFunc: func(ctx context.Context, id string, address Address) (err error) {
						gotIn = in{
							id:      id,
							address: address,
						}
						return c.out.err
					},
				},
				NewCodecs(),
			))

			if !reflect.DeepEqual(gotIn, c.wantIn) {
				t.Fatalf("In: got (%v), want (%v)", gotIn, c.wantIn)
			}

			if errStr := c.wantResponse.Equal(w); errStr != "" {
				t.Fatal(errStr)
			}
		})
	}
}

func TestHTTP_DeleteAddress(t *testing.T) {
	// in contains all the input parameters (except ctx) of DeleteAddress.
	type in struct {
		id        string
		addressID string
	}

	// out contains all the output parameters of DeleteAddress.
	type out struct {
		err error
	}

	cases := []struct {
		name         string
		request      request
		wantIn       in
		out          out
		wantResponse response
	}{
		{
			name: "ok",
			request: request{
				method: "DELETE",
				path:   "/profiles/1234/addresses/0",
			},
			wantIn: in{
				id:        "1234",
				addressID: "0",
			},
			out: out{
				err: nil,
			},
			wantResponse: response{
				statusCode: http.StatusOK,
				body:       []byte(`{}` + "\n"),
			},
		},
		{
			name: "err",
			request: request{
				method: "DELETE",
				path:   "/profiles/1234/addresses/9",
			},
			wantIn: in{
				id:        "1234",
				addressID: "9",
			},
			out: out{
				err: ErrNotFound,
			},
			wantResponse: response{
				statusCode: http.StatusNotFound,
				body:       []byte(`{"error":"not found"}` + "\n"),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var gotIn in
			w := c.request.ServedBy(NewHTTPRouter(
				&ServiceMock{
					DeleteAddressFunc: func(ctx context.Context, id string, addressID string) (err error) {
						gotIn = in{
							id:        id,
							addressID: addressID,
						}
						return c.out.err
					},
				},
				NewCodecs(),
			))

			if !reflect.DeepEqual(gotIn, c.wantIn) {
				t.Fatalf("In: got (%v), want (%v)", gotIn, c.wantIn)
			}

			if errStr := c.wantResponse.Equal(w); errStr != "" {
				t.Fatal(errStr)
			}
		})
	}
}
