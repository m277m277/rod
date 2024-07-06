// This file contains the methods that panics when error return value is not nil.
// Their function names are all prefixed with Must.
// A function here is usually a wrapper for the error version with fixed default options to make it easier to use.
//
// For example the source code of [Element.Click] and [Element.MustClick]. MustClick has no argument.
// But `Click` has a `button` argument to decide which button to click.
// `MustClick` feels like a version of `Click` with some default behaviors.

package rod

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-rod/rod/lib/devices"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/utils"
	"github.com/ysmood/gson"
)

// It must be generated by genE.
type eFunc func(args ...interface{})

// Generate a eFunc with the specified fail function.
// If the last arg of eFunc is error the fail will be called.
func genE(fail func(interface{})) eFunc {
	return func(args ...interface{}) {
		err, ok := args[len(args)-1].(error)
		if ok {
			fail(err)
		}
	}
}

// WithPanic returns a browser clone with the specified panic function.
// The fail must stop the current goroutine's execution immediately, such as use [runtime.Goexit] or panic inside it.
func (b *Browser) WithPanic(fail func(interface{})) *Browser {
	n := *b
	n.e = genE(fail)
	return &n
}

// MustConnect is similar to [Browser.Connect].
func (b *Browser) MustConnect() *Browser {
	b.e(b.Connect())
	return b
}

// MustClose is similar to [Browser.Close].
func (b *Browser) MustClose() {
	_ = b.Close()
}

// MustIncognito is similar to [Browser.Incognito].
func (b *Browser) MustIncognito() *Browser {
	p, err := b.Incognito()
	b.e(err)
	return p
}

// MustPage is similar to [Browser.Page].
// The url list will be joined by "/".
func (b *Browser) MustPage(url ...string) *Page {
	p, err := b.Page(proto.TargetCreateTarget{URL: strings.Join(url, "/")})
	b.e(err)
	return p
}

// MustPages is similar to [Browser.Pages].
func (b *Browser) MustPages() Pages {
	list, err := b.Pages()
	b.e(err)
	return list
}

// MustPageFromTargetID is similar to [Browser.PageFromTargetID].
func (b *Browser) MustPageFromTargetID(targetID proto.TargetTargetID) *Page {
	p, err := b.PageFromTarget(targetID)
	b.e(err)
	return p
}

// MustHandleAuth is similar to [Browser.HandleAuth].
func (b *Browser) MustHandleAuth(username, password string) (wait func()) {
	w := b.HandleAuth(username, password)
	return func() { b.e(w()) }
}

// MustIgnoreCertErrors is similar to [Browser.IgnoreCertErrors].
func (b *Browser) MustIgnoreCertErrors(enable bool) *Browser {
	b.e(b.IgnoreCertErrors(enable))
	return b
}

// MustGetCookies is similar to [Browser.GetCookies].
func (b *Browser) MustGetCookies() []*proto.NetworkCookie {
	nc, err := b.GetCookies()
	b.e(err)
	return nc
}

// MustSetCookies is similar to [Browser.SetCookies].
// If the len(cookies) is 0 it will clear all the cookies.
func (b *Browser) MustSetCookies(cookies ...*proto.NetworkCookie) *Browser {
	if len(cookies) == 0 {
		b.e(b.SetCookies(nil))
	} else {
		b.e(b.SetCookies(proto.CookiesToParams(cookies)))
	}
	return b
}

// MustWaitDownload is similar to [Browser.WaitDownload].
// It will read the file into bytes then remove the file.
func (b *Browser) MustWaitDownload() func() []byte {
	tmpDir := filepath.Join(os.TempDir(), "rod", "downloads")
	wait := b.WaitDownload(tmpDir)

	return func() []byte {
		info := wait()
		path := filepath.Join(tmpDir, info.GUID)
		defer func() { _ = os.Remove(path) }()
		data, err := os.ReadFile(path)
		b.e(err)
		return data
	}
}

// MustVersion is similar to [Browser.Version].
func (b *Browser) MustVersion() *proto.BrowserGetVersionResult {
	v, err := b.Version()
	b.e(err)
	return v
}

// MustFind is similar to [Browser.Find].
func (ps Pages) MustFind(selector string) *Page {
	p, err := ps.Find(selector)
	if err != nil {
		if len(ps) > 0 {
			ps[0].e(err)
		} else {
			// fallback to utils.E, because we don't have enough
			// context to call the scope `.e`.
			utils.E(err)
		}
	}
	return p
}

