package screenshotservice

import (
	"thumburl-service/internal/config"
	"thumburl-service/internal/pkg/cdpagent"

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

func Dispose() {
	pool.Dispose()
}

func ScreenShot(url string, width int, height int) ([]byte, error) {
	agent, err := pool.GetAgent()
	if err != nil {
		return nil, err
	}
	defer pool.ReleaseAgent(agent)
	c := agent.Agent.Client

	ctx, cancel := agent.Agent.CreateContext()
	defer cancel()

	if err := c.Emulation.SetDeviceMetricsOverride(ctx, emulation.NewSetDeviceMetricsOverrideArgs(width, height, 1, false)); err != nil {
		return nil, err
	}

	domContent, err := c.Page.DOMContentEventFired(ctx)
	if err != nil {
		return nil, err
	}
	defer domContent.Close()

	frame, err := c.Page.Navigate(ctx, page.NewNavigateArgs(url))
	if err != nil {
		return nil, err
	}

	if err := c.Page.BringToFront(ctx); err != nil {
		return nil, err
	}

	styleSheet, err := c.CSS.CreateStyleSheet(ctx, css.NewCreateStyleSheetArgs(frame.FrameID))
	if err != nil {
		return nil, err
	}
	injectedCSS := "::-webkit-scrollbar { display: none; }"
	if _, err := c.CSS.SetStyleSheetText(ctx, css.NewSetStyleSheetTextArgs(styleSheet.StyleSheetID, injectedCSS)); err != nil {
		return nil, err
	}

	// wait until the DOM content is loaded, or timeout
	select {
	case <-domContent.Ready():
	case <-ctx.Done():
	}

	screenshot, err := c.Page.CaptureScreenshot(ctx, page.NewCaptureScreenshotArgs().SetFormat("webp").SetClip(page.Viewport{
		X:      0,
		Y:      0,
		Width:  float64(width),
		Height: float64(height),
		Scale:  1,
	}))
	if err != nil {
		return nil, err
	}

	return screenshot.Data, nil
}
