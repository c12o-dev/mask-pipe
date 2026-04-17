package patterns

import "regexp"

var dbURLPassword = &Pattern{
	ID:          "db_url_password",
	Name:        "DB URL password",
	Regex:       regexp.MustCompile(`://[^:/\s]+:([^@/\s]+)@`),
	CaptureIdx:  1,
	Replacement: "****",
	Examples: []string{
		"postgres://admin:s3cretP4ss@db.example.com:5432/mydb",
		"mysql://root:hunter2@localhost:3306/app",
		"mongodb://user:p%40ssw0rd@cluster.mongodb.net/db",
		"redis://default:myredispass@redis.internal:6379",
		"amqp://guest:guest@rabbitmq.local:5672/vhost",
	},
	NonExamples: []string{
		"https://example.com/path:8080/foo@bar",
		"https://example.com:8080/api/v1/@latest",
		"git@github.com:org/repo.git",
		"mailto:user@example.com",
		"https://example.com/@user/profile",
		"file:///home/user/data.db",
	},
	Source: "https://datatracker.ietf.org/doc/html/rfc3986#section-3.2.1",
}
