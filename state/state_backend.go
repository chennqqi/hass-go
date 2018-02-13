package state

type stateDB interface {
	open() error
	close()
	setString(key string, value string) error
	getString(key string) (value string, err error)
}

type ContextStateDB struct {
	db stateDB
}

func NewDB(configjson string) *ContextStateDB {
	// Load JSON config
	// Determine the type of DB to create
	// Create the StateDB instance and initialize it
}

func (c *ContextStateDB) Open() error {
	return c.db.open()
}

func (c *ContextStateDB) Close() {
	c.db.close()
}

func (c *ContextStateDB) SetString(key string, value string) error {
	return c.db.setString(key, value)
}

func (c *ContextStateDB) GetString(key string) (value string, err error) {
	return c.db.getString(key)
}

// SetInt
// GetInt
// SetFloat
// GetFloat
// SetTime
// GetTime
