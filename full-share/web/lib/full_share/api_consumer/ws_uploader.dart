part of api_consumer;

/*
A successful save file hand shake will look like the following:
Client sends:
    {
        "FileSize": 1000,
        "FileName": "Chicken Dinner"
    }

Server replies:
	{
        "Status": 200,
        "Message": ""
  }

Then in a loop until file size is reached:
	Client sends file bytes.

	Server replies with:
		{
			"Status": 200,
			"Message": ""
    	}

Then the server will finish up with:
	{
        "Status": 200,
        "Message": "The url here"
    }

At any point in time, the server might respond with status indicating an error code and message containing an
error message.
*/

/// Uploads files to full share through it's websocket api.
/// Each instance is a one time use that automatically cleans itself up after uploading or failing to upload.
/// If a new upload is needed, a new [WsUploader] must be created.
class WsUploader extends SubCleaner {
  static const String uploadEndpoint = '/upload';
  static const int _sliceSize = 1024 * 1024 * 20; // 20 MiB

  final File _file;
  WebSocket _ws;
  final Completer<String> _downloadLinkCompleter = Completer();
  int _bytesSent = 0;
  bool _connected = false;
  Function onErrorMsg; // Must accept a map as it's only argument

  /// Create a [WsUploader] set to read and upload [_file].
  WsUploader.fromFile(this._file, {this.onErrorMsg}) {
    String url = '';
    if (window.location.protocol == 'https:') {
      url = getNewURL(window.location.href, '/', uploadEndpoint).replaceFirst('https', 'wss');
    } else if (window.location.protocol == 'http:') {
      url = getNewURL(window.location.href, '/', uploadEndpoint).replaceFirst('http', 'ws');
    }
    _ws = WebSocket(url);
    addSub(_ws.onOpen.listen((_) => _connected = true));
    addSub(_ws.onError.listen((e) => window.alert(e.toString())));
  }

  /// Reads [_file] and simultaneously uploads it to full share.
  Future<String> readAndUpload() async {
    _waitForOpenConn()
        .then(_readAndUpload)
        .catchError((e) => _downloadLinkCompleter.completeError(e, StackTrace.current));
    return _downloadLinkCompleter.future;
  }

  /// Waits a max of [maxSecondsToWait] for the websocket connection to open. If it fails to open,
  /// a [TimeoutException] is thrown.
  Future<Null> _waitForOpenConn({int maxSecondsToWait = 3}) async {
    int numMilliPerTry = 50;
    int maxTries = maxSecondsToWait * Duration.millisecondsPerSecond ~/ numMilliPerTry;
    for (int i = 0; !_connected; i++) {
      await Future.delayed(Duration(milliseconds: 50));

      if (i >= maxTries) {
        throw LeveledException.OpsError('upload websocket failed to open')
          ..attachedStackTrace = StackTrace.current
          ..userFacingMessage = 'Upload failed. Could not connect to server.';
      }
    }
  }

  /// Reads [_file] and simultaneously uploads it to full share.
  void _readAndUpload(_) {
    addSub(_ws.onMessage.listen(_handleServerResponse));
    _startUpload();
  }

  /// Handles the first handshake with the server. Sets up the upload.
  void _startUpload() {
    String firstMessage = jsonEncode({
      'FileSize': _file.size,
      'FileName': _file.name,
    });

    _ws.sendString(firstMessage);
  }

  /// Checks the response from the server and either triggers the next slice upload or completes the overall upload.
  void _handleServerResponse(MessageEvent msgEvent) {
    Map serverResponse = jsonDecode(msgEvent.data);
    if (_UploadRespParseUtil.getStatusCode(serverResponse) == _Status.ok) {
      if (_doneUploading()) {
        String link = _UploadRespParseUtil.getMessage(serverResponse);
        if (link.isNotEmpty) {
          _downloadLinkCompleter.complete(link);
        }
      } else {
        _sendNextFileSlice();
      }
    } else {
      print(serverResponse);
      cleanup();
      onErrorMsg?.call(serverResponse);
    }
  }

  /// Sends the next slice of the file over the websocket.
  void _sendNextFileSlice() {
    int bytesToSend = _bytesSent + _sliceSize > _file.size ? _file.size - _bytesSent : _sliceSize;
    Blob slice = _file.slice(_bytesSent, _bytesSent + bytesToSend);
    _bytesSent += bytesToSend;

    _ws.sendBlob(slice);
  }

  bool _doneUploading() => _bytesSent >= _file.size;

  /// Cancels all [StreamSubscription]s and closes the websocket connection.
  void cleanup() {
    cancelAllSubs();
    _ws.close();
  }
}

/// Contains static functions for parsing upload responses
class _UploadRespParseUtil {
  static const String statusCodeKey = 'StatusCode';
  static const String statusMsgKey = 'StatusMsg';
  static const String messageKey = 'Message';

  static int getStatusCode(Map serverResponse) => parseJsonFromKey(serverResponse, statusCodeKey);

  static String getStatusMsg(Map serverResponse) => parseJsonFromKey(serverResponse, statusMsgKey);

  static String getMessage(Map serverResponse) => parseJsonFromKey(serverResponse, messageKey);
}

/// Contains status codes used in responses from the server
class _Status {
  static const int ok = 200;
}
