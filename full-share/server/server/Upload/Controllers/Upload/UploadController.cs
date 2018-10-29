using System;
using System.Net.WebSockets;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using Newtonsoft.Json;
using server.Config;
using server.Upload.Util;

namespace server.Upload.Controllers.Upload {
    /*
    A successful save file hand shake will look like the following:
    Client sends:
        {
            "FileSize": 1000,
            "FileName": "Chicken Dinner"
        }
    Server replies:
	{
        "StatusCode": 200,
		"StatusMsg": "Ok",
        "Message": ""
    }
    Then in a loop until file size is reached:
        Client sends file bytes.
        The server responds with
        {
            "StatusCode": 200,
            "StatusMsg": "Ok",
            "Message": ""
        }
    Then the server will finish up with:
        {
            "StatusCode": 200,
            "StatusMsg": "Ok",
            "Message": "The url here"
        }
    At any point in time, the server might respond with status indicating an error code and message containing an
    error message.
    */
    [Route("/")]
    [ApiController]
    public class UploadController : ControllerBase {
        private readonly UploadConfig _config = new UploadConfig();

        [HttpGet]
        public async Task HandleWsUpload() {
            // Setup ws connection
            if (!HttpContext.WebSockets.IsWebSocketRequest) {
                BadRequest();
                return;
            }

            WebSocket ws = await HttpContext.WebSockets.AcceptWebSocketAsync();
            var buf = new byte[1024 * 4];

            await Setup(ws, buf);
            await ws.CloseAsync(WebSocketCloseStatus.NormalClosure, "Success", CancellationToken.None);
        }

        private async Task Setup(WebSocket ws, byte[] buf) {
            // Read setup msg
            WebSocketReceiveResult result = await ws.ReceiveAsync(new ArraySegment<byte>(buf), CancellationToken.None);
            SetupMsg setupMsg = JsonConvert.DeserializeObject<SetupMsg>(Encoding.UTF8.GetString(buf, 0, result.Count));

            if (setupMsg.FileSize > _config.MaxFileSize) {
                await SendResp(ws, Resp.BadRequest("File size too large."));
                throw new Exception("File size too large");
            }

            setupMsg.FileName = UploadUtil.SanitizeFileName(setupMsg.FileName);
            string code = GetGuid();
            
        }

        private string GetGuid() {
            Guid code = Guid.NewGuid();
            return code.ToString();
        }

        private static async Task SendResp(WebSocket ws, Resp resp) {
            byte[] respBytes = Encoding.UTF8.GetBytes(JsonConvert.SerializeObject(resp));
            await ws.SendAsync(new ArraySegment<byte>(respBytes), WebSocketMessageType.Text, true,
                CancellationToken.None);
        }
        
        [JsonObject(ItemRequired = Required.Always)]
        private class SetupMsg {
            public long FileSize { get; set; }
            public string FileName { get; set; }
        }
    }
}