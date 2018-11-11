using server.Upload.Db.Config;

namespace server.Upload.Db {
    public class DbUtil {
        public static string ConnStr { get; } = BuildConnStr(
            DbConfig.Config.Host,
            DbConfig.Config.Username,
            DbConfig.Config.Password,
            DbConfig.Config.DbName
        );

        public static string BuildConnStr(string host, string username, string password, string dbName) {
            return $"Host={host}; Username={username}; Password={password}; Database={dbName};";
        }
    }
}