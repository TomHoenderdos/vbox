package profile

import "github.com/TomHoenderdos/vbox/internal/config"

func init() {
	register(&Profile{
		Name:        "mysql",
		Description: "MySQL profile: installs MySQL server for development.",
		Ports:       []config.Port{{Guest: 3306, Host: 3306, Label: "MySQL"}},
		Provision: func(projectDir string) string {
			return `
    apt-get install -y mysql-server

    systemctl enable mysql
    systemctl start mysql

    # Create dev user (root/root)
    mysql -e "ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY 'root';"
    mysql -u root -proot -e "CREATE USER IF NOT EXISTS 'root'@'%' IDENTIFIED BY 'root';"
    mysql -u root -proot -e "GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' WITH GRANT OPTION;"
    mysql -u root -proot -e "FLUSH PRIVILEGES;"

    # Allow connections from host
    sed -i 's/^bind-address.*/bind-address = 0.0.0.0/' /etc/mysql/mysql.conf.d/mysqld.cnf
    systemctl restart mysql
`
		},
	})
}
