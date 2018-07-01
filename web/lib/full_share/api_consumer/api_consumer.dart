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
    return HttpRequest.request('$urlBase$uploadEndpoint', method: 'POST', requestHeaders: {
      'FileName': fileName,
    }, sendData: fileBytes).then((HttpRequest resp) => JSON.decode(resp.responseText));
  }
}