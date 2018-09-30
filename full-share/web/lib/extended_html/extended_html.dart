library extended_html;

import 'dart:html';
import 'dart:async';
import 'package:full_share/error/error.dart';
export 'dart:html';

part 'sub_cleaner.dart';

/// [querySelector] except throws a [String] if the element could not be found.
Element tryQuerySelector(String selectors) {
  var elem = querySelector(selectors);
  if (elem == null) {
    throw 'Could not query $selectors';
  }
  return elem;
}

/// Shows an element by setting display = ""
/// Will throw an exception if elem is null
void showElem(Element elem) => elem.style.display = '';

/// HIdes an element by setting display = "none"
/// Will throw an exception if elem is null
void hideElem(Element elem) => elem.style.display = 'none';

String getNewURL(String curURL, String curEndpoint, String newEndpoint) {
  String newURL = curURL;
  var index = newURL.lastIndexOf(curEndpoint);
  newURL = newURL.substring(0, index);
  return newURL + newEndpoint;
}

dynamic parseJsonFromKey(Map jsonMap, dynamic key) {
  var value = jsonMap[key];
  if (value == null) {
    throw LeveledException.Error('unable to parse $key from server response');
  }

  return value;
}
