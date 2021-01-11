package cloudformation

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	client "github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/yawn/alm/file"
)

var clock = time.Now

type cloudformation struct {
	client      *client.Client
	eventStamps map[string]*time.Time
	region      string
	started     time.Time
}

func New(cfg aws.Config) *cloudformation {

	return &cloudformation{
		client:      client.NewFromConfig(cfg),
		eventStamps: make(map[string]*time.Time),
		started:     clock(),
		region:      cfg.Region,
	}

}

func (c *cloudformation) Discover(ctx context.Context, maxAge time.Duration) ([]file.Target, error) {

	var targets []file.Target

	req := client.NewListStacksPaginator(c.client, &client.ListStacksInput{})

	for req.HasMorePages() {

		res, err := req.NextPage(ctx)

		if err != nil {
			return nil, errors.Wrapf(err, "failed to list stacks")
		}

		for _, stack := range res.StackSummaries {

			if stack.StackStatus == types.StackStatusDeleteComplete {
				continue
			}

			var ref = stack.LastUpdatedTime

			if ref == nil { // for new stacks
				ref = stack.CreationTime
			}

			if ref.Add(maxAge).After(c.started) {
				targets = append(targets, file.Target{
					ID:      stack.StackId,
					Name:    stack.StackName,
					Region:  c.region,
					Ref:     ref,
					Service: file.CloudFormation,
				})

			}

		}

	}

	return targets, nil

}

func (c *cloudformation) Log(ctx context.Context, writer *file.Writer, target file.Target) error {

	id := target.ID

	req := client.NewDescribeStackEventsPaginator(c.client, &client.DescribeStackEventsInput{
		StackName: id,
	})

	var events []types.StackEvent

	wc, err := writer.Add(target.Parts()...)

	if err != nil {
		return errors.Wrapf(err, "failed to open log writer")
	}

	defer wc.Close()

	for req.HasMorePages() {

		res, err := req.NextPage(ctx)

		if err != nil {
			return errors.Wrapf(err, "failed to list stack events for stack %q", *id)
		}

		for _, event := range res.StackEvents {
			events = append([]types.StackEvent{event}, events...)
		}

	}

	for _, event := range events {

		var (
			buf      bytes.Buffer
			last, ok = c.eventStamps[*id]
		)

		if !ok || last.Before(*event.Timestamp) {

			var (
				colorF func(string, ...interface{}) string
				status = string(event.ResourceStatus)
			)

			if strings.Contains(status, "COMPLETE") {
				colorF = color.GreenString
			} else if strings.Contains(status, "PROGRESS") {
				colorF = color.YellowString
			} else if strings.Contains(status, "FAILED") {
				colorF = color.RedString
			}

			buf.WriteString(fmt.Sprintf("%s %q %s (%s) %s",
				event.Timestamp.Format("2006-01-02T15:04:05"),
				*target.Name,
				*event.LogicalResourceId,
				*event.ResourceType,
				colorF(status),
			))

			if event.ResourceStatusReason != nil {
				buf.WriteString(fmt.Sprintf(" - %s", *event.ResourceStatusReason))
			}

			c.eventStamps[*id] = event.Timestamp

			fmt.Fprintln(wc, buf.String())

		}

	}

	return nil

}
