library api_consumer;

import 'dart:async';
import 'dart:convert';
import 'dart:html';
import 'dart:typed_data';

class FullShareApiConsumer {
//  static const String urlBase = 'http://nickanderssohn.com';
  static const String urlBase = 'http://localhost:8080';
  static const String uploadEndpoint = '/upload';

  static Future<Map> requestUpload(String fileName, Uint8List fileBytes) {
    int fileSize = fileBytes.length;
    if (fileSize > 6.554 * 10000000) {
      throw 'File size cannot be greater than 500 Mib';
    }
    return HttpRequest.request('$urlBase$uploadEndpoint', method: 'POST', requestHeaders: {
      'FileName': fileName, 'FileSize': fileBytes.length
    }, sendData: fileBytes).then((HttpRequest resp) => JSON.decode(resp.responseText));
  }
}