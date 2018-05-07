part of upload;

typedef void IDCallback(int id);

class AlbumCard extends material.StatelessWidget {
  final int id;

  material.Widget get _target => new material.Container(
        margin: const material.EdgeInsets.symmetric(horizontal: 10.0),
        child: new material.Card(
          child: new material.Container(
            child: new material.Column(
              children: <material.Widget>[
                _mainRow,
                _buttonBar,
              ],
            ),
            padding: const material.EdgeInsets.only(top: 20.0, left: 10.0, right: 5.0),
          ),
        ),
      );

  material.Widget get _buttonBar => new material.ButtonTheme.bar(
        child: new material.ButtonBar(
          children: <material.Widget>[
            new material.FlatButton(
              child: const material.Text('View'),
              onPressed: () => null,
            )
          ],
        ),
      );

  material.Widget get _mainRow => new material.Row(children: [
        _albumIcon,
        /*_codeAndExpiresCol,*/ _titleWidget, /*_removeWidget*/
      ]);

  material.Widget get _albumIcon => new material.Padding(
        padding: const material.EdgeInsets.only(right: 16.0),
        child: const material.Icon(
          material.Icons.album,
          color: material.Colors.black45,
        ),
      );

  material.Widget get _codeAndExpiresCol => new material.Column(
        children: [_titleWidget, _expiresWidget],
        crossAxisAlignment: material.CrossAxisAlignment.start,
      );

  material.Widget get _titleWidget => new material.Row(
        children: <material.Widget>[
          new material.Text(
            "TEMPPP",
            style: new material.TextStyle(fontSize: 20.0, color: material.Colors.black87),
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
      removeBtnClickedCallback(id);
    }
  }

  final IDCallback removeBtnClickedCallback;

  AlbumCard(this.id, {String code = "", this.removeBtnClickedCallback});

  @override
  material.Widget build(material.BuildContext context) {
    return new material.Center(child: _target);
  }
}
