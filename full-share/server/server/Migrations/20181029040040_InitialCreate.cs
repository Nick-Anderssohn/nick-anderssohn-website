using System;
using Microsoft.EntityFrameworkCore.Migrations;

namespace server.Migrations {
    public partial class InitialCreate : Migration {
        protected override void Up(MigrationBuilder migrationBuilder) {
            migrationBuilder.CreateTable(
                name: "files",
                columns: table => new {
                    code = table.Column<string>(maxLength: 36, nullable: false),
                    name = table.Column<string>(maxLength: 50, nullable: false),
                    filesize = table.Column<long>(nullable: true, defaultValueSql: "0"),
                    uploadedon = table.Column<DateTime>(nullable: true)
                },
                constraints: table => { table.PrimaryKey("PK_files", x => x.code); });
        }

        protected override void Down(MigrationBuilder migrationBuilder) {
            migrationBuilder.DropTable(
                name: "files");
        }
    }
}