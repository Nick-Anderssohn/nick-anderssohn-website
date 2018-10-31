using System;

namespace server.Upload.Db.EntityFramework.Model {
    public class Files {
        // Code is set as the primary key in FullShareContext
        public string Code { get; set; }
        public string Name { get; set; }
        public long? Filesize { get; set; }
        public DateTime? Uploadedon { get; set; }

        public static Files CreateWithUtcNow(string code, string name, long fileSize) {
            return new Files {
                Code = code,
                Name = name,
                Filesize = fileSize,
                Uploadedon = DateTime.UtcNow
            };
        }
    }
}