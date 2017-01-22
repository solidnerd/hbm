package sqlite

import (
	"fmt"
)

func (c *Config) AddHost(name string) {
	c.DB.Create(&Host{Name: name})
}

func (c *Config) RemoveHost(name string) error {
	if c.memberOfCluster(name) {
		return fmt.Errorf("host \"%s\" cannot be removed. It is being used by a cluster", name)
	}

	c.DB.Where("name = ?", name).Delete(Host{})

	return nil
}

func (c *Config) ListHosts() map[string][]string {
	var hosts []Host

	result := make(map[string][]string)

	c.DB.Find(&hosts)

	for _, host := range hosts {
		result[host.Name] = []string{}

		sql := c.DB.Table("cluster_hosts").Select("clusters.name").Joins("JOIN clusters ON clusters.id = cluster_hosts.cluster_id").Where("cluster_hosts.host_id = ?", host.ID)

		rows, _ := sql.Rows()
		defer rows.Close()

		for rows.Next() {
			var cluster string

			rows.Scan(&cluster)

			result[host.Name] = append(result[host.Name], cluster)
		}
	}

	return result
}

func (c *Config) FindHost(name string) bool {
	var count int64

	c.DB.Model(&Host{}).Where("name = ?", name).Count(&count)

	if count == 1 {
		return true
	}

	return false
}

func (c *Config) CountHost() int {
	var count int64

	c.DB.Model(&Host{}).Count(&count)

	return int(count)
}

func (c *Config) AddHostToCluster(cluster, host string) {
	cl := Cluster{}
	h := Host{}

	c.DB.Where("name = ?", host).Find(&h)
	c.DB.Where("name = ?", cluster).Find(&cl)

	c.DB.Model(&cl).Association("Hosts").Append(&h)
}

func (c *Config) RemoveHostFromCluster(cluster, host string) {
	cl := Cluster{}
	h := Host{}

	c.DB.Where("name = ?", host).Find(&h)
	c.DB.Where("name = ?", cluster).Find(&cl)

	c.DB.Model(&cl).Association("Hosts").Delete(&h)
}

func (c *Config) memberOfCluster(name string) bool {
	var count int64

	c.DB.Table("clusters").Joins("JOIN cluster_hosts ON cluster_hosts.cluster_id = clusters.id").Joins("JOIN hosts ON hosts.id = cluster_hosts.host_id").Where("hosts.name = ?", name).Count(&count)

	if count > 0 {
		return true
	}

	return false
}
