package screenshotservice

import (
	"bytes"
	"context"
	"image"
	"time"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/protocol/browser"
	"github.com/mafredri/cdp/protocol/css"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/rpcc"

	"github.com/chai2010/webp"
	"golang.org/x/image/draw"
)

const timeout = 10 * time.Second

func ScreenShot(
	url string,
	viewportWidth int,
	viewportHeight int,
	imageWidth int,
	imageHeight int,
) ([]byte, error) {
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
		Width:  &viewportWidth,
		Height: &viewportHeight,
	}))

	domContent, err := c.Page.DOMContentEventFired(ctx)
	if err != nil {
		return nil, err
	}
	defer domContent.Close()

	if err = c.Page.Enable(ctx); err != nil {
		return nil, err
	}

	frame, err := c.Page.Navigate(ctx, page.NewNavigateArgs(url))
	if err != nil {
		return nil, err
	}

	c.DOM.Enable(ctx, dom.NewEnableArgs())
	c.CSS.Enable(ctx)
	styleSheet, err := c.CSS.CreateStyleSheet(ctx, css.NewCreateStyleSheetArgs(frame.FrameID))
	if err != nil {
		return nil, err
	}
	injectedCSS := "::-webkit-scrollbar { display: none; }"
	if _, err := c.CSS.SetStyleSheetText(ctx, css.NewSetStyleSheetTextArgs(styleSheet.StyleSheetID, injectedCSS)); err != nil {
		return nil, err
	}

	select {
	case <-domContent.Ready():
	case <-ctx.Done():
	}

	screenshot, err := c.Page.CaptureScreenshot(ctx, page.NewCaptureScreenshotArgs().SetFormat("webp").SetClip(page.Viewport{
		X:      0,
		Y:      0,
		Width:  float64(viewportWidth),
		Height: float64(viewportHeight),
		Scale:  1,
	}))
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(screenshot.Data))
	if err != nil {
		return nil, err
	}
	bounds := img.Bounds()
	wRatio := float64(bounds.Dx()) / float64(imageWidth)
	hRatio := float64(bounds.Dy()) / float64(imageHeight)
	ratio := wRatio
	if hRatio > wRatio {
		ratio = hRatio
	}
	newWidth := int(float64(bounds.Dx()) / ratio)
	newHeight := int(float64(bounds.Dy()) / ratio)
	newImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	draw.CatmullRom.Scale(newImg, newImg.Bounds(), img, bounds, draw.Over, nil)

	result := new(bytes.Buffer)
	if err = webp.Encode(result, newImg, &webp.Options{Quality: 90}); err != nil {
		return nil, err
	}

	return result.Bytes(), nil
}
