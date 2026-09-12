package responseconverter

import (
	"github.com/rossigee/provider-http/apis/request/v1beta1"
	httpClient "github.com/rossigee/provider-http/internal/clients/http"
)

// Convert HttpResponse to Response
func HttpResponseToV1alpha1Response(httpResponse httpClient.HttpResponse) v1beta1.Response {
	return v1beta1.Response{
		StatusCode: httpResponse.StatusCode,
		Body:       httpResponse.Body,
		Headers:    httpResponse.Headers,
	}
}
