package screenshotservice

import (
	"fmt"
	"thumburl-service/internal/config"
	"thumburl-service/internal/pkg/cdpagent"
	"thumburl-service/internal/pkg/lockmap"

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
	fmt.Printf("set device metrics override %d %d\n", width, height)

	domContent, err := c.Page.DOMContentEventFired(ctx)
	if err != nil {
		return nil, err
	}
	defer domContent.Close()

	frame, err := c.Page.Navigate(ctx, page.NewNavigateArgs(url))
	if err != nil {
		return nil, err
	}
	fmt.Printf("navigated to %s\n", url)

	styleSheet, err := c.CSS.CreateStyleSheet(ctx, css.NewCreateStyleSheetArgs(frame.FrameID))
	if err != nil {
		return nil, err
	}
	injectedCSS := "::-webkit-scrollbar { display: none; }"
	if _, err := c.CSS.SetStyleSheetText(ctx, css.NewSetStyleSheetTextArgs(styleSheet.StyleSheetID, injectedCSS)); err != nil {
		return nil, err
	}
	fmt.Printf("injected css\n")

	// wait until the DOM content is loaded, or timeout
	select {
	case <-domContent.Ready():
		fmt.Printf("dom content ready\n")
		break
	case <-ctx.Done():
		fmt.Printf("timeout waiting for dom content\n")
		break
	}

	lockmap.Lock(agent.Agent.DevToolURL)
	defer lockmap.Unlock(agent.Agent.DevToolURL)

	if err := c.Page.BringToFront(ctx); err != nil {
		return nil, err
	}
	fmt.Printf("bring to front\n")

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
	fmt.Printf("captured screenshot\n")

	return screenshot.Data, nil
}