// MustFindByURL is similar to [Page.FindByURL].
func (ps Pages) MustFindByURL(regex string) *Page {
	p, err := ps.FindByURL(regex)
	if err != nil {
		if len(ps) > 0 {
			ps[0].e(err)
		} else {
			// fallback to utils.E, because we don't have enough
			// context to call the scope `.e`.
			utils.E(err)
		}
	}
	return p
}

// WithPanic returns a page clone with the specified panic function.
// The fail must stop the current goroutine's execution immediately, such as use [runtime.Goexit] or panic inside it.
func (p *Page) WithPanic(fail func(interface{})) *Page {
	n := *p
	n.e = genE(fail)
	return &n
}

// MustInfo is similar to [Page.Info].
func (p *Page) MustInfo() *proto.TargetTargetInfo {
	info, err := p.Info()
	p.e(err)
	return info
}

// MustHTML is similar to [Page.HTML].
func (p *Page) MustHTML() string {
	html, err := p.HTML()
	p.e(err)
	return html
}

// MustCookies is similar to [Page.Cookies].
func (p *Page) MustCookies(urls ...string) []*proto.NetworkCookie {
	cookies, err := p.Cookies(urls)
	p.e(err)
	return cookies
}

// MustSetCookies is similar to [Page.SetCookies].
// If the len(cookies) is 0 it will clear all the cookies.
func (p *Page) MustSetCookies(cookies ...*proto.NetworkCookieParam) *Page {
	if len(cookies) == 0 {
		cookies = nil
	}
	p.e(p.SetCookies(cookies))
	return p
}

// MustSetExtraHeaders is similar to [Page.SetExtraHeaders].
func (p *Page) MustSetExtraHeaders(dict ...string) (cleanup func()) {
	cleanup, err := p.SetExtraHeaders(dict)
	p.e(err)
	return
}

// MustSetUserAgent is similar to [Page.SetUserAgent].
func (p *Page) MustSetUserAgent(req *proto.NetworkSetUserAgentOverride) *Page {
	p.e(p.SetUserAgent(req))
	return p
}

// MustSetBlockedURLs is similar to [Page.SetBlockedURLs].
func (p *Page) MustSetBlockedURLs(urls ...string) *Page {
	p.e(p.SetBlockedURLs(urls))
	return p
}

// MustNavigate is similar to [Page.Navigate].
func (p *Page) MustNavigate(url string) *Page {
	p.e(p.Navigate(url))
	return p
}

// MustReload is similar to [Page.Reload].
func (p *Page) MustReload() *Page {
	p.e(p.Reload())
	return p
}

// MustActivate is similar to [Page.Activate].
func (p *Page) MustActivate() *Page {
	p.e(p.Activate())
	return p
}

// MustNavigateBack is similar to [Page.NavigateBack].
func (p *Page) MustNavigateBack() *Page {
	p.e(p.NavigateBack())
	return p
}

// MustNavigateForward is similar to [Page.NavigateForward].
func (p *Page) MustNavigateForward() *Page {
	p.e(p.NavigateForward())
	return p
}

// MustGetWindow is similar to [Page.GetWindow].
func (p *Page) MustGetWindow() *proto.BrowserBounds {
	bounds, err := p.GetWindow()
	p.e(err)
	return bounds
}

// MustSetWindow is similar to [Page.SetWindow].
func (p *Page) MustSetWindow(left, top, width, height int) *Page {
	p.e(p.SetWindow(&proto.BrowserBounds{
		Left:        gson.Int(left),
		Top:         gson.Int(top),
		Width:       gson.Int(width),
		Height:      gson.Int(height),
		WindowState: proto.BrowserWindowStateNormal,
	}))
	return p
}

// MustWindowMinimize is similar to [Page.WindowMinimize].
func (p *Page) MustWindowMinimize() *Page {
	p.e(p.SetWindow(&proto.BrowserBounds{
		WindowState: proto.BrowserWindowStateMinimized,
	}))
	return p
}

// MustWindowMaximize is similar to [Page.WindowMaximize].
func (p *Page) MustWindowMaximize() *Page {
	p.e(p.SetWindow(&proto.BrowserBounds{
		WindowState: proto.BrowserWindowStateMaximized,
	}))
	return p
}

