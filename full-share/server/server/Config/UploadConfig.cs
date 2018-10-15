using System;

namespace server.Config {
    public class UploadConfig {
        public long MaxFileSize { get; set; } = 1024L * 1024 * 1024 * 5; // 5 GiB
    }
}