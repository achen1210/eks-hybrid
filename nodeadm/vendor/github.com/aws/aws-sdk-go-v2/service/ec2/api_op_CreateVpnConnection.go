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

// Creates a VPN connection between an existing virtual private gateway or transit
// gateway and a customer gateway. The supported connection type is ipsec.1 .
//
// The response includes information that you need to give to your network
// administrator to configure your customer gateway.
//
// We strongly recommend that you use HTTPS when calling this operation because
// the response contains sensitive cryptographic information for configuring your
// customer gateway device.
//
// If you decide to shut down your VPN connection for any reason and later create
// a new VPN connection, you must reconfigure your customer gateway with the new
// information returned from this call.
//
// This is an idempotent operation. If you perform the operation more than once,
// Amazon EC2 doesn't return an error.
//
// For more information, see [Amazon Web Services Site-to-Site VPN] in the Amazon Web Services Site-to-Site VPN User
// Guide.
//
// [Amazon Web Services Site-to-Site VPN]: https://docs.aws.amazon.com/vpn/latest/s2svpn/VPC_VPN.html
func (c *Client) CreateVpnConnection(ctx context.Context, params *CreateVpnConnectionInput, optFns ...func(*Options)) (*CreateVpnConnectionOutput, error) {
	if params == nil {
		params = &CreateVpnConnectionInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "CreateVpnConnection", params, optFns, c.addOperationCreateVpnConnectionMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*CreateVpnConnectionOutput)
	out.ResultMetadata = metadata
	return out, nil
}

// Contains the parameters for CreateVpnConnection.
type CreateVpnConnectionInput struct {

	// The ID of the customer gateway.
	//
	// This member is required.
	CustomerGatewayId *string

	// The type of VPN connection ( ipsec.1 ).
	//
	// This member is required.
	Type *string

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have the
	// required permissions, the error response is DryRunOperation . Otherwise, it is
	// UnauthorizedOperation .
	DryRun *bool

	// The options for the VPN connection.
	Options *types.VpnConnectionOptionsSpecification

	// The tags to apply to the VPN connection.
	TagSpecifications []types.TagSpecification

	// The ID of the transit gateway. If you specify a transit gateway, you cannot
	// specify a virtual private gateway.
	TransitGatewayId *string

	// The ID of the virtual private gateway. If you specify a virtual private
	// gateway, you cannot specify a transit gateway.
	VpnGatewayId *string

	noSmithyDocumentSerde
}

// Contains the output of CreateVpnConnection.
type CreateVpnConnectionOutput struct {

	// Information about the VPN connection.
	VpnConnection *types.VpnConnection

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationCreateVpnConnectionMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsEc2query_serializeOpCreateVpnConnection{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsEc2query_deserializeOpCreateVpnConnection{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "CreateVpnConnection"); err != nil {
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
	if err = addOpCreateVpnConnectionValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opCreateVpnConnection(options.Region), middleware.Before); err != nil {
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

func newServiceMetadataMiddleware_opCreateVpnConnection(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "CreateVpnConnection",
	}
}
