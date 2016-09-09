package tmplbox

type AssetSourceFunc func(string) ([]byte, error)

func (f AssetSourceFunc) Get(s string) ([]byte, error) {
	return f(s)
}
