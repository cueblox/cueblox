package repo

func createSchema(path, version, name string) error {
	return nil
}

func CreateSchema(version, name string) error {
	return createSchema("", version, name)
}
