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

// Shuts down the specified instances. This operation is idempotent; if you
// terminate an instance more than once, each call succeeds.
//
// If you specify multiple instances and the request fails (for example, because
// of a single incorrect instance ID), none of the instances are terminated.
//
// If you terminate multiple instances across multiple Availability Zones, and one
// or more of the specified instances are enabled for termination protection, the
// request fails with the following results:
//
//   - The specified instances that are in the same Availability Zone as the
//     protected instance are not terminated.
//
//   - The specified instances that are in different Availability Zones, where no
//     other specified instances are protected, are successfully terminated.
//
// For example, say you have the following instances:
//
//   - Instance A: us-east-1a ; Not protected
//
//   - Instance B: us-east-1a ; Not protected
//
//   - Instance C: us-east-1b ; Protected
//
//   - Instance D: us-east-1b ; not protected
//
// If you attempt to terminate all of these instances in the same request, the
// request reports failure with the following results:
//
//   - Instance A and Instance B are successfully terminated because none of the
//     specified instances in us-east-1a are enabled for termination protection.
//
//   - Instance C and Instance D fail to terminate because at least one of the
//     specified instances in us-east-1b (Instance C) is enabled for termination
//     protection.
//
// Terminated instances remain visible after termination (for approximately one
// hour).
//
// By default, Amazon EC2 deletes all EBS volumes that were attached when the
// instance launched. Volumes attached after instance launch continue running.
//
// You can stop, start, and terminate EBS-backed instances. You can only terminate
// instance store-backed instances. What happens to an instance differs if you stop
// it or terminate it. For example, when you stop an instance, the root device and
// any other devices attached to the instance persist. When you terminate an
// instance, any attached EBS volumes with the DeleteOnTermination block device
// mapping parameter set to true are automatically deleted. For more information
// about the differences between stopping and terminating instances, see [Instance lifecycle]in the
// Amazon EC2 User Guide.
//
// For more information about troubleshooting, see [Troubleshooting terminating your instance] in the Amazon EC2 User Guide.
//
// [Instance lifecycle]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-lifecycle.html
// [Troubleshooting terminating your instance]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/TroubleshootingInstancesShuttingDown.html
func (c *Client) TerminateInstances(ctx context.Context, params *TerminateInstancesInput, optFns ...func(*Options)) (*TerminateInstancesOutput, error) {
	if params == nil {
		params = &TerminateInstancesInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "TerminateInstances", params, optFns, c.addOperationTerminateInstancesMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*TerminateInstancesOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type TerminateInstancesInput struct {

	// The IDs of the instances.
	//
	// Constraints: Up to 1000 instance IDs. We recommend breaking up this request
	// into smaller batches.
	//
	// This member is required.
	InstanceIds []string

	// Checks whether you have the required permissions for the operation, without
	// actually making the request, and provides an error response. If you have the
	// required permissions, the error response is DryRunOperation . Otherwise, it is
	// UnauthorizedOperation .
	DryRun *bool

	noSmithyDocumentSerde
}

type TerminateInstancesOutput struct {

	// Information about the terminated instances.
	TerminatingInstances []types.InstanceStateChange

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationTerminateInstancesMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsEc2query_serializeOpTerminateInstances{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsEc2query_deserializeOpTerminateInstances{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "TerminateInstances"); err != nil {
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
	if err = addOpTerminateInstancesValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opTerminateInstances(options.Region), middleware.Before); err != nil {
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

func newServiceMetadataMiddleware_opTerminateInstances(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "TerminateInstances",
	}
}
