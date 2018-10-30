using System;
using System.Net.WebSockets;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using Newtonsoft.Json;
using server.Upload.Db;
using server.Upload.Db.EntityFramework.Model;
using server.Upload.Util;
using Serilog;

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
        private readonly FullShareDbHelper _dbHelper;

        public UploadController(FullShareDbHelper dbHelper) {
            _dbHelper = dbHelper;
        }

        [HttpGet]
        public async Task HandleWsUpload() {
            try {
                Log.Information("New upload.");
                // Setup ws connection
                if (!HttpContext.WebSockets.IsWebSocketRequest) {
                    BadRequest();
                    return;
                }

                WebSocket ws = await HttpContext.WebSockets.AcceptWebSocketAsync();
                SetupMsg setupMsg = await Setup(ws);
                FileProcessor processor = new FileProcessor(ws, setupMsg.FileSize, setupMsg.Code, setupMsg.FileName);
                await processor.Run();
                Log.Information("Normal exit");
            }
            catch (Exception e) {
                Log.Error("Upload failed {exception}", e);
            }
        }

        private async Task<SetupMsg> Setup(WebSocket ws) {
            var buf = new byte[1024 * 4];
            // Read setup msg
            WebSocketReceiveResult result = await ws.ReceiveAsync(new ArraySegment<byte>(buf), CancellationToken.None);
            SetupMsg setupMsg = JsonConvert.DeserializeObject<SetupMsg>(Encoding.UTF8.GetString(buf, 0, result.Count));
            if (setupMsg.FileSize > UploadConfig.MaxFileSize) {
                await ws.CloseAsync(WebSocketCloseStatus.PolicyViolation, "File size too large", CancellationToken.None);
                throw new Exception("File size too large");
            }
            setupMsg.FileName = UploadUtil.SanitizeFileName(setupMsg.FileName);
            setupMsg.Code = GetGuid();
            _dbHelper.InsertFilesEntry(Files.CreateWithUtcNow(setupMsg.Code, setupMsg.FileName, setupMsg.FileSize));
            await UploadUtil.SendResp(ws, Resp.Ok());
            return setupMsg;
        }

        private string GetGuid() {
            Guid code = Guid.NewGuid();
            return code.ToString();
        }
        
        private class SetupMsg {
            [JsonProperty(Required = Required.Always)]
            public long FileSize { get; set; }
            [JsonProperty(Required = Required.Always)]
            public string FileName { get; set; }
            
            [JsonIgnore]
            public string Code { get; set; }
        }
    }
}