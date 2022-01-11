package plugins

type PluginDefinition interface {
	Execute() error
}
