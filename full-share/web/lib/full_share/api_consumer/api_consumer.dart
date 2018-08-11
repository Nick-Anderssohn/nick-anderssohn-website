library api_consumer;

import 'dart:async';
import 'dart:convert';
import 'dart:typed_data';

import 'package:full_share/extended_html/extended_html.dart';

class FullShareApiConsumer {
  static const String uploadEndpoint = '/upload';

  static Future<Map> requestUpload(String fileName, Uint8List fileBytes) {
    int fileSize = fileBytes.length;
    if (fileSize > 5.243 * 100000000) {
      window.alert('File size cannot be greater than 500 Mib');
      throw 'File size cannot be greater than 500 Mib';
    }

    return HttpRequest.request(getNewURL(window.location.href, '/', uploadEndpoint), method: 'POST', requestHeaders: {
      'FileName': fileName, 'FileSize': fileBytes.length
    }, sendData: fileBytes).then((HttpRequest resp) => JSON.decode(resp.responseText));
  }
}