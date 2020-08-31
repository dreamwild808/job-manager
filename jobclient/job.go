package jobclient

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"

	apiv1 "github.com/jeffrom/job-manager/pkg/api/v1"
	"github.com/jeffrom/job-manager/pkg/job"
)

func (c *Client) EnqueueJob(ctx context.Context, name string, args ...interface{}) (string, error) {
	argList, err := structpb.NewList(args)
	if err != nil {
		return "", err
	}
	params := &apiv1.EnqueueParams{
		Jobs: []*apiv1.EnqueueParamArgs{{Job: name, Args: argList.Values}},
	}

	uri := fmt.Sprintf("/api/v1/jobs/%s/enqueue", name)
	req, err := c.newRequestProto("POST", uri, params)
	if err != nil {
		return "", err
	}

	jobs := &job.Jobs{}
	err = c.doRequest(ctx, req, jobs)
	if err != nil {
		return "", err
	}

	return jobs.Jobs[0].Id, nil
}

func (c *Client) DequeueJobs(ctx context.Context, num int, jobName string, selectors ...string) (*job.Jobs, error) {
	params := &apiv1.DequeueParams{
		Selectors: selectors,
	}
	if num > 0 {
		params.Num = proto.Int32(int32(num))
	}
	if jobName != "" {
		params.Job = proto.String(jobName)
	}

	uri := "/api/v1/jobs/dequeue"
	if jobName != "" {
		uri = fmt.Sprintf("/api/v1/jobs/%s/dequeue", jobName)
	}
	req, err := c.newRequestProto("POST", uri, params)
	if err != nil {
		return nil, err
	}

	jobs := &job.Jobs{}
	err = c.doRequest(ctx, req, jobs)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

type AckJobOpts struct {
	Data map[string]interface{}
}

func (c *Client) AckJob(ctx context.Context, id string, status job.Status) error {
	return c.AckJobOpts(ctx, id, status, AckJobOpts{})
}

func (c *Client) AckJobOpts(ctx context.Context, id string, status job.Status, opts AckJobOpts) error {
	args := &apiv1.AckParamArgs{
		Id:     id,
		Status: status,
	}
	if len(opts.Data) > 0 {
		data, err := structpb.NewStruct(opts.Data)
		if err != nil {
			return err
		}
		args.Data = data.Fields
	}

	uri := "/api/v1/jobs/ack"
	params := &apiv1.AckParams{Acks: []*apiv1.AckParamArgs{args}}
	req, err := c.newRequestProto("POST", uri, params)
	if err != nil {
		return err
	}
	return c.doRequest(ctx, req, nil)
}

func (c *Client) GetJob(ctx context.Context, id string) (*job.Job, error) {
	uri := fmt.Sprintf("/api/v1/job/%s", id)
	req, err := c.newRequestProto("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	jobData := &job.Job{}
	if err := c.doRequest(ctx, req, jobData); err != nil {
		return nil, err
	}
	return jobData, nil
}
