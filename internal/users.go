package internal

type User struct {
	Name     string   `yaml:"name"`
	Policies []string `yaml:"token_policies"`
	Method   string   `yaml:"method"`
}

type UserContainer struct {
	Type          string `yaml:"type"`
	UserContainer []User `yaml:"users"`
}

func NewUser(name string, policies []string, method string) (user *User, err error) {
	return &User{
		Name:     name,
		Policies: policies,
		Method:   method,
	}, nil
}

func (uc *UserContainer) AppendUser(user User) []User {
	uc.UserContainer = append(uc.UserContainer, user)
	return uc.UserContainer
}

func GetLocalUsers(path string) (users []User, dir string, err error) {
	return nil, "", nil
}
