package conf

import "fmt"

func (c *serverConfig) URL() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// func (c *databaseConfig) URI() string {
// 	return fmt.Sprintf("mongodb+srv://%s:%s@%s/%s",
// 		c.Username,
// 		c.Password,
// 		c.Host,
// 		c.Name)
// }
