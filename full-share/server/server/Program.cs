using System;
using Microsoft.AspNetCore;
using Microsoft.AspNetCore.Hosting;
using server.Upload.Reaper;

namespace server {
    public class Program {
        public static void Main(string[] args) {
            CreateWebHostBuilder(args).Build().Run();
        }

        public static IWebHostBuilder CreateWebHostBuilder(string[] args) =>
            WebHost.CreateDefaultBuilder(args)
                .UseUrls("http://0.0.0.0:8080")
                .UseKestrel()
                .UseIISIntegration() // should I be using this?
                .UseStartup<Startup>();
    }
}