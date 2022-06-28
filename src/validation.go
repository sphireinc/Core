package core

// validate helps to validate our Application
func (c *Config) validate() error {
	//return ozzo.ValidateStruct(c,
	//ozzo.Field(&c.Application.Name, ozzo.Required, ozzo.Length(2, 50)),
	//ozzo.Field(&c.Application.UUID, ozzo.Required, ozzo.Length(1, 0), is.UUID),
	//ozzo.Field(&c.Application.Version, ozzo.Required, is.Semver),
	//ozzo.Field(&c.Log, ozzo.Required),
	//ozzo.Field(&c.Application.Runtime, ozzo.Required),
	//)
	return nil
}
