// Code generated by smithy-go-codegen DO NOT EDIT.

package ec2

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Modifies the attributes of the specified VPC endpoint service configuration.
//
// If you set or modify the private DNS name, you must prove that you own the
// private DNS domain name.
func (c *Client) ModifyVpcEndpointServiceConfiguration(ctx context.Context, params *ModifyVpcEndpointServiceConfigurationInput, optFns ...func(*Options)) (*ModifyVpcEndpointServiceConfigurationOutput, error) {
	if params == nil {
		params = &ModifyVpcEndpointServiceConfigurationInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "ModifyVpcEndpointServiceConfiguration", params, optFns, c.addOperationModifyVpcEndpointServiceConfigurationMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*ModifyVpcEndpointServiceConfigurationOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type ModifyVpcEndpointServiceConfigurationInput struct {

	// The ID of the service.
	//
	// This member is required.
	ServiceId *string

	// Indicates whether requests to create an endpoint to the service must be
	// accepted.
	AcceptanceRequired *bool

	// The Amazon Resource Names (ARNs) of Gateway Load Balancers to add to the
	// service configuration.
	AddGatewayLoadBalancerArns []string

	// The Amazon Resource Names (ARNs) of Network Load Balancers to add to the
	// service configuration.
	AddNetworkLoadBalancerArns []string

	// The IP address types to add to the service configuration.
	AddSupportedIpAddressTypes []string

	// The supported Regions to add to the service configuration.
	AddSupportedRegions []string

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have the
	// required permissions, the error response is DryRunOperation . Otherwise, it is
	// UnauthorizedOperation .
	DryRun *bool

	// (Interface endpoint configuration) The private DNS name to assign to the
	// endpoint service.
	PrivateDnsName *string

	// The Amazon Resource Names (ARNs) of Gateway Load Balancers to remove from the
	// service configuration.
	RemoveGatewayLoadBalancerArns []string

	// The Amazon Resource Names (ARNs) of Network Load Balancers to remove from the
	// service configuration.
	RemoveNetworkLoadBalancerArns []string

	// (Interface endpoint configuration) Removes the private DNS name of the endpoint
	// service.
	RemovePrivateDnsName *bool

	// The IP address types to remove from the service configuration.
	RemoveSupportedIpAddressTypes []string

	// The supported Regions to remove from the service configuration.
	RemoveSupportedRegions []string

	noSmithyDocumentSerde
}

type ModifyVpcEndpointServiceConfigurationOutput struct {

	// Returns true if the request succeeds; otherwise, it returns an error.
	Return *bool

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationModifyVpcEndpointServiceConfigurationMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsEc2query_serializeOpModifyVpcEndpointServiceConfiguration{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsEc2query_deserializeOpModifyVpcEndpointServiceConfiguration{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "ModifyVpcEndpointServiceConfiguration"); err != nil {
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
	if err = addOpModifyVpcEndpointServiceConfigurationValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opModifyVpcEndpointServiceConfiguration(options.Region), middleware.Before); err != nil {
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

func newServiceMetadataMiddleware_opModifyVpcEndpointServiceConfiguration(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "ModifyVpcEndpointServiceConfiguration",
	}
}