// MustWindowFullscreen is similar to [Page.WindowFullscreen].
func (p *Page) MustWindowFullscreen() *Page {
	p.e(p.SetWindow(&proto.BrowserBounds{
		WindowState: proto.BrowserWindowStateFullscreen,
	}))
	return p
}

// MustWindowNormal is similar to [Page.WindowNormal].
func (p *Page) MustWindowNormal() *Page {
	p.e(p.SetWindow(&proto.BrowserBounds{
		WindowState: proto.BrowserWindowStateNormal,
	}))
	return p
}

// MustSetViewport is similar to [Page.SetViewport].
func (p *Page) MustSetViewport(width, height int, deviceScaleFactor float64, mobile bool) *Page {
	p.e(p.SetViewport(&proto.EmulationSetDeviceMetricsOverride{
		Width:             width,
		Height:            height,
		DeviceScaleFactor: deviceScaleFactor,
		Mobile:            mobile,
	}))
	return p
}

// MustEmulate is similar to [Page.Emulate].
func (p *Page) MustEmulate(device devices.Device) *Page {
	p.e(p.Emulate(device))
	return p
}

// MustStopLoading is similar to [Page.StopLoading].
func (p *Page) MustStopLoading() *Page {
	p.e(p.StopLoading())
	return p
}

// MustClose is similar to [Page.Close].
func (p *Page) MustClose() {
	p.e(p.Close())
}

// MustHandleDialog is similar to [Page.HandleDialog].
func (p *Page) MustHandleDialog() (wait func() *proto.PageJavascriptDialogOpening, handle func(bool, string)) {
	w, h := p.HandleDialog()
	return w, func(accept bool, promptText string) {
		p.e(h(&proto.PageHandleJavaScriptDialog{
			Accept:     accept,
			PromptText: promptText,
		}))
	}
}

// MustHandleFileDialog is similar to [Page.HandleFileDialog].
func (p *Page) MustHandleFileDialog() func(...string) {
	setFiles, err := p.HandleFileDialog()
	p.e(err)
	return func(paths ...string) {
		p.e(setFiles(paths))
	}
}

// MustScreenshot is similar to [Page.Screenshot].
// If the toFile is "", it Page.will save output to "tmp/screenshots" folder, time as the file name.
func (p *Page) MustScreenshot(toFile ...string) []byte {
	bin, err := p.Screenshot(false, nil)
	p.e(err)
	p.e(saveFile(saveFileTypeScreenshot, bin, toFile))
	return bin
}

// MustCaptureDOMSnapshot is similar to [Page.CaptureDOMSnapshot].
func (p *Page) MustCaptureDOMSnapshot() (domSnapshot *proto.DOMSnapshotCaptureSnapshotResult) {
	domSnapshot, err := p.CaptureDOMSnapshot()
	p.e(err)
	return domSnapshot
}

// MustTriggerFavicon is similar to [PageTriggerFavicon].
func (p *Page) MustTriggerFavicon() *Page {
	p.e(p.TriggerFavicon())
	return p
}

// MustScreenshotFullPage is similar to [Page.ScreenshotFullPage].
// If the toFile is "", it Page.will save output to "tmp/screenshots" folder, time as the file name.
func (p *Page) MustScreenshotFullPage(toFile ...string) []byte {
	bin, err := p.Screenshot(true, nil)
	p.e(err)
	p.e(saveFile(saveFileTypeScreenshot, bin, toFile))
	return bin
}

// MustScrollScreenshotPage is similar to [Page.ScrollScreenshot].
// If the toFile is "", it Page.will save output to "tmp/screenshots" folder, time as the file name.
func (p *Page) MustScrollScreenshotPage(toFile ...string) []byte {
	bin, err := p.ScrollScreenshot(nil)
	p.e(err)
	p.e(saveFile(saveFileTypeScreenshot, bin, toFile))
	return bin
}

// MustPDF is similar to [Page.PDF].
// If the toFile is "", it Page.will save output to "tmp/pdf" folder, time as the file name.
func (p *Page) MustPDF(toFile ...string) []byte {
	r, err := p.PDF(&proto.PagePrintToPDF{})
	p.e(err)
	bin, err := io.ReadAll(r)
	p.e(err)

	p.e(saveFile(saveFileTypePDF, bin, toFile))
	return bin
}

// MustWaitOpen is similar to [Page.WaitOpen].
func (p *Page) MustWaitOpen() (wait func() (newPage *Page)) {
	w := p.WaitOpen()
	return func() *Page {
		page, err := w()
		p.e(err)
		return page
	}
}

