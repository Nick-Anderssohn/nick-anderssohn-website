library extended_html;
import 'dart:html';
export 'dart:html';

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