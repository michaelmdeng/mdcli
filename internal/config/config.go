package config

type Config struct {
	// Whether to automatically enable using cluster-admin role for non-read-only
	// commands that require it in test kubecontexts
	EnableClusterAdminForTest bool
}

func NewConfig() Config {
	return Config{
		EnableClusterAdminForTest: true,
	}
}