// MustWaitNavigation is similar to [Page.WaitNavigation].
func (p *Page) MustWaitNavigation() func() {
	return p.WaitNavigation(proto.PageLifecycleEventNameNetworkAlmostIdle)
}

// MustWaitRequestIdle is similar to [Page.WaitRequestIdle].
func (p *Page) MustWaitRequestIdle(excludes ...string) (wait func()) {
	return p.WaitRequestIdle(300*time.Millisecond, nil, excludes, nil)
}

// MustWaitIdle is similar to [Page.WaitIdle].
func (p *Page) MustWaitIdle() *Page {
	p.e(p.WaitIdle(time.Minute))
	return p
}

// MustWaitDOMStable is similar to [Page.WaitDOMStable].
func (p *Page) MustWaitDOMStable() *Page {
	p.e(p.WaitDOMStable(time.Second, 0))
	return p
}

// MustWaitStable is similar to [Page.WaitStable].
func (p *Page) MustWaitStable() *Page {
	p.e(p.WaitStable(time.Second))
	return p
}

// MustWaitLoad is similar to [Page.WaitLoad].
func (p *Page) MustWaitLoad() *Page {
	p.e(p.WaitLoad())
	return p
}

// MustAddScriptTag is similar to [Page.AddScriptTag].
func (p *Page) MustAddScriptTag(url string) *Page {
	p.e(p.AddScriptTag(url, ""))
	return p
}

// MustAddStyleTag is similar to [Page.AddStyleTag].
func (p *Page) MustAddStyleTag(url string) *Page {
	p.e(p.AddStyleTag(url, ""))
	return p
}

// MustEvalOnNewDocument is similar to [Page.EvalOnNewDocument].
func (p *Page) MustEvalOnNewDocument(js string) {
	_, err := p.EvalOnNewDocument(js)
	p.e(err)
}

// MustExpose is similar to [Page.Expose].
func (p *Page) MustExpose(name string, fn func(gson.JSON) (interface{}, error)) (stop func()) {
	s, err := p.Expose(name, fn)
	p.e(err)
	return func() { p.e(s()) }
}

// MustEval is similar to [Page.Eval].
func (p *Page) MustEval(js string, params ...interface{}) gson.JSON {
	res, err := p.Eval(js, params...)
	p.e(err)
	return res.Value
}

// MustEvaluate is similar to [Page.Evaluate].
func (p *Page) MustEvaluate(opts *EvalOptions) *proto.RuntimeRemoteObject {
	res, err := p.Evaluate(opts)
	p.e(err)
	return res
}

// MustWait is similar to [Page.Wait].
func (p *Page) MustWait(js string, params ...interface{}) *Page {
	p.e(p.Wait(Eval(js, params...)))
	return p
}

// MustWaitElementsMoreThan is similar to [Page.WaitElementsMoreThan].
func (p *Page) MustWaitElementsMoreThan(selector string, num int) *Page {
	p.e(p.WaitElementsMoreThan(selector, num))
	return p
}

// MustObjectToJSON is similar to [Page.ObjectToJSON].
func (p *Page) MustObjectToJSON(obj *proto.RuntimeRemoteObject) gson.JSON {
	j, err := p.ObjectToJSON(obj)
	p.e(err)
	return j
}

// MustObjectsToJSON is similar to [Page.ObjectsToJSON].
func (p *Page) MustObjectsToJSON(list []*proto.RuntimeRemoteObject) gson.JSON {
	arr := []interface{}{}
	for _, obj := range list {
		j, err := p.ObjectToJSON(obj)
		p.e(err)
		arr = append(arr, j.Val())
	}
	return gson.New(arr)
}

// MustElementFromNode is similar to [Page.ElementFromNode].
func (p *Page) MustElementFromNode(node *proto.DOMNode) *Element {
	el, err := p.ElementFromNode(node)
	p.e(err)
	return el
}

// MustElementFromPoint is similar to [Page.ElementFromPoint].
func (p *Page) MustElementFromPoint(left, top int) *Element {
	el, err := p.ElementFromPoint(left, top)
	p.e(err)
	return el
}

// MustRelease is similar to [Page.Release].
func (p *Page) MustRelease(obj *proto.RuntimeRemoteObject) *Page {
	p.e(p.Release(obj))
	return p
}

