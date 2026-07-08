package responseconverter

import (
	"github.com/google/go-cmp/cmp"
	"github.com/rossigee/provider-http/apis/request/v1alpha2"
	"github.com/rossigee/provider-http/internal/clients/http"
	"testing"
)

var testHeaders = map[string][]string{
	"fruits":                {"apple", "banana", "orange"},
	"colors":                {"red", "green", "blue"},
	"countries":             {"USA", "UK", "India", "Germany"},
	"programming_languages": {"Go", "Python", "JavaScript"},
}

func Test_HttpResponseToV1alpha1Response(t *testing.T) {
	type args struct {
		httpResponse httpClient.HttpResponse
	}
	type want struct {
		result v1alpha2.Response
	}
	cases := map[string]struct {
		args args
		want want
	}{
		"Success": {
			args: args{
				httpResponse: httpClient.HttpResponse{
					Body:       `{"email":"john.doe@example.com","name":"john_doe"}`,
					Headers:    testHeaders,
					StatusCode: 200,
				},
			},
			want: want{
				result: v1alpha2.Response{
					Body:       `{"email":"john.doe@example.com","name":"john_doe"}`,
					Headers:    testHeaders,
					StatusCode: 200,
				},
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := HttpResponseToV1alpha1Response(tc.args.httpResponse)
			if diff := cmp.Diff(tc.want.result, got); diff != "" {
				t.Fatalf("HttpResponseToV1alpha1Response(...): -want result, +got result: %s", diff)
			}
		})
	}

}
