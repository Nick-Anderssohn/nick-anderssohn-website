using System;

namespace server.Upload {
    public class UploadConfig {
        public static long MaxFileSize { get; private set; } = 1024L * 1024 * 1024 * 5; // 5 GiB
        public static string WebsiteDomain { get; } = Environment.GetEnvironmentVariable("WEBSITE_DOMAIN");
        public static string WebsiteProtocol { get; } = Environment.GetEnvironmentVariable("WEBSITE_PROTOCOL");
    }
}