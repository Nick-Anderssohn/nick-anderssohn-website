library quickshare;

import '../../extended_html/extended_html.dart';
import 'package:clippy/browser.dart' as clippy;
import 'package:full_share/full_share/api_consumer/api_consumer.dart';

class QuickShare {
  InputElement _fileInput;
  AnchorElement _copyBtn;
  AnchorElement _sendBtn;
  DivElement _progContainer;
  AnchorElement _shareBtn;
  FileReader _fileReader = new FileReader();
  DivElement _successfulContainer;
  ParagraphElement _linkElem;

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
    _successfulContainer = tryQuerySelector('#successful-container');
    _linkElem = tryQuerySelector('#link');
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
    _copyBtn.onClick.listen((_) => clippy.write(_linkElem.text));
  }

  void _handleLoadFile(ProgressEvent e) {
    showElem(_progContainer);
    FullShareApiConsumer.requestUpload(_fileInput.files.first.name, _fileReader.result).then(_handleUploadResp).catchError(print).whenComplete(([_]) => hideElem(_progContainer));
  }

  void _handleUploadResp(Map responseJson) {
    if (responseJson['Success']) {
      showElem(_successfulContainer);
      _linkElem.text = responseJson['DownloadLink'];
      _sendBtn.href = _getSMSHref();
    }
  }

  String _getSMSHref() {
    String href = 'sms:&body=${Uri.encodeFull(_linkElem.text)}';

    if (isAndroid()) {
      href = 'sms:?body=${Uri.encodeFull(_linkElem.text)}';
    }

    return href;
  }

  bool isAndroid() {
    var mua = window.navigator.userAgent;
    if (mua == null) {
      mua = window.navigator.vendor;
    }
    mua = mua?.toLowerCase();
    return mua?.contains('android');
  }
}