// NOTE: This is in an ops folder because of a weird combination between nginx and golang
// hosting...
import 'package:nick_ops/extended_html/extended_html.dart';
import 'dart:convert';

AnchorElement _loginBtn = querySelector('#login-btn');
InputElement _usernameInput = querySelector('#username');
InputElement _passwordInput = querySelector('#password');

void main() {
  _setupLoginSubmission();
}

void _setupLoginSubmission() {
  _loginBtn.onClick.listen(_requestLogin);
}

void _requestLogin([_]) {
  HttpRequest.request(getNewURL(window.location.href, "/ops/Login", "/ops/SubmitLogin"),
      method: 'POST',
      sendData: jsonEncode({
        'username': _usernameInput.value,
        'password': _passwordInput.value,
      })).then(_handleLoginResponse);
}

void _handleLoginResponse(HttpRequest request) {
  print('status: ${request.status}: ${request.statusText}');
  if (request.status == 200) {
    window.location.href = getNewURL(window.location.href, '/ops/Login', '/ops/Kibana');
  }
}
