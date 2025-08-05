package observe

import (
	"context"
	"fmt"

	"github.com/rossigee/provider-http/apis/request/v1alpha2"
	httpClient "github.com/rossigee/provider-http/internal/clients/http"
	"github.com/rossigee/provider-http/internal/controller/request/requestgen"
	"github.com/rossigee/provider-http/internal/controller/request/responseconverter"
	datapatcher "github.com/rossigee/provider-http/internal/data-patcher"
	"github.com/rossigee/provider-http/internal/jq"
	"github.com/rossigee/provider-http/internal/utils"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// responseCheck is an interface for performing response checks.
type responseCheck interface {
	Check(ctx context.Context, cr *v1alpha2.Request, details httpClient.HttpDetails, responseErr error) (bool, error)
}

// customCheck performs a custom response check using JQ logic.
type customCheck struct {
	localKube client.Client
	logger    logging.Logger
	http      httpClient.Client
}

// Check performs a custom response check using JQ logic.
func (c *customCheck) check(ctx context.Context, cr *v1alpha2.Request, details httpClient.HttpDetails, logic string) (bool, error) {
	// Convert response to a map and apply JQ logic
	response := responseconverter.HttpResponseToV1alpha1Response(details.HttpResponse)

	sensitiveResponse, err := datapatcher.PatchSecretsIntoResponse(ctx, c.localKube, response, c.logger)
	if err != nil {
		return false, err
	}

	sensitiveRequestContext := requestgen.GenerateRequestContext(cr.Spec.ForProvider, sensitiveResponse)

	jqQuery := utils.NormalizeWhitespace(logic)
	sensitiveJQQuery, err := datapatcher.PatchSecretsIntoString(ctx, c.localKube, jqQuery, c.logger)
	if err != nil {
		return false, err
	}

	isExpected, err := jq.ParseBool(sensitiveJQQuery, sensitiveRequestContext)

	c.logger.Debug(fmt.Sprintf("Applying JQ filter %s, result is %v", jqQuery, isExpected))
	if err != nil {
		return false, err
	}

	return isExpected, nil
}
