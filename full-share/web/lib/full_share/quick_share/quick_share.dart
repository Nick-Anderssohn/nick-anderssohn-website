library quickshare;

import 'package:clippy/browser.dart' as clippy;
import 'package:full_share/extended_html/extended_html.dart';
import 'package:full_share/full_share/api_consumer/api_consumer.dart';
import 'package:full_share/error/error.dart';

class QuickShare {
  InputElement _fileInput;
  AnchorElement _copyBtn;
  AnchorElement _sendBtn;
  DivElement _progContainer;
  AnchorElement _shareBtn;
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
    _fileInput.onChange.listen(([_]) {
      if (_fileInput.files.length > 1) {
        window.alert('Please only select 1 file for upload.');
        return;
      }
      if (_fileInput.files.isNotEmpty) {
        _uploadViaWebsocket(_fileInput.files.first);
      }
    });

    _shareBtn.onClick.listen((_) => _fileInput.click());
    _copyBtn.onClick.listen((_) {
      clippy.write(_linkElem.text);
      window.alert('Copied to clipboard.');
    });
  }

  void _uploadViaWebsocket(File file) {
    var uploader = WsUploader.fromFile(file, onErrorMsg: _handleErrorMsg);
    showElem(_progContainer);
    uploader
        .readAndUpload()
        .then(_handleWsUploadResp)
        .catchError(_handleGenericError)
        .whenComplete(() => hideElem(_progContainer));
  }

  void _handleGenericError(dynamic error, [dynamic stackTrace]) {
    print('$error: $stackTrace');
    if (error is UserFacingMessage && error.hasUserFacingMessage()) {
      window.alert(error.getUserFacingMessage());
    }
  }

  void _handleErrorMsg(Map errorMsg) {
    hideElem(_progContainer);
    window.alert('Upload failed');
  }

  void _handleWsUploadResp(String downloadUrl) {
    showElem(_successfulContainer);
    _linkElem.text = downloadUrl;
    _sendBtn.href = _getSMSHref();
  }

  String _getSMSHref() {
    String href = 'sms:&body=${Uri.encodeComponent(_linkElem.text)}';

    if (isAndroid()) {
      href = 'sms:?body=${Uri.encodeComponent(_linkElem.text)}';
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
