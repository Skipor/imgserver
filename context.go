package imgserver

import (
	"errors"
	"net/http"
	"net/url"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

type ctxValueKeyType int

//private keys
const (
	ctxURLParamKey ctxValueKeyType = iota
)

// public keys upper handler can
const (
	CtxHTTPClientKey = "httpclient"
	CtxLoggerKey     = "logger"
)

func cxtAwareGet(ctx context.Context, URL string) (*http.Response, error) {
	// request will be canceled on context cancel or timeout
	return ctxhttp.Get(ctx, getClient(ctx), URL)

	// another way to do context-aware request.
	// Way to set req.Cancel = ctx.Done seems have better performance, but return not ctx.Err() on ctx.Done
	//req, err := http.NewRequest("GET", reqUrl, nil)
	//req.Cancel = ctx.Done()
	//if err != nil {
	//	return nil, &HandlerError{500, "Can't get requested page", err}
	//}
	//return h.Client.Do(req)
}

func newImgLogicContext(ctx context.Context, client *http.Client, urlParam *url.URL) context.Context {
	//don't override passed context
	ctx = context.WithValue(ctx, CtxLoggerKey, SetEmitter(getLogger(ctx), "ImgLogicHandler"))
	if _, ok := ctx.Value(CtxHTTPClientKey).(*http.Client); !ok {
		ctx = context.WithValue(ctx, CtxHTTPClientKey, client)
	}
	ctx = context.WithValue(ctx, ctxURLParamKey, urlParam)
	return ctx
}

func getLogger(ctx context.Context) Logger {
	log, ok := ctx.Value(CtxLoggerKey).(Logger)
	if !ok {
		panic(errors.New("No logger in context"))
	}
	return log
}

func setLogger(ctx context.Context, log Logger) context.Context {
	return context.WithValue(ctx, CtxLoggerKey, log)
}

func getClient(ctx context.Context) *http.Client {
	client, ok := ctx.Value(CtxHTTPClientKey).(*http.Client)
	if !ok {
		panic(errors.New("No client in context"))
	}
	return client
}

func getURLParam(ctx context.Context) *url.URL {
	urlParam, ok := ctx.Value(ctxURLParamKey).(*url.URL)
	if !ok {
		panic(errors.New("No urlParam in context"))
	}
	return urlParam
}

func getLocalLogger(ctx context.Context, emitter string) Logger {
	return SetEmitter(getLogger(ctx), emitter)
}
