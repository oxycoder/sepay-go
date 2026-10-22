package einvoice

import (
	"context"
	"time"
)

// settlingPollInterval is used once SePay has no retry scheduled and is settling the result.
const settlingPollInterval = time.Minute

// WaitForCreate polls GetCreateStatus until the create, adjust or replace
// request reaches Success or Failed, or ctx is done.
//
// A request can stay Pending for up to 2 hours while SePay retries transient
// failures, so give ctx a generous deadline. A Failed status is returned as a
// result, not an error; check ProcessResult.Status (and Unconfirmed) yourself.
// On ctx expiry the last result is returned with ctx.Err().
func (c *Client) WaitForCreate(ctx context.Context, trackingCode string) (*ProcessResult, error) {
	return c.wait(ctx, func(ctx context.Context) (*ProcessResult, error) {
		return c.GetCreateStatus(ctx, trackingCode)
	})
}

// WaitForIssue polls GetIssueStatus until the issue request reaches Success or
// Failed, or ctx is done. See WaitForCreate.
func (c *Client) WaitForIssue(ctx context.Context, trackingCode string) (*ProcessResult, error) {
	return c.wait(ctx, func(ctx context.Context) (*ProcessResult, error) {
		return c.GetIssueStatus(ctx, trackingCode)
	})
}

func (c *Client) wait(ctx context.Context, check func(context.Context) (*ProcessResult, error)) (*ProcessResult, error) {
	for {
		res, err := check(ctx)
		if err != nil {
			return res, err
		}
		if res.Done() {
			return res, nil
		}

		timer := time.NewTimer(c.pollDelay(res))
		select {
		case <-ctx.Done():
			timer.Stop()
			return res, ctx.Err()
		case <-timer.C:
		}
	}
}

// pollDelay honours retry.next_retry_at so pending retries are not polled needlessly.
func (c *Client) pollDelay(res *ProcessResult) time.Duration {
	if res.Retry == nil {
		return c.pollInterval
	}
	next, ok := res.Retry.NextRetryTime()
	if !ok {
		return max(settlingPollInterval, c.pollInterval)
	}
	return max(time.Until(next), c.pollInterval)
}
