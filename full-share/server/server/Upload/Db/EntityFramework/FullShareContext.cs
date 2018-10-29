using Microsoft.EntityFrameworkCore;
using server.Upload.Db.EntityFramework.Model;

namespace server.Upload.Db.EntityFramework {
    // Note: generated class name was full_shareContext
    public class FullShareContext : DbContext {
        public FullShareContext() { }

        public FullShareContext(DbContextOptions<FullShareContext> options)
            : base(options) { }

        public virtual DbSet<Files> Files { get; set; }

        protected override void OnConfiguring(DbContextOptionsBuilder optionsBuilder) {
            if (!optionsBuilder.IsConfigured) {
                optionsBuilder.UseNpgsql(DbUtil.ConnStr);
            }
        }

        protected override void OnModelCreating(ModelBuilder modelBuilder) {
            modelBuilder.Entity<Files>(entity => {
                entity.HasKey(e => e.Code);

                entity.ToTable("files");

                entity.Property(e => e.Code)
                    .HasColumnName("code")
                    .HasMaxLength(36)
                    .ValueGeneratedNever();

                entity.Property(e => e.Filesize)
                    .HasColumnName("filesize")
                    .HasDefaultValueSql("0");

                entity.Property(e => e.Name)
                    .IsRequired()
                    .HasColumnName("name")
                    .HasMaxLength(50);

                entity.Property(e => e.Uploadedon).HasColumnName("uploadedon");
            });
        }
    }
}