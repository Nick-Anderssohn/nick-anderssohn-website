using System.Linq;
using Microsoft.EntityFrameworkCore;
using server.Upload.Db.EntityFramework;
using server.Upload.Db.EntityFramework.Model;

namespace server.Upload.Db {
    public sealed class FullShareDbHelper {
        public void Migrate() {
            using (var ctx = new FullShareContext()) {
                ctx.Database.Migrate();
            }
        }

        public void InsertFilesEntry(Files fileInfo) {
            using (var ctx = new FullShareContext()) {
                ctx.Files.Add(fileInfo);
                ctx.SaveChanges();
            }
        }

        // Throws an exception if a matching file is not found.
        public Files GetFileInfo(string code) {
            using (var ctx = new FullShareContext()) {
                return ctx.Files.Single(file => file.Code == code);
            }
        }

        public bool FileInfoExists(string code) {
            using (var ctx = new FullShareContext()) {
                return ctx.Files.Any(file => file.Code == code);
            }
        }

        public void DeleteFileInfo(string code) {
            Files fileInfo = new Files {Code = code};
            DeleteFileInfo(fileInfo);
        }

        public void DeleteFileInfo(Files fileInfo) {
            using (var ctx = new FullShareContext()) {
                ctx.Entry(fileInfo).State = EntityState.Deleted;
                ctx.SaveChanges();
            }
        }
    }
}