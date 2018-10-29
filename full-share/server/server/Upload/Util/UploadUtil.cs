using System;
using System.Linq;
using System.Net.WebSockets;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using Newtonsoft.Json;
using server.Upload.Controllers.Upload;

namespace server.Upload.Util {
    public static class UploadUtil {
        public static string SanitizeFileName(string fileName) {
            char[] chars = fileName.ToCharArray()
                .Select(c => char.IsLetterOrDigit(c) || c == '.' ? c : '_')
                .ToArray();
            return new string(chars);
        }
        
        public static async Task SendResp(WebSocket ws, Resp resp) {
            byte[] respBytes = Encoding.UTF8.GetBytes(JsonConvert.SerializeObject(resp));
            await ws.SendAsync(new ArraySegment<byte>(respBytes), WebSocketMessageType.Text, true,
                CancellationToken.None);
        }
    }
}