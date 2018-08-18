package heroku

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

type Addons []Addon
