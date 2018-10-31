using System;
using System.Collections.Generic;
using System.Collections.Immutable;
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

        // Optionally, this can be optimized with DELETE FROM files WHERE date_part('day', age($1, uploadedOn)) >= $2 RETURNING code;
        // Returns the codes of the deleted records.
        public List<string> DeleteFilesOlderThan(TimeSpan duration) {
            using (var ctx = new FullShareContext()) {
                var filesToDelete =
                    from f in ctx.Files
                    where f.Uploadedon != null && DateTime.UtcNow.Subtract(f.Uploadedon.Value) >= duration
                    select f;

                List<string> codes = filesToDelete.Select(f => f.Code).ToList();
                ctx.Files.RemoveRange(filesToDelete);
                ctx.SaveChanges();
                return codes;
            }
        }
    }
}