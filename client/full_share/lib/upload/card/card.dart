part of upload;

typedef void IDCallback(int id);

class Card extends material.StatelessWidget {
  final int id;

  material.Widget get _target => new material.Container(
        margin: const material.EdgeInsets.symmetric(horizontal: 10.0),
        child: new material.Card(
          child: new material.Container(
            child: _mainRow,
//          margin: const material.EdgeInsets.only(left: 10.0, right: 10.0, bottom: 10.0),
            padding: const material.EdgeInsets.only(top: 10.0, bottom: 10.0, left: 10.0, right: 5.0),
//        decoration: new material.BoxDecoration(
//          color: style.materialBlueGrey100, // material bluegrey100
////          border: new material.Border.all(color: new material.Color(0xFF000000)),
//          borderRadius: new material.BorderRadius.all(new material.Radius.circular(20.0)),
//        ),
          ),
        ),
      );

  material.Widget get _mainRow => new material.Row(children: [_imageContainer, _codeAndExpiresCol, _removeWidget]);

  material.Widget get _imageContainer => new material.Container(
        height: 90.0,
        width: 60.0,
        margin: const material.EdgeInsets.only(right: 5.0),
        color: style.translucentCrystalBlue,
      );

  material.Widget get _codeAndExpiresCol => new material.Column(
        children: [_codeWidget, _expiresWidget],
        crossAxisAlignment: material.CrossAxisAlignment.start,
      );

  material.Widget get _codeWidget => new material.Row(
        children: <material.Widget>[
          new material.Text("Code: "),
          new material.Text(
            "TEMPPP",
            style: new material.TextStyle(fontSize: 20.0, color: material.Colors.green),
          ),
        ],
      );

  material.Widget get _expiresWidget => new material.Container(
        child: new material.Text("Expires: 00/00/0000 00:00pm id: $id"),
      );

  material.Widget get _removeWidget => new material.Expanded(
        child: new material.Container(
          child: new material.Align(
            alignment: material.Alignment.centerRight,
            child: new material.IconButton(
              icon: new material.Icon(material.Icons.close),
              onPressed: _handleRemoveWidgetPressed,
            ),
          ),
        ),
      );

  void _handleRemoveWidgetPressed() {
    if (removeBtnClickedCallback != null) {
      print('calling callback');
      removeBtnClickedCallback(id);
    }
  }

  final IDCallback removeBtnClickedCallback;

  Card(this.id, {String code = "", this.removeBtnClickedCallback});

  @override
  material.Widget build(material.BuildContext context) {
    return new material.Center(child: _target);
  }
}
