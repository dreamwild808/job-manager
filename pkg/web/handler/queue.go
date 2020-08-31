package handler

import (
	"net/http"

	"github.com/go-chi/chi"
	"google.golang.org/protobuf/types/known/durationpb"

	apiv1 "github.com/jeffrom/job-manager/pkg/api/v1"
	"github.com/jeffrom/job-manager/pkg/job"
	"github.com/jeffrom/job-manager/pkg/web/middleware"
)

func SaveQueue(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	cfg := middleware.ConfigFromContext(ctx)
	reqLog := middleware.RequestLogFromContext(ctx)
	be := middleware.GetBackend(ctx)
	name := chi.URLParam(r, "queueName")
	reqLog.Str("queue", name)

	var params apiv1.SaveQueueParamArgs
	if err := UnmarshalBody(r, &params, false); err != nil {
		return err
	}

	var concurrency int32 = int32(cfg.DefaultConcurrency)
	if conc := params.Concurrency; conc != nil {
		concurrency = *conc
	}
	var maxRetries int32 = int32(cfg.DefaultMaxRetries)
	if mr := params.MaxRetries; mr != nil {
		maxRetries = *mr
	}

	dur := durationpb.New(cfg.DefaultMaxJobTimeout)
	if d := params.Duration; d != nil {
		dur = d
	}

	queue := &job.Queue{
		Name:        name,
		Concurrency: concurrency,
		MaxRetries:  maxRetries,
		Labels:      params.Labels,
		Duration:    dur,
	}
	if err := be.SaveQueue(ctx, queue); err != nil {
		return err
	}
	return MarshalResponse(w, r, queue)
}

func DeleteQueue(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func ListQueues(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	be := middleware.GetBackend(ctx)
	jobs, err := be.ListQueues(ctx, nil)
	if err != nil {
		return err
	}
	return MarshalResponse(w, r, jobs)
}
