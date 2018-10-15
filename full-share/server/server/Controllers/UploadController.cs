using System;
using System.Net.WebSockets;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using Newtonsoft.Json;
using server.Config;

namespace server.Controllers {
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
        }

        private static async Task SendResp(WebSocket ws, Resp resp) {
            byte[] respBytes = Encoding.UTF8.GetBytes(JsonConvert.SerializeObject(resp));
            await ws.SendAsync(new ArraySegment<byte>(respBytes), WebSocketMessageType.Text, true,
                CancellationToken.None);
        }

        private class SetupMsg {
            [JsonProperty(Required = Required.Always)]
            public long FileSize { get; set; }

            [JsonProperty(Required = Required.Always)]
            public string FileName { get; set; }
        }

        private class Resp {
            public int StatusCode { get; set; }
            public string StatusMsg { get; set; } = "";
            public string Message { get; set; } = "";

            public static Resp BadRequest() {
                return new Resp {
                    StatusCode = 400,
                    StatusMsg = "Bad Request"
                };
            }

            public static Resp BadRequest(string message) {
                Resp resp = BadRequest();
                resp.Message = message;
                return resp;
            }

            public static Resp Ok() {
                return new Resp {
                    StatusCode = 200,
                    StatusMsg = "Ok"
                };
            }

            public static Resp Ok(string message) {
                Resp resp = Ok();
                resp.Message = message;
                return resp;
            }
        }
    }
}