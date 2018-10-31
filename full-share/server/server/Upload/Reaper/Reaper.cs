using System;
using System.Collections.Generic;
using System.IO;
using System.Threading.Tasks;
using System.Timers;
using server.Upload.Db;
using Serilog;

namespace server.Upload.Reaper {
    public class Reaper {
        private readonly Timer _timer;
        private readonly FullShareDbHelper _dbHelper = new FullShareDbHelper();
        private readonly TimeSpan _timeToLive;

        public Reaper(int reapFreqHours, TimeSpan timeToLive) {
            _timer = new Timer(reapFreqHours * 60 * 60 * 1000);
            _timeToLive = timeToLive;
        }

        public void ReapAndStart() {
            Reap();
            _timer.Elapsed += Reap;
            _timer.Start();
        }

        private void Reap(object sender = null, ElapsedEventArgs e = null) {
            try {
                Log.Information("Reaping");
                List<string> codes = _dbHelper.DeleteFilesOlderThan(_timeToLive);
                Log.Information("Codes {@codes} deleted from database.", codes);

                codes.ForEach(DeleteFolderFromCode);
            } catch (Exception error) {
                Log.Error("Error when reaping: {error}", error);
            }
        }

        private static void DeleteFolderFromCode(string code) {
            Directory.Delete(Path.Combine(UploadConfig.DownloadDir, code), true);
        }
    }
}