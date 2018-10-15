using System;
using System.Net.WebSockets;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;

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
        [HttpGet]
        public async Task<ActionResult> HandleWsUpload() {
            // Setup ws connection
            if (!HttpContext.WebSockets.IsWebSocketRequest) {
                return BadRequest();
            }

            var respBytes = Encoding.UTF8.GetBytes("test");
            WebSocket ws = await HttpContext.WebSockets.AcceptWebSocketAsync();

            var buf = new byte[1024 * 4];
            
            // Read msg from client
            await ws.ReceiveAsync(new ArraySegment<byte>(buf), CancellationToken.None);

            // Reply with test
            // msg content, msg type, end of msg, cancellation token
            await ws.SendAsync(new ReadOnlyMemory<byte>(respBytes), WebSocketMessageType.Text, true, CancellationToken.None);
            await ws.CloseAsync(WebSocketCloseStatus.NormalClosure, "Success", CancellationToken.None);
            return Ok();
        }
    }
}