// MustHas is similar to [Page.Has].
func (p *Page) MustHas(selector string) bool {
	has, _, err := p.Has(selector)
	p.e(err)
	return has
}

// MustHasX is similar to [Page.HasX].
func (p *Page) MustHasX(selector string) bool {
	has, _, err := p.HasX(selector)
	p.e(err)
	return has
}

// MustHasR is similar to [Page.HasR].
func (p *Page) MustHasR(selector, regex string) bool {
	has, _, err := p.HasR(selector, regex)
	p.e(err)
	return has
}

// MustSearch is similar to [Page.Search].
// It only returns the first element in the search result.
func (p *Page) MustSearch(query string) *Element {
	res, err := p.Search(query)
	p.e(err)
	res.Release()
	return res.First
}

// MustElement is similar to [Page.Element].
func (p *Page) MustElement(selector string) *Element {
	el, err := p.Element(selector)
	p.e(err)
	return el
}

// MustElementR is similar to [Page.ElementR].
func (p *Page) MustElementR(selector, jsRegex string) *Element {
	el, err := p.ElementR(selector, jsRegex)
	p.e(err)
	return el
}

// MustElementX is similar to [Page.ElementX].
func (p *Page) MustElementX(xPath string) *Element {
	el, err := p.ElementX(xPath)
	p.e(err)
	return el
}

// MustElementByJS is similar to [Page.ElementByJS].
func (p *Page) MustElementByJS(js string, params ...interface{}) *Element {
	el, err := p.ElementByJS(Eval(js, params...))
	p.e(err)
	return el
}

// MustElements is similar to [Page.Elements].
func (p *Page) MustElements(selector string) Elements {
	list, err := p.Elements(selector)
	p.e(err)
	return list
}

// MustElementsX is similar to [Page.ElementsX].
func (p *Page) MustElementsX(xpath string) Elements {
	list, err := p.ElementsX(xpath)
	p.e(err)
	return list
}

// MustElementsByJS is similar to [Page.ElementsByJS].
func (p *Page) MustElementsByJS(js string, params ...interface{}) Elements {
	list, err := p.ElementsByJS(Eval(js, params...))
	p.e(err)
	return list
}

// MustElementByJS is similar to [RaceContext.ElementByJS].
func (rc *RaceContext) MustElementByJS(js string, params []interface{}) *RaceContext {
	return rc.ElementByJS(Eval(js, params...))
}

// MustHandle is similar to [RaceContext.Handle].
func (rc *RaceContext) MustHandle(callback func(*Element)) *RaceContext {
	return rc.Handle(func(e *Element) error {
		callback(e)
		return nil
	})
}

// MustDo is similar to [RaceContext.Do].
func (rc *RaceContext) MustDo() *Element {
	el, err := rc.Do()
	rc.page.e(err)
	return el
}

// MustMoveTo is similar to [Mouse.Move].
func (m *Mouse) MustMoveTo(x, y float64) *Mouse {
	m.page.e(m.MoveTo(proto.NewPoint(x, y)))
	return m
}

// MustScroll is similar to [Mouse.Scroll].
func (m *Mouse) MustScroll(x, y float64) *Mouse {
	m.page.e(m.Scroll(x, y, 0))
	return m
}

// MustDown is similar to [Mouse.Down].
func (m *Mouse) MustDown(button proto.InputMouseButton) *Mouse {
	m.page.e(m.Down(button, 1))
	return m
}

// MustUp is similar to [Mouse.Up].
func (m *Mouse) MustUp(button proto.InputMouseButton) *Mouse {
	m.page.e(m.Up(button, 1))
	return m
}

// MustClick is similar to [Mouse.Click].
func (m *Mouse) MustClick(button proto.InputMouseButton) *Mouse {
	m.page.e(m.Click(button, 1))
	return m
}

// MustType is similar to [Keyboard.Type].
func (k *Keyboard) MustType(key ...input.Key) *Keyboard {
	k.page.e(k.Type(key...))
	return k
}

// MustDo is similar to [KeyActions.Do].
func (ka *KeyActions) MustDo() {
	ka.keyboard.page.e(ka.Do())
}

// MustInsertText is similar to [Page.InsertText].
func (p *Page) MustInsertText(text string) *Page {
	p.e(p.InsertText(text))
	return p
}

