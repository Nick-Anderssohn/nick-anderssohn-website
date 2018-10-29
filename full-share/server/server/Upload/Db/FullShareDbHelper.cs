using Microsoft.EntityFrameworkCore;
using server.Upload.Db.EntityFramework;

namespace server.Upload.Db {
    public sealed class FullShareDbHelper {
        public void Migrate() {
            using (var ctx = new FullShareContext()) {
                ctx.Database.Migrate();
            }
        }
    }
}