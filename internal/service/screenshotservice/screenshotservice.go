package screenshotservice

import (
	"context"

	"thumburl-service/internal/config"
	"thumburl-service/internal/pkg/cdpagent"
	"thumburl-service/internal/pkg/lockmap"
	"thumburl-service/internal/pkg/logger"

	"github.com/mafredri/cdp/protocol/css"
	"github.com/mafredri/cdp/protocol/emulation"
	"github.com/mafredri/cdp/protocol/page"
)

var pool *cdpagent.Pool

func InitPool() error {
	p, err := cdpagent.InitPool(config.PoolConfig)
	pool = p
	return err
}

func DisposePool() {
	pool.Dispose()
}

func ScreenShot(ctx context.Context, url string, width int, height int) ([]byte, error) {
	agent, err := pool.GetAgent(ctx)
	if err != nil {
		return nil, err
	}
	defer pool.ReleaseAgent(ctx, agent, false)
	c := agent.Agent.Client

	agentCtx, cancelAgentCtx := agent.Agent.CreateContext()
	defer cancelAgentCtx()

	if err := c.Emulation.SetDeviceMetricsOverride(agentCtx, emulation.NewSetDeviceMetricsOverrideArgs(width, height, 1, false)); err != nil {
		return nil, err
	}
	logger.Infow(
		ctx,
		"device metrics override",
		"width", width,
		"height", height,
	)

	domContent, err := c.Page.DOMContentEventFired(agentCtx)
	if err != nil {
		return nil, err
	}
	defer domContent.Close()

	frame, err := c.Page.Navigate(agentCtx, page.NewNavigateArgs(url))
	if err != nil {
		return nil, err
	}
	logger.Infow(
		ctx,
		"navigated",
		"url", url,
	)

	styleSheet, err := c.CSS.CreateStyleSheet(agentCtx, css.NewCreateStyleSheetArgs(frame.FrameID))
	if err != nil {
		pool.ReleaseAgent(ctx, agent, true)
		return nil, err
	}
	injectedCSS := "::-webkit-scrollbar { display: none; }"
	if _, err := c.CSS.SetStyleSheetText(agentCtx, css.NewSetStyleSheetTextArgs(styleSheet.StyleSheetID, injectedCSS)); err != nil {
		pool.ReleaseAgent(ctx, agent, true)
		return nil, err
	}
	logger.Infow(
		ctx,
		"css injected",
		"stylesheet_id", styleSheet.StyleSheetID,
	)

	// wait until the DOM content is loaded, or timeout
	select {
	case <-domContent.Ready():
		logger.Infow(
			ctx,
			"dom content ready",
		)
		break
	case <-agentCtx.Done():
		logger.Infow(
			ctx,
			"timeout loading dom content",
		)
		pool.ReleaseAgent(ctx, agent, true)
		break
	}

	// for chromium, c.Page.CaptureScreenshot is not thread-safe
	lockmap.Lock(ctx, agent.Agent.DevToolURL)
	defer lockmap.Unlock(ctx, agent.Agent.DevToolURL)

	screenshot, err := c.Page.CaptureScreenshot(agentCtx, page.NewCaptureScreenshotArgs().SetFormat("webp").SetClip(page.Viewport{
		X:      0,
		Y:      0,
		Width:  float64(width),
		Height: float64(height),
		Scale:  1,
	}))
	if err != nil {
		pool.ReleaseAgent(ctx, agent, true)
		return nil, err
	}
	logger.Infow(
		ctx,
		"screenshot captured",
		"size_bytes", len(screenshot.Data),
	)

	return screenshot.Data, nil
}
