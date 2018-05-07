part of home;

const int uploadTabIndex = 0;
const int downloadTabIndex = 1;

class Page extends material.StatefulWidget {
  final String title;
  Page({material.Key key, this.title = ""}) : super(key: key);

  @override
  _PageState createState() => new _PageState(title: title);
}

class _PageState extends material.State<Page> with material.SingleTickerProviderStateMixin {
  final List<material.Tab> myTabs = [new material.Tab(text: 'Upload'), new material.Tab(text: 'Download')];
  material.TabController _tabController;

  final String title;

  int currentTabIndex = 0;

  upload.AlbumPageModel uploadPageModel = new upload.AlbumPageModel();

  bool _newAlbumPageIsShown = false;

  _PageState({this.title = ""});

  void _getImage() async {
    image_picker.ImagePicker.pickImage(source: image_picker.ImageSource.gallery).then(addCard);
  }

  void _showNewAlbumPage() {
//    return material.Navigator.push(context, new material.MaterialPageRoute(builder: (_) => new upload.NewAlbumPage()));
    setState(() => _newAlbumPageIsShown = true);
  }

  material.FloatingActionButton _getFloatingActionButton() {
    //test
    material.FloatingActionButton btn;
    if (currentTabIndex == uploadTabIndex) {
      btn = new material.FloatingActionButton(
//        onPressed: _showNewAlbumPage,
//        onPressed: addCard,
        onPressed: _getImage,
        tooltip: 'New Album',
        child: new material.Icon(material.Icons.add),
      );
    }

    return btn;
  }

  void _detectTabChange() {
    if (_tabController.index != currentTabIndex) {
      setState(() => currentTabIndex = _tabController.index);
    }
  }

  void addCard([io.File imageFile]) {
    setState(() => uploadPageModel.cardModels[uploadPageModel.nextID] = new upload.CardModel(uploadPageModel.nextID)
      ..imageFile = imageFile);
    uploadPageModel.nextID++;
  }

  void removeCard(int cardID) => setState(() => uploadPageModel.cardModels.remove(cardID));

  @override
  void initState() {
    super.initState();
    _tabController = new material.TabController(vsync: this, length: myTabs.length)..addListener(_detectTabChange);
  }

  void _popNewAlbumPage() {
    setState(() => _newAlbumPageIsShown = false);
  }

  material.Widget _getLeadingIcon() {
    material.Widget _leadingIcon;
    if (_newAlbumPageIsShown) {
      _leadingIcon =
          new material.IconButton(icon: new material.Icon(material.Icons.arrow_back), onPressed: _popNewAlbumPage);
    }
    return _leadingIcon;
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  material.Widget build(material.BuildContext context) {
    return new material.Scaffold(
      appBar: new material.AppBar(
        leading: _getLeadingIcon(),
        title: new material.Text(title),
        bottom: new material.TabBar(
          controller: _tabController,
          tabs: myTabs,
        ),
      ),
      body: new material.TabBarView(
        controller: _tabController,
        children: [new upload.AlbumPage(uploadPageModel, removeCard, _newAlbumPageIsShown), new download.Page()],
      ),
      floatingActionButton: _getFloatingActionButton(),
    );
  }
}
