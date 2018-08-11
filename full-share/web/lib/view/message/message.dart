part of view;

/// A const class used to display a message. Add [target] to the DOM to display.
class Message {
  DivElement get target => new DivElement()..classes.add('message')..innerHtml =
  '''
    <span class="message-icon-holder"">${showIcon ? MaterialIcons.smallAttachFile : ""}</span>
    <div class="message-text-area-container">
    <div>$senderName &gt $recipientName</div>
    <div>$message</div>
    </div>
    ''';

  final String senderName;
  final String recipientName;
  final String message;
  final String fileName;
  final bool showIcon;

  const Message(this.senderName, this.recipientName, this.message, this.fileName, {this.showIcon = true});
}