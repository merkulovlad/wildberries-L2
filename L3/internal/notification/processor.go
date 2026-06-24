package notification

import "context"

type Processor struct {
	repo   Repository
	sender Sender
}

func NewProcessor(repo Repository, sender Sender) *Processor {
	return &Processor{
		repo:   repo,
		sender: sender,
	}
}

func (p *Processor) Process(ctx context.Context, id string) error {
	n, err := p.repo.Get(ctx, id)
	if err != nil {
		return err
	}

	if n.Status != StatusQueued {
		return nil
	}

	if err := p.sender.Send(ctx, n); err != nil {
		_ = p.repo.MarkFailed(ctx, id)
		return err
	}

	return p.repo.MarkSent(ctx, id)
}
