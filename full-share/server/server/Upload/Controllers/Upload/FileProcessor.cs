using System;
using System.IO;
using System.Net.WebSockets;
using System.Threading;
using System.Threading.Tasks;
using server.Upload.Concurrency;
using server.Upload.Util;
using Serilog;

namespace server.Upload.Controllers.Upload {
    sealed class FileProcessor {
        private static readonly string DownloadDir = "download";
        private readonly ConcurrentFileWriter _writer;
        private readonly WebSocket _ws;
        private readonly string _code;
        private readonly string _fileName;
        private readonly long _targetFileSize;

        public FileProcessor(WebSocket ws, long targetFileSize, string code, string fileName) {
            _ws = ws;
            _code = code;
            _fileName = fileName;
            _targetFileSize = targetFileSize;
            string dir = Path.Combine(DownloadDir, code);
            Directory.CreateDirectory(dir);
            _writer = new ConcurrentFileWriter(Path.Combine(dir, fileName), targetFileSize);
            _writer.RegisterFinally(Finish);
        }

        public async Task Run() {
            try {
                _writer.RunAsync();

                for (int bytesReceived = 0; bytesReceived < _targetFileSize;) {
                    var buf = new byte[1024 * 1024]; // 1 MiB
                    WebSocketReceiveResult result =
                        await _ws.ReceiveAsync(new ArraySegment<byte>(buf), CancellationToken.None);
                    bytesReceived += result.Count;
                    _writer.Process(new ArraySegment<byte>(buf, 0, result.Count));
                }
            }
            catch (Exception e) {
                _writer.Cancel();
                Log.Error("could not process file {exception}", e);
            }
        }

        private void Finish(bool success) {
            if (success) {
                HandleSuccess();
            }
            else {
                HandleUnknownFailure();
            }
        }

        private async void HandleSuccess() {
            await UploadUtil.SendResp(_ws, Resp.Ok(MakeDownloadUrl()));
            await _ws.CloseAsync(WebSocketCloseStatus.NormalClosure, "Success", CancellationToken.None);
        }

        private async void HandleUnknownFailure() {
            await UploadUtil.SendResp(_ws, Resp.InternalError());
            await _ws.CloseAsync(WebSocketCloseStatus.InternalServerError, "Unknown Error", CancellationToken.None);
        }

        private string MakeDownloadUrl() {
            string protocol = UploadConfig.WebsiteProtocol;
            string domain = UploadConfig.WebsiteDomain;
            string downloadDir = Uri.EscapeDataString(DownloadDir);
            string code = Uri.EscapeDataString(_code);
            string fileName = Uri.EscapeDataString(_fileName);
            return $"{protocol}{domain}/{downloadDir}/{code}/{fileName}";
        }
    }
}