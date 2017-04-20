package shared

type NoPluginRepositoriesError struct{}

// list outdated plugins
// can't get version info on plugins because no plugin repositories were day.

func (e NoPluginRepositoriesError) Error() string {
	return "No plugin repositories registered to search for plugin updates."
}

func (e NoPluginRepositoriesError) Translate(translate func(string, ...interface{}) string) string {
	return translate(e.Error())
}

// PluginRepositoryNotFoundError is returned when the plugin Repository was not found.
type PluginRepositoryNotFoundError struct {
	Name string
}

func (e PluginRepositoryNotFoundError) Error() string {
	return "Plugin repository '{{.RepositoryName}}' was not found."
}

func (e PluginRepositoryNotFoundError) Translate(translate func(string, ...interface{}) string) string {
	return translate(e.Error(), map[string]interface{}{"RepositoryName": e.Name})
}

// GettingPluginRepositoryError is returned when there's an error
// accessing the plugin repository
type GettingPluginRepositoryError struct {
	Name    string
	Message string
}

func (e GettingPluginRepositoryError) Error() string {
	return "Could not get plugin repository '{{.RepositoryName}}': {{.ErrorMessage}}"
}

func (e GettingPluginRepositoryError) Translate(translate func(string, ...interface{}) string) string {
	return translate(e.Error(), map[string]interface{}{"RepositoryName": e.Name, "ErrorMessage": e.Message})
}