package plugin

type Meta struct {
	TableName           string
	MainTypeCapitalized string
	MainTypeLowercased  string
	Types               []Fields
}

type Fields struct {
	Name               string
	GoType             string
	MainTypeLowercased string
}
