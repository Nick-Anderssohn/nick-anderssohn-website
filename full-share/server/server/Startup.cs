using System;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.DependencyInjection;
using server.Upload.Db;
using Serilog;

namespace server {
    public class Startup {
        static Startup() {
            Log.Logger = new LoggerConfiguration()
                .MinimumLevel.Information()
                .WriteTo.Console()
                .CreateLogger();
            
            DbStaticInitializer.SetupDb();
        }
        
        // This method gets called by the runtime. Use this method to add services to the container.
        // For more information on how to configure your application, visit https://go.microsoft.com/fwlink/?LinkID=398940
        public void ConfigureServices(IServiceCollection services) {
            services.AddMvc().SetCompatibilityVersion(CompatibilityVersion.Version_2_1);
        }

        // This method gets called by the runtime. Use this method to configure the HTTP request pipeline.
        public void Configure(IApplicationBuilder app, IHostingEnvironment env) {
            if (env.IsDevelopment()) {
                app.UseDeveloperExceptionPage();
            }
            
            // Must enable before calling UseMvc() !!!!!!
            EnableWebSockets(app);

            app.UseMvc();
        }

        private static void EnableWebSockets(IApplicationBuilder app) {
            var webSocketOptions = new WebSocketOptions {
                KeepAliveInterval = TimeSpan.FromSeconds(120),
                ReceiveBufferSize = 4 * 1024 // 4 KiB
            };

            app.UseWebSockets(webSocketOptions);
        }
    }
}