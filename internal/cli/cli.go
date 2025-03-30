package cli

import (
	"bufio"
	"context"
	"dockerdb/internal/databases"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dockerdb [database-type]",
	Short: "A command-line utility to set up Docker containers for various databases",
	Long:  `dockerdb is a CLI tool that simplifies the setup of Docker containers for databases like MySQL, MariaDB, PostgreSQL, MongoDB, and Redis.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to dockerdb! Please specify a database type.")
		fmt.Println("Available database types: mysql, mariadb, postgres, mongodb, redis")
		fmt.Println("Usage: dockerdb [database-type]")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(mysqlCmd)
	rootCmd.AddCommand(mariadbCmd)
	rootCmd.AddCommand(postgresCmd)
	rootCmd.AddCommand(mongodbCmd)
	rootCmd.AddCommand(redisCmd)
}

// promptForInput asks the user for input with the given prompt text
func promptForInput(prompt string, defaultValue string) string {
	reader := bufio.NewReader(os.Stdin)
	if defaultValue != "" {
		fmt.Printf("%s (%s): ", prompt, defaultValue)
	} else {
		fmt.Printf("%s: ", prompt)
	}
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue
	}
	return input
}

var mysqlCmd = &cobra.Command{
    Use:   "mysql",
    Short: "Set up a MySQL Docker container",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Setting up MySQL Docker container...")

        // Prompt for configuration
        containerName := promptForInput("Container Name", "mysql-db")
        imageTag := promptForInput("Image Tag (latest, 8.0, 5.7, etc)", "latest")
        port := promptForInput("DB Port", "3306")
        rootPassword := promptForInput("DB Root Password", "")
        dbName := promptForInput("Database Name", "mydb")
        user := promptForInput("DB User", "user")
        userPassword := promptForInput("DB User Password", "")
        volume := promptForInput("Data Volume", "mysql_data")
        network := promptForInput("Docker Network (leave empty for no specific network)", "")

        if rootPassword == "" {
            fmt.Println("Error: Root password cannot be empty")
            return
        }

        if userPassword == "" {
            fmt.Println("Error: User password cannot be empty")
            return
        }

        // Set up MySQL container
        config := databases.MySQLConfig{
            Name:         containerName,
            Image:        "mysql:" + imageTag,
            Port:         port,
            RootPassword: rootPassword,
            DatabaseName: dbName,
            User:         user,
            Password:     userPassword,
            Volume:       volume,
            Network:      network,
        }

        err := databases.SetupMySQLContainer(config)
        if err != nil {
            fmt.Printf("Error setting up MySQL container: %v\n", err)
            return
        }

        fmt.Println("MySQL container set up successfully!")
        fmt.Printf("Connection details:\n")
        fmt.Printf("  Host: localhost\n")
        fmt.Printf("  Port: %s\n", port)
        fmt.Printf("  Database: %s\n", dbName)
        fmt.Printf("  User: %s\n", user)
        if network != "" {
            fmt.Printf("  Network: %s\n", network)
        }
    },
}

var mariadbCmd = &cobra.Command{
    Use:   "mariadb",
    Short: "Set up a MariaDB Docker container",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Setting up MariaDB Docker container...")

        // Prompt for configuration
        containerName := promptForInput("Container Name", "mariadb-db")
        imageTag := promptForInput("Image Tag (latest, 10.11, 10.6, etc)", "latest")
        port := promptForInput("DB Port", "3306")
        rootPassword := promptForInput("DB Root Password", "")
        dbName := promptForInput("Database Name", "mydb")
        user := promptForInput("DB User", "user")
        userPassword := promptForInput("DB User Password", "")
        volume := promptForInput("Data Volume", "mariadb_data")
        network := promptForInput("Docker Network (leave empty for no specific network)", "")

        if rootPassword == "" {
            fmt.Println("Error: Root password cannot be empty")
            return
        }

        if userPassword == "" {
            fmt.Println("Error: User password cannot be empty")
            return
        }

        // Set up MariaDB container
        config := databases.MariaDBConfig{
            Name:         containerName,
            Image:        "mariadb:" + imageTag,
            Port:         port,
            RootPassword: rootPassword,
            DatabaseName: dbName,
            User:         user,
            Password:     userPassword,
            Volume:       volume,
            Network:      network,
        }

        err := databases.SetupMariaDBContainer(config)
        if err != nil {
            fmt.Printf("Error setting up MariaDB container: %v\n", err)
            return
        }

        fmt.Println("MariaDB container set up successfully!")
        fmt.Printf("Connection details:\n")
        fmt.Printf("  Host: localhost\n")
        fmt.Printf("  Port: %s\n", port)
        fmt.Printf("  Database: %s\n", dbName)
        fmt.Printf("  User: %s\n", user)
        if network != "" {
            fmt.Printf("  Network: %s\n", network)
        }
    },
}

var postgresCmd = &cobra.Command{
    Use:   "postgres",
    Short: "Set up a PostgreSQL Docker container",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Setting up PostgreSQL Docker container...")

        // Prompt for configuration
        containerName := promptForInput("Container Name", "postgres-db")
        imageTag := promptForInput("Image Tag (latest, 16, 15, 14, etc)", "latest")
        port := promptForInput("DB Port", "5432")
        dbName := promptForInput("Database Name", "postgres")
        user := promptForInput("DB User", "postgres")
        password := promptForInput("DB User Password", "")
        volume := promptForInput("Data Volume", "postgres_data")
        network := promptForInput("Docker Network (leave empty for no specific network)", "")

        if password == "" {
            fmt.Println("Error: Password cannot be empty")
            return
        }

        // Set up PostgreSQL container
        config := databases.PostgresConfig{
            Name:     containerName,
            Image:    "postgres:" + imageTag,
            Port:     port,
            User:     user,
            Password: password,
            Database: dbName,
            Volume:   volume,
            Network:  network,
        }

        err := databases.SetupPostgresContainer(config)
        if err != nil {
            fmt.Printf("Error setting up PostgreSQL container: %v\n", err)
            return
        }

        fmt.Println("PostgreSQL container set up successfully!")
        fmt.Printf("Connection details:\n")
        fmt.Printf("  Host: localhost\n")
        fmt.Printf("  Port: %s\n", port)
        fmt.Printf("  Database: %s\n", dbName)
        fmt.Printf("  User: %s\n", user)
        if network != "" {
            fmt.Printf("  Network: %s\n", network)
        }
    },
}

var mongodbCmd = &cobra.Command{
    Use:   "mongodb",
    Short: "Set up a MongoDB Docker container",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Setting up MongoDB Docker container...")

        // Prompt for configuration
        containerName := promptForInput("Container Name", "mongodb")
        imageTag := promptForInput("Image Tag (latest, 7.0, 6.0, 5.0, etc)", "latest")
        port := promptForInput("DB Port", "27017")
        volume := promptForInput("Data Volume", "mongodb_data")
        network := promptForInput("Docker Network (leave empty for no specific network)", "")

        useAuth := promptForInput("Enable Authentication? (yes/no)", "no")
        var user, password string
        if strings.ToLower(useAuth) == "yes" {
            user = promptForInput("Admin Username", "admin")
            password = promptForInput("Admin Password", "")

            if password == "" {
                fmt.Println("Error: Admin password cannot be empty when authentication is enabled")
                return
            }
        }

        // Set up MongoDB container
        config := &databases.MongoDBConfig{
            Name:     containerName,
            Image:    "mongo:" + imageTag,
            Port:     port,
            Volume:   volume,
            User:     user,
            Password: password,
            Auth:     strings.ToLower(useAuth) == "yes",
            Network:  network,
        }

        ctx := context.Background()
        err := databases.SetupMongoDB(ctx, config)
        if err != nil {
            fmt.Printf("Error setting up MongoDB container: %v\n", err)
            return
        }

        fmt.Println("MongoDB container set up successfully!")
        fmt.Printf("Connection details:\n")
        fmt.Printf("  Host: localhost\n")
        fmt.Printf("  Port: %s\n", port)
        if strings.ToLower(useAuth) == "yes" {
            fmt.Printf("  User: %s\n", user)
        }
        if network != "" {
            fmt.Printf("  Network: %s\n", network)
        }
    },
}

var redisCmd = &cobra.Command{
    Use:   "redis",
    Short: "Set up a Redis Docker container",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Setting up Redis Docker container...")

        // Prompt for configuration
        containerName := promptForInput("Container Name", "redis")
        imageTag := promptForInput("Image Tag (latest, 7.2, 7.0, alpine, etc)", "latest")
        port := promptForInput("DB Port", "6379")
        volume := promptForInput("Data Volume", "redis_data")
        password := promptForInput("Password (optional)", "")
        network := promptForInput("Docker Network (leave empty for no specific network)", "")

        // Set up Redis container
        config := &databases.RedisConfig{
            Name:     containerName,
            Image:    "redis:" + imageTag,
            Port:     port,
            Volume:   volume,
            Password: password,
            Network:  network,
        }

        err := databases.SetupRedisContainer(config)
        if err != nil {
            fmt.Printf("Error setting up Redis container: %v\n", err)
            return
        }

        fmt.Println("Redis container set up successfully!")
        fmt.Printf("Connection details:\n")
        fmt.Printf("  Host: localhost\n")
        fmt.Printf("  Port: %s\n", port)
        if password != "" {
            fmt.Printf("  Password: (configured)\n")
        }
        if network != "" {
            fmt.Printf("  Network: %s\n", network)
        }
    },
}