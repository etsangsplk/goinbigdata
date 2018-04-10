package tags

type Tag struct {
	TagKey   string
	TagValue string
}

type Tags struct {
	Tags []Tag
}

func GetTagKey(t Tag) string {
	return t.TagKey
}

func GetTagValue(t Tag) string {
	return t.TagValue
}
