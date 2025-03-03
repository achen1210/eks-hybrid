// Code generated by smithy-go-codegen DO NOT EDIT.

package ec2

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Describes the specified Amazon Web Services Verified Access endpoints.
func (c *Client) DescribeVerifiedAccessEndpoints(ctx context.Context, params *DescribeVerifiedAccessEndpointsInput, optFns ...func(*Options)) (*DescribeVerifiedAccessEndpointsOutput, error) {
	if params == nil {
		params = &DescribeVerifiedAccessEndpointsInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "DescribeVerifiedAccessEndpoints", params, optFns, c.addOperationDescribeVerifiedAccessEndpointsMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*DescribeVerifiedAccessEndpointsOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type DescribeVerifiedAccessEndpointsInput struct {

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have the
	// required permissions, the error response is DryRunOperation . Otherwise, it is
	// UnauthorizedOperation .
	DryRun *bool

	// One or more filters. Filter names and values are case-sensitive.
	Filters []types.Filter

	// The maximum number of results to return with a single call. To retrieve the
	// remaining results, make another call with the returned nextToken value.
	MaxResults *int32

	// The token for the next page of results.
	NextToken *string

	// The ID of the Verified Access endpoint.
	VerifiedAccessEndpointIds []string

	// The ID of the Verified Access group.
	VerifiedAccessGroupId *string

	// The ID of the Verified Access instance.
	VerifiedAccessInstanceId *string

	noSmithyDocumentSerde
}

type DescribeVerifiedAccessEndpointsOutput struct {

	// The token to use to retrieve the next page of results. This value is null when
	// there are no more results to return.
	NextToken *string

	// Details about the Verified Access endpoints.
	VerifiedAccessEndpoints []types.VerifiedAccessEndpoint

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationDescribeVerifiedAccessEndpointsMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsEc2query_serializeOpDescribeVerifiedAccessEndpoints{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsEc2query_deserializeOpDescribeVerifiedAccessEndpoints{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "DescribeVerifiedAccessEndpoints"); err != nil {
		return fmt.Errorf("add protocol finalizers: %v", err)
	}

	if err = addlegacyEndpointContextSetter(stack, options); err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = addClientRequestID(stack); err != nil {
		return err
	}
	if err = addComputeContentLength(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = addComputePayloadSHA256(stack); err != nil {
		return err
	}
	if err = addRetry(stack, options); err != nil {
		return err
	}
	if err = addRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = addRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addSpanRetryLoop(stack, options); err != nil {
		return err
	}
	if err = addClientUserAgent(stack, options); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addSetLegacyContextSigningOptionsMiddleware(stack); err != nil {
		return err
	}
	if err = addTimeOffsetBuild(stack, c); err != nil {
		return err
	}
	if err = addUserAgentRetryMode(stack, options); err != nil {
		return err
	}
	if err = addCredentialSource(stack, options); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opDescribeVerifiedAccessEndpoints(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addRecursionDetection(stack); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	if err = addDisableHTTPSMiddleware(stack, options); err != nil {
		return err
	}
	if err = addSpanInitializeStart(stack); err != nil {
		return err
	}
	if err = addSpanInitializeEnd(stack); err != nil {
		return err
	}
	if err = addSpanBuildRequestStart(stack); err != nil {
		return err
	}
	if err = addSpanBuildRequestEnd(stack); err != nil {
		return err
	}
	return nil
}

// DescribeVerifiedAccessEndpointsPaginatorOptions is the paginator options for
// DescribeVerifiedAccessEndpoints
type DescribeVerifiedAccessEndpointsPaginatorOptions struct {
	// The maximum number of results to return with a single call. To retrieve the
	// remaining results, make another call with the returned nextToken value.
	Limit int32

	// Set to true if pagination should stop if the service returns a pagination token
	// that matches the most recent token provided to the service.
	StopOnDuplicateToken bool
}

// DescribeVerifiedAccessEndpointsPaginator is a paginator for
// DescribeVerifiedAccessEndpoints
type DescribeVerifiedAccessEndpointsPaginator struct {
	options   DescribeVerifiedAccessEndpointsPaginatorOptions
	client    DescribeVerifiedAccessEndpointsAPIClient
	params    *DescribeVerifiedAccessEndpointsInput
	nextToken *string
	firstPage bool
}

// NewDescribeVerifiedAccessEndpointsPaginator returns a new
// DescribeVerifiedAccessEndpointsPaginator
func NewDescribeVerifiedAccessEndpointsPaginator(client DescribeVerifiedAccessEndpointsAPIClient, params *DescribeVerifiedAccessEndpointsInput, optFns ...func(*DescribeVerifiedAccessEndpointsPaginatorOptions)) *DescribeVerifiedAccessEndpointsPaginator {
	if params == nil {
		params = &DescribeVerifiedAccessEndpointsInput{}
	}

	options := DescribeVerifiedAccessEndpointsPaginatorOptions{}
	if params.MaxResults != nil {
		options.Limit = *params.MaxResults
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return &DescribeVerifiedAccessEndpointsPaginator{
		options:   options,
		client:    client,
		params:    params,
		firstPage: true,
		nextToken: params.NextToken,
	}
}

// HasMorePages returns a boolean indicating whether more pages are available
func (p *DescribeVerifiedAccessEndpointsPaginator) HasMorePages() bool {
	return p.firstPage || (p.nextToken != nil && len(*p.nextToken) != 0)
}

// NextPage retrieves the next DescribeVerifiedAccessEndpoints page.
func (p *DescribeVerifiedAccessEndpointsPaginator) NextPage(ctx context.Context, optFns ...func(*Options)) (*DescribeVerifiedAccessEndpointsOutput, error) {
	if !p.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	params := *p.params
	params.NextToken = p.nextToken

	var limit *int32
	if p.options.Limit > 0 {
		limit = &p.options.Limit
	}
	params.MaxResults = limit

	optFns = append([]func(*Options){
		addIsPaginatorUserAgent,
	}, optFns...)
	result, err := p.client.DescribeVerifiedAccessEndpoints(ctx, &params, optFns...)
	if err != nil {
		return nil, err
	}
	p.firstPage = false

	prevToken := p.nextToken
	p.nextToken = result.NextToken

	if p.options.StopOnDuplicateToken &&
		prevToken != nil &&
		p.nextToken != nil &&
		*prevToken == *p.nextToken {
		p.nextToken = nil
	}

	return result, nil
}

// DescribeVerifiedAccessEndpointsAPIClient is a client that implements the
// DescribeVerifiedAccessEndpoints operation.
type DescribeVerifiedAccessEndpointsAPIClient interface {
	DescribeVerifiedAccessEndpoints(context.Context, *DescribeVerifiedAccessEndpointsInput, ...func(*Options)) (*DescribeVerifiedAccessEndpointsOutput, error)
}

var _ DescribeVerifiedAccessEndpointsAPIClient = (*Client)(nil)

func newServiceMetadataMiddleware_opDescribeVerifiedAccessEndpoints(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "DescribeVerifiedAccessEndpoints",
	}
}
