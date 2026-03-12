package domain

type Item struct {
	Code   string
	Name   string
	Icon   string
	Symbol string
}

type Items interface {
	Get(code string) (Item, bool)
	All() []Item
	Add(item Item)
}

type itemGroup map[string]Item

func (g itemGroup) Get(code string) (Item, bool) {
	val, ok := g[code]
	return val, ok
}

func (g itemGroup) All() []Item {
	items := make([]Item, 0, len(g))
	for _, v := range g {
		items = append(items, v)
	}
	return items
}

func (g itemGroup) Add(item Item) {
	g[item.Code] = item
}

func NewItems() Items {
	return make(itemGroup)
}
