package database

type Option func(*Database)

func WithEncryptedConnection() Option {
	return func(db *Database) { db.encrypted = true }
}

func WithHost(host string) Option {
	return func(db *Database) { db.host = host }
}

func WithPort(port int) Option {
	return func(db *Database) { db.port = port }
}

func WithoutPort() Option {
	return func(db *Database) { db.port = 0 }
}

func WithCredentials(username, password string) Option {
	return func(db *Database) {
		db.username = &username
		db.password = &password
	}
}
