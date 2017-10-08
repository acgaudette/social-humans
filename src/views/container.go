package views

// Map from a string indentifier to a view, used in templates
type Container map[string]View

/* Set views */

func (this Container) SetActive(active Active) {
	this["active"] = active
}

func (this Container) SetBase(base Base) {
	this["base"] = base
}

func (this Container) SetContent(content View) {
	this["content"] = content
}

func (this Container) SetStatus(status Status) {
	this["status"] = status
}

func NewContainer() Container {
	return make(Container)
}
