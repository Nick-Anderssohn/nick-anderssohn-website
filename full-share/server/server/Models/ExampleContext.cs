using Microsoft.EntityFrameworkCore;

namespace server.Models
{
    public class ExampleContext : DbContext
    {
        public DbSet<ExampleModel> Examples { get; set; }
        
        public ExampleContext(DbContextOptions<ExampleContext> options) : base(options)
        {
            
        }
    }
}