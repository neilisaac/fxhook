fxhook
======

Convenience method to create a [fx](https://go.uber.org/fx) Lifecycle hook
to run a background goroutine of the form `func(context.Context) error`
which is started via Hook.OnStart and is canceled via Hook.OnStop.


## Warning

This package is currently experimental.
Its API or semantics are subject to change until version 1.0.


## Example

```go
func DoSomethingUsefulUntilCtxCanceled(ctx context.Context) error {
  ...
}

func Register(lc fx.Lifecycle) {
  lc.Append(fxhook.CtxErrFunc(DoSomethingUsefulUntilCtxCanceled))
}

func main() {
  app := fx.New(fx.Invoke(Register)).Run()
}
```
