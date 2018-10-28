namespace server.Upload.Db {
    public class SqlStr {
        public static readonly string CreateFilesTableIFNotExists = @"
                                                                            CREATE TABLE IF NOT EXISTS files(
                                                                              code VARCHAR(36) PRIMARY KEY,
                                                                              name VARCHAR(50) NOT NULL,
                                                                              fileSize INTEGER DEFAULT 0,
                                                                              uploadedOn TIMESTAMP 
                                                                            );
                                                                        ";

        public static readonly string CreateUpgradesTableIfNotExists = @"
                                                                            CREATE TABLE IF NOT EXISTS upgrades(
                                                                              upgradeNumber INTEGER PRIMARY KEY
                                                                            );
                                                                        ";

        public static readonly string InsertIntoFilesTable = @"
                                                                     INSERT INTO files
                                                                     VALUES ($1, $2, $3, $4);
                                                                  ";

        public static readonly string SelectFileInfo = @"SELECT * FROM files WHERE code = $1;";

        public static readonly string DeleteOlderThan =
            @"DELETE FROM files WHERE date_part('day', age($1, uploadedOn)) >= $2 RETURNING code;";

        public static readonly string DeleteByCode = @"DELETE FROM files WHERE code = $1";
    }
}