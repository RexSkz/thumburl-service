package cdpagent

import (
	"context"
	"time"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/target"
	"github.com/mafredri/cdp/rpcc"
)

type Agent struct {
	Client     *cdp.Client
	timeoutSec time.Duration
	closeConn  func() error
}

func newAgent(url string, timeoutSec time.Duration) (*Agent, error) {
	ctx := context.Background()

	devt := devtool.New(url)
	pt, err := devt.Get(ctx, devtool.Page)
	if err != nil {
		pt, err = devt.Create(ctx)
		if err != nil {
			return nil, err
		}
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

	// open new tab
	newTarget, err := c.Target.CreateTarget(ctx, target.NewCreateTargetArgs(url))
	if err != nil {
		conn.Close()
		return nil, err
	}

	// attach to the new tab
	if _, err := c.Target.AttachToTarget(ctx, target.NewAttachToTargetArgs(newTarget.TargetID)); err != nil {
		conn.Close()
		return nil, err
	}

	agent := &Agent{
		Client:     c,
		timeoutSec: timeoutSec,
		closeConn:  conn.Close,
	}

	return agent, nil
}

func (agent *Agent) GetContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), agent.timeoutSec)
}

func (agent *Agent) close() error {
	agent.closeConn()
	return agent.Client.Browser.Close(context.Background())
}
