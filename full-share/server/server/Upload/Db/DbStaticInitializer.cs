using System;
using Npgsql;
using server.Upload.Db.Config;

// Uses Npgsql to setup the database/tables. It's nice to use raw sql to actually
// setup the database/schema. This way, we know exactly what is going on with it.
namespace server.Upload.Db {
    static class DbStaticInitializer {
        public static void SetupDb() {
                CreateDbIfNotExist();
        }
        
        private static string ConnStr { get; } = BuildConnStr(
            DbConfig.Config.Host,
            DbConfig.Config.Username,
            DbConfig.Config.Password,
            DbConfig.Config.DbName
        );

        private static string BuildConnStr(string host, string username, string password, string dbName) {
            return $"Host={host}; Username={username}; Password={password}; Database={dbName};";
        }

        private static void CreateDbIfNotExist() {
            string noDbConnStr = BuildConnStr(
                DbConfig.Config.Host,
                "postgres",
                DbConfig.Config.PostgresPassword,
                "postgres"
            );

            using (NpgsqlConnection conn = new NpgsqlConnection(noDbConnStr)) {
                if (DatabaseExists(conn, DbConfig.Config.DbName)) {
                    return;
                }

                using (var cmd = new NpgsqlCommand(SqlCreateDbIfNotExist, conn)) {
                    cmd.Parameters.AddWithValue("dbName", DbConfig.Config.DbName);
                    cmd.Parameters.AddWithValue("owner", DbConfig.Config.Username);
                    cmd.ExecuteNonQuery();
                }
            }
        }

        private static bool DatabaseExists(NpgsqlConnection conn, string dbName) {
            using (var cmd = new NpgsqlCommand(SqlCountPgDatabaseWhereDatname, conn)) {
                cmd.Parameters.AddWithValue("dbName", dbName);
                using (NpgsqlDataReader reader = cmd.ExecuteReader()) {
                    reader.Read();
                    int numDatabases = reader.GetInt32(0);
                    Console.WriteLine("Number of databases found: " + numDatabases);
                    return numDatabases > 0;
                }
            }
        }

        // ******************** SQL Strings and Variables ********************
        private static readonly string SqlCountPgDatabaseWhereDatname = @"SELECT COUNT(*) FROM pg_database WHERE datname = @dbName";
        private static readonly string SqlCreateDbIfNotExist = @"CREATE DATABASE @dbName OWNER @owner";
    }
}