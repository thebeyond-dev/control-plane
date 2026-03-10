package ports

type Selectable interface {
	GetCode() string
	GetName() string
	GetIcon() string
}
