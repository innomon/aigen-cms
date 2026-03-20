package descriptors

type PublicationStatus string

const (
	Draft       PublicationStatus = "Draft"
	Published   PublicationStatus = "Published"
	Unpublished PublicationStatus = "Unpublished"
	Scheduled   PublicationStatus = "Scheduled"
)

type Entity struct {
	Attributes               []Attribute       `json:"attributes"`
	Name                     string            `json:"name"`
	DisplayName              string            `json:"displayName"`
	TableName                string            `json:"tableName"`
	LabelAttributeName       string            `json:"labelAttributeName"`
	PrimaryKey               string            `json:"primaryKey"`
	DefaultPageSize          int               `json:"defaultPageSize"`
	DefaultPublicationStatus PublicationStatus `json:"defaultPublicationStatus"`
	PageUrl                  string            `json:"pageUrl"`
	TagsQuery                string            `json:"tagsQuery"`
	TagsQueryParam           string            `json:"tagsQueryParam"`
	TitleTagField            string            `json:"titleTagField"`
	SubtitleTagField         string            `json:"subtitleTagField"`
	ContentTagField          string            `json:"contentTagField"`
	ImageTagField            string            `json:"imageTagField"`
	PublishTimeTagField      string            `json:"publishTimeTagField"`
}

type LoadedEntity struct {
	Entity
	LoadedAttributes           []LoadedAttribute
	PrimaryKeyAttribute        LoadedAttribute
	LabelAttribute             LoadedAttribute
	DeletedAttribute           LoadedAttribute
	PublicationStatusAttribute LoadedAttribute
	UpdatedAtAttribute         LoadedAttribute
}
