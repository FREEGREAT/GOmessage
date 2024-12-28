package database

import "github.com/sirupsen/logrus"

func Exec(query string, values ...interface{}) error {
	if err := connection.session.Query(query).Bind(values...).Exec(); err != nil {
		logrus.Fatal(err)
		return err
	}
	return nil
}
