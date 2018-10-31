using System;
using System.Collections.Concurrent;
using System.IO;
using System.Threading.Tasks;
using Serilog;

namespace server.Upload.Concurrency {
    public sealed class ConcurrentFileWriter {
        private BlockingCollection<CancelableValue<ArraySegment<byte>>> _fileSlices =
            new BlockingCollection<CancelableValue<ArraySegment<byte>>>();

        private string _filePath;
        private long _targetFileSize;

        public delegate void Finally(bool success);

        private Finally _runFinally;

        public ConcurrentFileWriter(string filePath, long targetFileSize) {
            _targetFileSize = targetFileSize;
            _filePath = filePath;
        }

        public void Process(ArraySegment<byte> fileSlice) => _fileSlices.Add(CancelableValue<ArraySegment<byte>>.CreateWithValue(fileSlice));

        public void Cancel() => _fileSlices.Add(CancelableValue<ArraySegment<byte>>.CreateCanceledValue());

        public void RegisterFinally(Finally handler) => _runFinally += handler;

        public void RunAsync() => Task.Run(() => Run());

        private void Run() {
            bool success = false;
            try {
                using (FileStream writer = new FileStream(_filePath, FileMode.Create, FileAccess.Write)) {
                    int bytesWritten = 0;
                    while (bytesWritten < _targetFileSize) {
                        CancelableValue<ArraySegment<byte>> nextProcessMsg = _fileSlices.Take();
                        if (nextProcessMsg.Cancel) {
                            Log.Information("Write process was canceled");
                            // Todo: remove partial file
                            return;
                        }
                        
                        ArraySegment<byte> fileSlice = nextProcessMsg.Value;
                        writer.Write(fileSlice);
                        bytesWritten += fileSlice.Count;
                    }

                    success = bytesWritten == _targetFileSize;
                }
            }
            catch (Exception e) {
                Log.Error("could not write file {exception}", e);
            }
            finally {
                _runFinally?.Invoke(success);
            }
        }
    }
}