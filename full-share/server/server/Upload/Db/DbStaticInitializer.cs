using System;
using Npgsql;
using server.Upload.Db.Config;
using Serilog;

// Uses Npgsql to setup the database/tables. It's nice to use raw sql to actually
// setup the database/schema. This way, we know exactly what is going on with it.
namespace server.Upload.Db {
    static class DbStaticInitializer {
        public static void SetupDb() {
            CreateDbIfNotExist();
        }

        private static void CreateDbIfNotExist() {
            string noDbConnStr = DbUtil.BuildConnStr(
                DbConfig.Config.Host,
                "postgres",
                DbConfig.Config.PostgresPassword,
                "postgres"
            );

            using (NpgsqlConnection conn = new NpgsqlConnection(noDbConnStr)) {
                conn.Open();
                if (DatabaseExists(conn, DbConfig.Config.DbName)) {
                    return;
                }

                // Don't need to sanitize values from our config...
                using (var cmd =
                    new NpgsqlCommand($"CREATE DATABASE {DbConfig.Config.DbName} OWNER {DbConfig.Config.Username}",
                        conn)) {
                    cmd.ExecuteNonQuery();
                }
            }
        }

        private static bool DatabaseExists(NpgsqlConnection conn, string dbName) {
            using (var cmd = new NpgsqlCommand(SqlCountPgDatabaseWhereDatname, conn)) {
                cmd.Parameters.AddWithValue("@dbName", dbName);
                using (NpgsqlDataReader reader = cmd.ExecuteReader()) {
                    reader.Read();
                    int numDatabases = reader.GetInt32(0);
                    Log.Information("Number of databases found: {NumDatabases}", numDatabases);
                    return numDatabases > 0;
                }
            }
        }

        // ******************** SQL Strings and Variables ********************
        private static readonly string SqlCountPgDatabaseWhereDatname =
            @"SELECT COUNT(*) FROM pg_database WHERE datname = @dbName";
    }
}