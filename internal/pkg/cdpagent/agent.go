package cdpagent

import (
	"context"
	"time"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/rpcc"
)

type Agent struct {
	Client     *cdp.Client
	targetID   string
	timeoutSec int
}

func newAgent(url string, timeoutSec int) (*Agent, error) {
	ctx := context.Background()

	devt := devtool.New(url)
	pt, err := devt.Create(ctx)
	if err != nil {
		return nil, err
	}

	conn, err := rpcc.DialContext(ctx, pt.WebSocketDebuggerURL)
	if err != nil {
		return nil, err
	}

	c := cdp.NewClient(conn)
	if err = c.DOM.Enable(ctx, dom.NewEnableArgs()); err != nil {
		conn.Close()
		return nil, err
	}
	if err = c.CSS.Enable(ctx); err != nil {
		conn.Close()
		return nil, err
	}
	if err = c.Page.Enable(ctx); err != nil {
		conn.Close()
		return nil, err
	}

	agent := &Agent{
		Client:     c,
		targetID:   pt.ID,
		timeoutSec: timeoutSec,
	}

	return agent, nil
}

func (agent *Agent) CreateContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(agent.timeoutSec)*time.Second)
}

func (agent *Agent) close() error {
	return agent.Client.Page.Close(context.Background())
}
