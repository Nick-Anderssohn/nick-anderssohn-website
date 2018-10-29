using System;

namespace server.Upload.Db.EntityFramework.Model {
    public class Files {
        public string Code { get; set; }
        public string Name { get; set; }
        public long? Filesize { get; set; }
        public DateTime? Uploadedon { get; set; }
    }
}