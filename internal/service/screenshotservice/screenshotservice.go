package screenshotservice

import (
	"context"
	"time"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/protocol/browser"
	"github.com/mafredri/cdp/protocol/css"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/rpcc"
)

const timeout = 30 * time.Second

func ScreenShot(url string, width int, height int) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	devt := devtool.New("http://127.0.0.1:9222")
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
	defer conn.Close()

	c := cdp.NewClient(conn)
	c.Browser.SetWindowBounds(ctx, browser.NewSetWindowBoundsArgs(browser.WindowID(1), browser.Bounds{
		Width:  &width,
		Height: &height,
	}))
	c.CSS.CreateStyleSheet(ctx, css.NewCreateStyleSheetArgs("::-webkit-scrollbar { display: none; }"))

	domContent, err := c.Page.DOMContentEventFired(ctx)
	if err != nil {
		return nil, err
	}
	defer domContent.Close()

	if _, err := c.Page.Navigate(ctx, page.NewNavigateArgs(url)); err != nil {
		return nil, err
	}

	// if _, err = domContent.Recv(); err != nil {
	// 	return nil, err
	// }
	time.Sleep(10 * time.Second)

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
