using System.IO;
using Newtonsoft.Json;

namespace server.Upload.Db.Config {
    [JsonObject(ItemRequired = Required.Always)]
    public sealed class DbConfig {
        // Singleton instance of the config
        public static DbConfig Config { get; } =
            JsonConvert.DeserializeObject<DbConfig>(File.ReadAllText("ConfigFiles/db_config.json"));
        
        public DbConfig(string username, string password, string dbName, string postgresPassword, string host) {
            Username = username;
            Password = password;
            DbName = dbName;
            PostgresPassword = postgresPassword;
            Host = host;
        }
        
        public string Username { get; }
        public string Password { get; }
        public string DbName { get; }
        public string PostgresPassword { get; }
        public string Host { get; }
    }
}