part of upload;

class Page extends material.StatelessWidget {
  final IDCallback removeCardCallback;
  final PageModel model;

  Page(this.model, this.removeCardCallback);

  @override
  material.Widget build(material.BuildContext context) {
//    List<material.Widget> listChildren = model.cardModels
//        .map((cardModel) => new Card(
//              cardModel.cardID,
//              removeBtnClickedCallback: removeCardCallback,
//            ))
//        .toList();

    List listChildren = model.cardModels.values
        .toList()
        .map((cardModel) => new Card(
              cardModel.cardID,
              removeBtnClickedCallback: removeCardCallback,
            ))
        .toList();
    listChildren.sort((c1, c2) => c1.id < c2.id ? 1 : -1);

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
