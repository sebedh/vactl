package internal

type SshRole struct {
	Name               string   `yaml:"name"`
	Key_type           string   `yaml:"key_type"`
	Default_user       string   `yaml:"default_user"`
	Cidr_list          []string `yaml:"cidr_list"`
	Allowed_users      []string `yaml:"allowed_users"`
	Port               int      `yaml:"port"`
	Excluded_cidr_list []string `yaml:"excluded_cidr_list"`
}

type SshRoleContainer struct {
	Type             string    `yaml:"type"`
	Path             string    `yaml:"path"`
	SshRoleContainer []SshRole `yaml:"sshroles"`
}

func (r *SshRoleContainer) AppendSshRole(role SshRole) []SshRole {
	r.SshRoleContainer = append(r.SshRoleContainer, role)
	return r.SshRoleContainer
}
