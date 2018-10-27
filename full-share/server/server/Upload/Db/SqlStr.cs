namespace server.Upload.Db {
    public class SqlStr {
        private static readonly string CreateFilesTableIFNotExistsSql = @"
                                                                            CREATE TABLE IF NOT EXISTS files(
                                                                              code VARCHAR(36) PRIMARY KEY,
                                                                              name VARCHAR(50) NOT NULL,
                                                                              fileSize INTEGER DEFAULT 0,
                                                                              uploadedOn TIMESTAMP 
                                                                            );
                                                                        ";

        private static readonly string CreateUpgradesTableIfNotExists = @"
                                                                            CREATE TABLE IF NOT EXISTS upgrades(
                                                                              upgradeNumber INTEGER PRIMARY KEY
                                                                            );
                                                                        ";

        private static readonly string SelectFileInfoSql = @"                SELECT * FROM files WHERE code = $1;";

        private static readonly string DeleteOlderThan = @"                 DELETE FROM files WHERE date_part('day', age($1, uploadedOn)) >= $2 RETURNING code;";

        private static readonly string DeleteByCode = @"                    DELETE FROM files WHERE code = $1";
    }
}