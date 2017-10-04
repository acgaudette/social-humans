package views

type Container map[string]View

func (this Container) SetActive(active Active) {
	this["active"] = active
}

func (this Container) SetContent(content View) {
	this["content"] = content
}

func (this Container) SetStatus(status Status) {
	this["status"] = status
}

func NewContainer() Container{
	return make(Container)
}
