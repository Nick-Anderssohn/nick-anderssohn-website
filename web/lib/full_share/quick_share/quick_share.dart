library quickshare;

import '../../extended_html/extended_html.dart';
import 'package:full_share/full_share/api_consumer/api_consumer.dart';

class QuickShare {
  InputElement _fileInput;
  AnchorElement _copyBtn;
  AnchorElement _sendBtn;
  DivElement _progContainer;
  AnchorElement _shareBtn;
  FileReader _fileReader = new FileReader();

  void run() {
    try {
      _queryElements();
      _setupHandlers();
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
    _fileReader.onLoad.listen(_handleLoadFile);
    _fileInput.onChange.listen(([_]) {
      try {
        if (_fileInput.files.length > 1) {
          throw 'Please only select 1 file for upload.';
        }
        if (_fileInput.files.isNotEmpty) {
          _fileReader.readAsArrayBuffer(_fileInput.files.first);
        }
      } catch (e) {
        window.alert(e);
      }
    });

    _shareBtn.onClick.listen((_) => _fileInput.click());
  }

  void _handleLoadFile(ProgressEvent e) {
    FullShareApiConsumer.requestUpload(_fileInput.files.first.name, _fileReader.result).then(print).catchError(window.alert);
  }
}