library quickshare;

import '../../extended_html/extended_html.dart';

class QuickShare {
  InputElement _fileInput;
  AnchorElement _copyBtn;
  AnchorElement _sendBtn;
  DivElement _progContainer;
  AnchorElement _shareBtn;

  void run() {
    try {
      _queryElements();
    } catch (e) {
      print(e);
      // add a second param s if stacktrace is desired
      // print(s); // prints stacktrace
    }
  }

  /// Queries all elements from DOM. Will throw an error message if an element could not be queried.
  void _queryElements() {
    _fileInput = tryQuerySelector('#file-input');
    _copyBtn = tryQuerySelector('#copy-btn');
    _sendBtn = tryQuerySelector('#send-btn');
    _progContainer = tryQuerySelector('#prog-container');
    _shareBtn = tryQuerySelector('#share-btn');
  }

  void _setupHandlers() {
    
  }
}