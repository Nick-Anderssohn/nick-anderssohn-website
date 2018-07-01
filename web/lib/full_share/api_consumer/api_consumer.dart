library api_consumer;

import 'dart:async';
import 'dart:convert';
import '../../extended_html/extended_html.dart';
import 'dart:typed_data';

class FullShareApiConsumer {
  static const String uploadEndpoint = '/upload';

  static Future<Map> requestUpload(String fileName, Uint8List fileBytes) {
    int fileSize = fileBytes.length;
    if (fileSize > 6.554 * 10000000) {
      window.alert('File size cannot be greater than 500 Mib');
      throw 'File size cannot be greater than 500 Mib';
    }

    return HttpRequest.request(getNewURL(window.location.href, '/', uploadEndpoint), method: 'POST', requestHeaders: {
      'FileName': fileName, 'FileSize': fileBytes.length
    }, sendData: fileBytes).then((HttpRequest resp) => JSON.decode(resp.responseText));
  }
}