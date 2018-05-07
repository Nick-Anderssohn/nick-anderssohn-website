part of upload;

class AlbumPage extends material.StatelessWidget {
  final IDCallback removeCardCallback;
  final AlbumPageModel model;
  final bool showNewAlbumUI;

  AlbumPage(this.model, this.removeCardCallback, this.showNewAlbumUI);

  @override
  material.Widget build(material.BuildContext context) {
//    List<material.Widget> listChildren = model.cardModels
//        .map((cardModel) => new Card(
//              cardModel.cardID,
//              removeBtnClickedCallback: removeCardCallback,
//            ))
//        .toList();

    if (showNewAlbumUI) {
      return new material.Center(
        child: new material.Text('CHICKEN'),
      );
    }

    List listChildren = model.cardModels.values.toList().map((cardModel) {
      material.Widget val = new ImageCard(
        cardModel.cardID,
        removeBtnClickedCallback: removeCardCallback,
        imageFile: cardModel.imageFile,
      );
      return val;
    }).toList();

    // At this point, the only values should be of type AlbumCard
    listChildren.sort((c1, c2) => (c1 as ImageCard).id < (c2 as ImageCard).id ? 1 : -1);

    // add space so that you can scroll past end a bit
    listChildren.add(new material.Container(
      height: 180.0,
    ));

    int itemCount = listChildren.length;

    return new material.ListView.builder(
      itemBuilder: (context, i) {
        if (i < itemCount) {
          return listChildren[i];
        }
        return null;
      },
    );
  }
}
