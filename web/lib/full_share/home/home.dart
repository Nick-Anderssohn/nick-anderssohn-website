library home;

import 'dart:html';
import '../../view/view.dart';

/// Sets up the user's homepage.
class Home {
  DivElement _mainDiv;

  void run() {
    _mainDiv = querySelector('#main-container');

    // Test
    _mainDiv.children.add(const Message('Sender', 'Recipient', 'Test message.', 'Test filename').target);
    _mainDiv.children.add(const Message('Sender', 'Recipient', '''
    Test message. chicken dinner oh yeahhhhhh monkey balls kappa kappa feelsbadman
    ''', 'Test filename', showIcon: false).target);
    _mainDiv.children.add(const Message('Sender', 'Recipient', 'Test message.', 'Test filename').target);
  }
}