// MustStart is similar to [Touch.Start].
func (t *Touch) MustStart(points ...*proto.InputTouchPoint) *Touch {
	t.page.e(t.Start(points...))
	return t
}

// MustMove is similar to [Touch.Move].
func (t *Touch) MustMove(points ...*proto.InputTouchPoint) *Touch {
	t.page.e(t.Move(points...))
	return t
}

// MustEnd is similar to [Touch.End].
func (t *Touch) MustEnd() *Touch {
	t.page.e(t.End())
	return t
}

// MustCancel is similar to [Touch.Cancel].
func (t *Touch) MustCancel() *Touch {
	t.page.e(t.Cancel())
	return t
}

// MustTap is similar to [Touch.Tap].
func (t *Touch) MustTap(x, y float64) *Touch {
	t.page.e(t.Tap(x, y))
	return t
}

// WithPanic returns an element clone with the specified panic function.
// The fail must stop the current goroutine's execution immediately, such as use [runtime.Goexit] or panic inside it.
func (el *Element) WithPanic(fail func(interface{})) *Element {
	n := *el
	n.e = genE(fail)
	return &n
}

// MustDescribe is similar to [Element.Describe].
func (el *Element) MustDescribe() *proto.DOMNode {
	node, err := el.Describe(1, false)
	el.e(err)
	return node
}

// MustShadowRoot is similar to [Element.ShadowRoot].
func (el *Element) MustShadowRoot() *Element {
	node, err := el.ShadowRoot()
	el.e(err)
	return node
}

// MustFrame is similar to [Element.Frame].
func (el *Element) MustFrame() *Page {
	p, err := el.Frame()
	el.e(err)
	return p
}

// MustFocus is similar to [Element.Focus].
func (el *Element) MustFocus() *Element {
	el.e(el.Focus())
	return el
}

// MustScrollIntoView is similar to [Element.ScrollIntoView].
func (el *Element) MustScrollIntoView() *Element {
	el.e(el.ScrollIntoView())
	return el
}

// MustHover is similar to [Element.Hover].
func (el *Element) MustHover() *Element {
	el.e(el.Hover())
	return el
}

// MustClick is similar to [Element.Click].
func (el *Element) MustClick() *Element {
	el.e(el.Click(proto.InputMouseButtonLeft, 1))
	return el
}

// MustDoubleClick is similar to [Element.Click].
func (el *Element) MustDoubleClick() *Element {
	el.e(el.Click(proto.InputMouseButtonLeft, 2))
	return el
}

// MustTap is similar to [Element.Tap].
func (el *Element) MustTap() *Element {
	el.e(el.Tap())
	return el
}

// MustInteractable is similar to [Element.Interactable].
func (el *Element) MustInteractable() bool {
	_, err := el.Interactable()
	if errors.Is(err, &NotInteractableError{}) {
		return false
	}
	el.e(err)
	return true
}

// MustWaitInteractable is similar to [Element.WaitInteractable].
func (el *Element) MustWaitInteractable() *Element {
	el.e(el.WaitInteractable())
	return el
}

// MustType is similar to [Element.Type].
func (el *Element) MustType(keys ...input.Key) *Element {
	el.e(el.Type(keys...))
	return el
}

// MustKeyActions is similar to [Element.KeyActions].
func (el *Element) MustKeyActions() *KeyActions {
	ka, err := el.KeyActions()
	el.e(err)
	return ka
}

// MustSelectText is similar to [Element.SelectText].
func (el *Element) MustSelectText(regex string) *Element {
	el.e(el.SelectText(regex))
	return el
}

// MustSelectAllText is similar to [Element.SelectAllText].
func (el *Element) MustSelectAllText() *Element {
	el.e(el.SelectAllText())
	return el
}

// MustInput is similar to [Element.Input].
func (el *Element) MustInput(text string) *Element {
	el.e(el.Input(text))
	return el
}

// MustInputTime is similar to [Element.Input].
func (el *Element) MustInputTime(t time.Time) *Element {
	el.e(el.InputTime(t))
	return el
}

// MustInputColor is similar to [Element.InputColor].
func (el *Element) MustInputColor(color string) *Element {
	el.e(el.InputColor(color))
	return el
}

// MustBlur is similar to [Element.Blur].
func (el *Element) MustBlur() *Element {
	el.e(el.Blur())
	return el
}

// MustSelect is similar to [Element.Select].
func (el *Element) MustSelect(selectors ...string) *Element {
	el.e(el.Select(selectors, true, SelectorTypeText))
	return el
}

