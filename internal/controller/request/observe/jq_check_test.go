package observe

import (
	"context"
	"testing"

	"github.com/crossplane/crossplane-runtime/v2/pkg/logging"
	"github.com/crossplane/crossplane-runtime/v2/pkg/test"
	"github.com/google/go-cmp/cmp"
	"github.com/rossigee/provider-http/apis/request/v1beta1"
	httpClient "github.com/rossigee/provider-http/internal/clients/http"
)

func Test_CustomCheck(t *testing.T) {
	type args struct {
		ctx     context.Context
		cr      *v1beta1.Request
		details httpClient.HttpDetails
		logic   string
	}

	type want struct {
		result bool
		err    error
	}

	cases := map[string]struct {
		args args
		want want
	}{
		"CustomCheckPasses": {
			args: args{
				ctx: context.Background(),
				cr: &v1beta1.Request{
					Spec: v1beta1.RequestSpec{
						ForProvider: v1beta1.RequestParameters{
							Payload: v1beta1.Payload{
								Body: `{"password": "password"}`,
							},
							ExpectedResponseCheck: v1beta1.ExpectedResponseCheck{
								Type:  v1beta1.ExpectedResponseCheckTypeCustom,
								Logic: `.response.body.password == .payload.body.password`,
							},
						},
					},
				},
				details: httpClient.HttpDetails{
					HttpResponse: httpClient.HttpResponse{
						Body:       `{"password":"password"}`,
						Headers:    nil,
						StatusCode: 0,
					},
				},
				logic: `.response.body.password == .payload.body.password`,
			},
			want: want{
				result: true,
				err:    nil,
			},
		},
		"CustomCheckFails": {
			args: args{
				ctx: context.Background(),
				cr: &v1beta1.Request{
					Spec: v1beta1.RequestSpec{
						ForProvider: v1beta1.RequestParameters{
							Payload: v1beta1.Payload{
								Body: `{"password": "password"}`,
							},
							ExpectedResponseCheck: v1beta1.ExpectedResponseCheck{
								Type:  v1beta1.ExpectedResponseCheckTypeCustom,
								Logic: `.response.body.password == .payload.body.password`,
							},
							IsRemovedCheck: v1beta1.ExpectedResponseCheck{
								Type:  v1beta1.ExpectedResponseCheckTypeCustom,
								Logic: `.response.body.password == .payload.body.password`,
							},
						},
					},
				},
				details: httpClient.HttpDetails{
					HttpResponse: httpClient.HttpResponse{
						Body:       `{"password":"wrong_password"}`,
						Headers:    nil,
						StatusCode: 0,
					},
				},
				logic: `.response.body.password == .payload.body.password`,
			},
			want: want{
				result: false,
				err:    nil,
			},
		},
	}

	for name, tc := range cases {
		tc := tc // Create local copies of loop variables

		t.Run(name, func(t *testing.T) {
			e := &customCheck{
				localKube: nil,
				http:      nil,
				logger:    logging.NewNopLogger(),
			}
			got, gotErr := e.check(tc.args.ctx, tc.args.cr, tc.args.details, tc.args.logic)
			if diff := cmp.Diff(tc.want.err, gotErr, test.EquateErrors()); diff != "" {
				t.Fatalf("Check(...): -want error, +got error: %s", diff)
			}

			if diff := cmp.Diff(tc.want.result, got); diff != "" {
				t.Fatalf("Check(...): -want result, +got result: %s", diff)
			}
		})
	}
}
