# Go Go Gadget Musings

## Project Layout

Attempting to comply with the golang-standards/project-layout guidance as much as possible: https://github.com/golang-standards/project-layout/tree/master/docs

## Plugins

Plugins might be the right way to extend functionality to ggg, then we can make some kind of build tool that lets you select which plugins you want in your binary or have alternative implementations built for specific architectures.

I added the "gadget.Gadget" struct - implement that in your own `internal/<gadget>` folder. See `internal/{download,pivot,server,...}` for examples. When you add a gadget you need to add your gadget to the `gadgets` slice in `main.go`

### Plugin References
- https://eli.thegreenplace.net/2021/plugins-in-go/
- https://eli.thegreenplace.net/2019/design-patterns-in-gos-databasesql-package/
- https://en.wikipedia.org/wiki/Strategy_pattern