// MustMatches is similar to [Element.Matches].
func (el *Element) MustMatches(selector string) bool {
	res, err := el.Matches(selector)
	el.e(err)
	return res
}

// MustAttribute is similar to [Element.Attribute].
func (el *Element) MustAttribute(name string) *string {
	attr, err := el.Attribute(name)
	el.e(err)
	return attr
}

// MustProperty is similar to [Element.Property].
func (el *Element) MustProperty(name string) gson.JSON {
	prop, err := el.Property(name)
	el.e(err)
	return prop
}

// MustDisabled is similar to [Element.Disabled].
func (el *Element) MustDisabled() bool {
	disabled, err := el.Disabled()
	el.e(err)
	return disabled
}

// MustContainsElement is similar to [Element.ContainsElement].
func (el *Element) MustContainsElement(target *Element) bool {
	contains, err := el.ContainsElement(target)
	el.e(err)
	return contains
}

// MustSetFiles is similar to [Element.SetFiles].
func (el *Element) MustSetFiles(paths ...string) *Element {
	el.e(el.SetFiles(paths))
	return el
}

// MustSetDocumentContent is similar to [Page.SetDocumentContent].
func (p *Page) MustSetDocumentContent(html string) *Page {
	p.e(p.SetDocumentContent(html))
	return p
}

// MustText is similar to [Element.Text].
func (el *Element) MustText() string {
	s, err := el.Text()
	el.e(err)
	return s
}

// MustHTML is similar to [Element.HTML].
func (el *Element) MustHTML() string {
	s, err := el.HTML()
	el.e(err)
	return s
}

// MustVisible is similar to [Element.Visible].
func (el *Element) MustVisible() bool {
	v, err := el.Visible()
	el.e(err)
	return v
}

// MustWaitLoad is similar to [Element.WaitLoad].
func (el *Element) MustWaitLoad() *Element {
	el.e(el.WaitLoad())
	return el
}

// MustWaitStable is similar to [Element.WaitStable].
func (el *Element) MustWaitStable() *Element {
	el.e(el.WaitStable(300 * time.Millisecond))
	return el
}

// MustWait is similar to [Element.Wait].
func (el *Element) MustWait(js string, params ...interface{}) *Element {
	el.e(el.Wait(Eval(js, params...)))
	return el
}

// MustWaitVisible is similar to [Element.WaitVisible].
func (el *Element) MustWaitVisible() *Element {
	el.e(el.WaitVisible())
	return el
}

// MustWaitInvisible is similar to [Element.WaitInvisible]..
func (el *Element) MustWaitInvisible() *Element {
	el.e(el.WaitInvisible())
	return el
}

// MustWaitEnabled is similar to [Element.WaitEnabled].
func (el *Element) MustWaitEnabled() *Element {
	el.e(el.WaitEnabled())
	return el
}

// MustWaitWritable is similar to [Element.WaitWritable].
func (el *Element) MustWaitWritable() *Element {
	el.e(el.WaitWritable())
	return el
}

// MustShape is similar to [Element.Shape].
func (el *Element) MustShape() *proto.DOMGetContentQuadsResult {
	shape, err := el.Shape()
	el.e(err)
	return shape
}

// MustCanvasToImage is similar to [Element.CanvasToImage].
func (el *Element) MustCanvasToImage() []byte {
	bin, err := el.CanvasToImage("", -1)
	el.e(err)
	return bin
}

// MustResource is similar to [Element.Resource].
func (el *Element) MustResource() []byte {
	bin, err := el.Resource()
	el.e(err)
	return bin
}

// MustBackgroundImage is similar to [Element.BackgroundImage].
func (el *Element) MustBackgroundImage() []byte {
	bin, err := el.BackgroundImage()
	el.e(err)
	return bin
}

// MustScreenshot is similar to [Element.Screenshot].
func (el *Element) MustScreenshot(toFile ...string) []byte {
	bin, err := el.Screenshot(proto.PageCaptureScreenshotFormatPng, 0)
	el.e(err)
	el.e(saveFile(saveFileTypeScreenshot, bin, toFile))
	return bin
}

// MustRelease is similar to [Element.Release].
func (el *Element) MustRelease() {
	el.e(el.Release())
}

// MustRemove is similar to [Element.Remove].
func (el *Element) MustRemove() {
	el.e(el.Remove())
}

