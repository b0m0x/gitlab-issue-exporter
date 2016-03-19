package gitlab

type ByCreationDate []GitlabComment

func (p ByCreationDate) Len() int {
	return len(p)
}

func (p ByCreationDate) Less(i, j int) bool {
	return p[i].CreatedAt.Before(p[j].CreatedAt)
}

func (p ByCreationDate) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
