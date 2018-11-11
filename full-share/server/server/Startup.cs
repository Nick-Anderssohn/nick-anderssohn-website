using System;
using System.Net;
using System.Net.Mime;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.HttpOverrides;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.DependencyInjection;
using server.Upload.Controllers.Upload;
using server.Upload.Db;
using server.Upload.Db.Config;
using server.Upload.Reaper;
using Serilog;

namespace server {
    public class Startup {
        // Every 24 hours, remove files 2 days or older.
        private static readonly Reaper Grim = new Reaper(24, TimeSpan.FromDays(2));
        static Startup() {
            Log.Logger = new LoggerConfiguration()
                .MinimumLevel.Information()
                .WriteTo.Console()
                .CreateLogger();

            DbStaticInitializer.SetupDb();
            var dbHelper = new FullShareDbHelper();
            dbHelper.Migrate();
            
            Grim.ReapAndStart();
        }

        // This method gets called by the runtime. Use this method to add services to the container.
        // For more information on how to configure your application, visit https://go.microsoft.com/fwlink/?LinkID=398940
        public void ConfigureServices(IServiceCollection services) {
            services.AddSingleton<FullShareDbHelper, FullShareDbHelper>();
            services.AddMvc().SetCompatibilityVersion(CompatibilityVersion.Version_2_1);
        }

        // This method gets called by the runtime. Use this method to configure the HTTP request pipeline.
        public void Configure(IApplicationBuilder app, IHostingEnvironment env) {
            if (env.IsDevelopment()) {
                app.UseDeveloperExceptionPage();
            }

            // Must use websockets before calling UseMvc() !!!!!!
            UseWebSockets(app);

            app.UseMvc();
        }

        private static void UseWebSockets(IApplicationBuilder app) {
            var webSocketOptions = new WebSocketOptions {
                KeepAliveInterval = TimeSpan.FromSeconds(120),
                ReceiveBufferSize = 4 * 1024 // 4 KiB
            };

            app.UseWebSockets(webSocketOptions);
        }
    }
}