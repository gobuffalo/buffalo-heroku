package heroku

import "encoding/json"

var DefaultAddons = Addons{
	{"heroku-postgresql", []string{"hobby-dev", "hobby-basic", "standard-0"}, "hobby-dev"},
	{"sendgrid", []string{"starter"}, ""},
	{"heroku-redis", []string{"hobby-dev"}, ""},
}

type Addon struct {
	Name      string
	Available []string
	Level     string
}

func (a Addon) String() string {
	b, _ := json.Marshal(a)
	return string(b)
}

type Addons []Addon
