using Newtonsoft.Json;

namespace server.Upload.Controllers.Upload {
    public class Resp {
        public int StatusCode { get; set; }
        public string StatusMsg { get; set; } = "";
        public string Message { get; set; } = "";
        
        [JsonProperty(NullValueHandling = NullValueHandling.Ignore)]
        public long? ValueLong { get; set; } = null;

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

        public static Resp InternalError() {
            return new Resp {
                StatusCode = 500,
                StatusMsg = "Internal Server Error"
            };
        }

        public Resp WithValueLong(long value) {
            ValueLong = value;
            return this;
        }
    }
}