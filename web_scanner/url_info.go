package webscanner

type URLInfo struct {
	Url             string
	NumOfCharacters int
}

type URLInfos []URLInfo

func (u URLInfos) Len() int           { return len(u) }
func (u URLInfos) Less(i, j int) bool { return u[i].Url < u[j].Url }
func (u URLInfos) Swap(i, j int)      { u[i], u[j] = u[j], u[i] }