// MustEval is similar to [Element.Eval].
func (el *Element) MustEval(js string, params ...interface{}) gson.JSON {
	res, err := el.Eval(js, params...)
	el.e(err)
	return res.Value
}

// MustHas is similar to [Element.Has].
func (el *Element) MustHas(selector string) bool {
	has, _, err := el.Has(selector)
	el.e(err)
	return has
}

// MustHasX is similar to [Element.HasX].
func (el *Element) MustHasX(selector string) bool {
	has, _, err := el.HasX(selector)
	el.e(err)
	return has
}

// MustHasR is similar to [Element.HasR].
func (el *Element) MustHasR(selector, regex string) bool {
	has, _, err := el.HasR(selector, regex)
	el.e(err)
	return has
}

// MustElement is similar to [Element.Element].
func (el *Element) MustElement(selector string) *Element {
	el, err := el.Element(selector)
	el.e(err)
	return el
}

// MustElementX is similar to [Element.ElementX].
func (el *Element) MustElementX(xpath string) *Element {
	el, err := el.ElementX(xpath)
	el.e(err)
	return el
}

// MustElementByJS is similar to [Element.ElementByJS].
func (el *Element) MustElementByJS(js string, params ...interface{}) *Element {
	el, err := el.ElementByJS(Eval(js, params...))
	el.e(err)
	return el
}

// MustParent is similar to [Element.Parent].
func (el *Element) MustParent() *Element {
	parent, err := el.Parent()
	el.e(err)
	return parent
}

// MustParents is similar to [Element.Parents].
func (el *Element) MustParents(selector string) Elements {
	list, err := el.Parents(selector)
	el.e(err)
	return list
}

// MustNext is similar to [Element.Next].
func (el *Element) MustNext() *Element {
	parent, err := el.Next()
	el.e(err)
	return parent
}

// MustPrevious is similar to [Element.Previous].
func (el *Element) MustPrevious() *Element {
	parent, err := el.Previous()
	el.e(err)
	return parent
}

// MustElementR is similar to [Element.ElementR].
func (el *Element) MustElementR(selector, regex string) *Element {
	sub, err := el.ElementR(selector, regex)
	el.e(err)
	return sub
}

// MustElements is similar to [Element.Elements].
func (el *Element) MustElements(selector string) Elements {
	list, err := el.Elements(selector)
	el.e(err)
	return list
}

// MustElementsX is similar to [Element.ElementsX].
func (el *Element) MustElementsX(xpath string) Elements {
	list, err := el.ElementsX(xpath)
	el.e(err)
	return list
}

// MustElementsByJS is similar to [Element.ElementsByJS].
func (el *Element) MustElementsByJS(js string, params ...interface{}) Elements {
	list, err := el.ElementsByJS(Eval(js, params...))
	el.e(err)
	return list
}

// MustAdd is similar to [HijackRouter.Add].
func (r *HijackRouter) MustAdd(pattern string, handler func(*Hijack)) *HijackRouter {
	r.browser.e(r.Add(pattern, "", handler))
	return r
}

// MustRemove is similar to [HijackRouter.Remove].
func (r *HijackRouter) MustRemove(pattern string) *HijackRouter {
	r.browser.e(r.Remove(pattern))
	return r
}

// MustStop is similar to [HijackRouter.Stop].
func (r *HijackRouter) MustStop() {
	r.browser.e(r.Stop())
}

// MustLoadResponse is similar to [Hijack.LoadResponse].
func (h *Hijack) MustLoadResponse() {
	h.browser.e(h.LoadResponse(http.DefaultClient, true))
}

// MustEqual is similar to [Element.Equal].
func (el *Element) MustEqual(elm *Element) bool {
	res, err := el.Equal(elm)
	el.e(err)
	return res
}

// MustMoveMouseOut is similar to [Element.MoveMouseOut].
func (el *Element) MustMoveMouseOut() *Element {
	el.e(el.MoveMouseOut())
	return el
}

// MustGetXPath is similar to [Element.GetXPath].
func (el *Element) MustGetXPath(optimized bool) string {
	xpath, err := el.GetXPath(optimized)
	el.e(err)
	return xpath
}

// MustGet an elem from the pool. Use the [Pool[T].Put] to make it reusable later.
func (p Pool[T]) MustGet(create func() *T) *T {
	elem := <-p
	if elem == nil {
		elem = create()
	}
	return elem
}
