using System.Linq;

namespace server.Controllers {
    public class UploadUtil {
        internal static string SanitizeFileName(string fileName) {
            char[] chars = fileName.ToCharArray()
                .Select(c => char.IsLetterOrDigit(c) || c == '.' ? c : '_')
                .ToArray();
            return new string(chars);
        }
    }
}