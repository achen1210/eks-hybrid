// Code generated by smithy-go-codegen DO NOT EDIT.

package ec2

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"time"
)

// Describes Capacity Block offerings available for purchase in the Amazon Web
// Services Region that you're currently using. With Capacity Blocks, you purchase
// a specific instance type for a period of time.
//
// To search for an available Capacity Block offering, you specify a reservation
// duration and instance count. You must select one of the following options.
//
//   - For reservation durations 1-day increments up 14 days and 7-day increments
//     up to 182 days total
//
//   - For instance count 1, 2, 4, 8, 16, 32, or 64 instances
func (c *Client) DescribeCapacityBlockOfferings(ctx context.Context, params *DescribeCapacityBlockOfferingsInput, optFns ...func(*Options)) (*DescribeCapacityBlockOfferingsOutput, error) {
	if params == nil {
		params = &DescribeCapacityBlockOfferingsInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "DescribeCapacityBlockOfferings", params, optFns, c.addOperationDescribeCapacityBlockOfferingsMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*DescribeCapacityBlockOfferingsOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type DescribeCapacityBlockOfferingsInput struct {

	// The number of hours for which to reserve Capacity Block.
	//
	// This member is required.
	CapacityDurationHours *int32

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have the
	// required permissions, the error response is DryRunOperation . Otherwise, it is
	// UnauthorizedOperation .
	DryRun *bool

	// The latest end date for the Capacity Block offering.
	EndDateRange *time.Time

	// The number of instances for which to reserve capacity.
	InstanceCount *int32

	// The type of instance for which the Capacity Block offering reserves capacity.
	InstanceType *string

	// The maximum number of items to return for this request. To get the next page of
	// items, make another request with the token returned in the output. For more
	// information, see [Pagination].
	//
	// [Pagination]: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/Query-Requests.html#api-pagination
	MaxResults *int32

	// The token to use to retrieve the next page of results.
	NextToken *string

	// The earliest start date for the Capacity Block offering.
	StartDateRange *time.Time

	noSmithyDocumentSerde
}

type DescribeCapacityBlockOfferingsOutput struct {

	// The recommended Capacity Block offering for the dates specified.
	CapacityBlockOfferings []types.CapacityBlockOffering

	// The token to use to retrieve the next page of results. This value is null when
	// there are no more results to return.
	NextToken *string

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationDescribeCapacityBlockOfferingsMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsEc2query_serializeOpDescribeCapacityBlockOfferings{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsEc2query_deserializeOpDescribeCapacityBlockOfferings{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "DescribeCapacityBlockOfferings"); err != nil {
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
	if err = addOpDescribeCapacityBlockOfferingsValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opDescribeCapacityBlockOfferings(options.Region), middleware.Before); err != nil {
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

// DescribeCapacityBlockOfferingsPaginatorOptions is the paginator options for
// DescribeCapacityBlockOfferings
type DescribeCapacityBlockOfferingsPaginatorOptions struct {
	// The maximum number of items to return for this request. To get the next page of
	// items, make another request with the token returned in the output. For more
	// information, see [Pagination].
	//
	// [Pagination]: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/Query-Requests.html#api-pagination
	Limit int32

	// Set to true if pagination should stop if the service returns a pagination token
	// that matches the most recent token provided to the service.
	StopOnDuplicateToken bool
}

// DescribeCapacityBlockOfferingsPaginator is a paginator for
// DescribeCapacityBlockOfferings
type DescribeCapacityBlockOfferingsPaginator struct {
	options   DescribeCapacityBlockOfferingsPaginatorOptions
	client    DescribeCapacityBlockOfferingsAPIClient
	params    *DescribeCapacityBlockOfferingsInput
	nextToken *string
	firstPage bool
}

// NewDescribeCapacityBlockOfferingsPaginator returns a new
// DescribeCapacityBlockOfferingsPaginator
func NewDescribeCapacityBlockOfferingsPaginator(client DescribeCapacityBlockOfferingsAPIClient, params *DescribeCapacityBlockOfferingsInput, optFns ...func(*DescribeCapacityBlockOfferingsPaginatorOptions)) *DescribeCapacityBlockOfferingsPaginator {
	if params == nil {
		params = &DescribeCapacityBlockOfferingsInput{}
	}

	options := DescribeCapacityBlockOfferingsPaginatorOptions{}
	if params.MaxResults != nil {
		options.Limit = *params.MaxResults
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return &DescribeCapacityBlockOfferingsPaginator{
		options:   options,
		client:    client,
		params:    params,
		firstPage: true,
		nextToken: params.NextToken,
	}
}

// HasMorePages returns a boolean indicating whether more pages are available
func (p *DescribeCapacityBlockOfferingsPaginator) HasMorePages() bool {
	return p.firstPage || (p.nextToken != nil && len(*p.nextToken) != 0)
}

// NextPage retrieves the next DescribeCapacityBlockOfferings page.
func (p *DescribeCapacityBlockOfferingsPaginator) NextPage(ctx context.Context, optFns ...func(*Options)) (*DescribeCapacityBlockOfferingsOutput, error) {
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
	result, err := p.client.DescribeCapacityBlockOfferings(ctx, &params, optFns...)
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

// DescribeCapacityBlockOfferingsAPIClient is a client that implements the
// DescribeCapacityBlockOfferings operation.
type DescribeCapacityBlockOfferingsAPIClient interface {
	DescribeCapacityBlockOfferings(context.Context, *DescribeCapacityBlockOfferingsInput, ...func(*Options)) (*DescribeCapacityBlockOfferingsOutput, error)
}

var _ DescribeCapacityBlockOfferingsAPIClient = (*Client)(nil)

func newServiceMetadataMiddleware_opDescribeCapacityBlockOfferings(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "DescribeCapacityBlockOfferings",
	}
}
