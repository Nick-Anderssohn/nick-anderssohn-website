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
  io.File imageFile;
  material.TabController _tabController;

  final String title;

  int currentTabIndex = 0;

  upload.PageModel uploadPageModel = new upload.PageModel();

  _PageState({this.title = ""});

  void _getImage() async {
    var _fileName = await image_picker.ImagePicker.pickImage();
    setState(() {
      imageFile = _fileName;
    });
  }

  material.FloatingActionButton _getFloatingActionButton() {
    //test
    material.FloatingActionButton btn;
    if (currentTabIndex == uploadTabIndex) {
      btn = new material.FloatingActionButton(
//        onPressed: addCard,
        onPressed: _getImage,
        tooltip: 'Upload File',
        child: new material.Icon(material.Icons.image),
      );
    }

    return btn;
  }

  void _detectTabChange() {
    if (_tabController.index != currentTabIndex) {
      setState(() => currentTabIndex = _tabController.index);
    }
  }

  void addCard() {
    setState(() => uploadPageModel.cardModels[uploadPageModel.nextID] = new upload.CardModel(uploadPageModel.nextID));
    uploadPageModel.nextID++;
  }

  void removeCard(int cardID) => setState(() => uploadPageModel.cardModels.remove(cardID));

  @override
  void initState() {
    super.initState();
    _tabController = new material.TabController(vsync: this, length: myTabs.length)..addListener(_detectTabChange);
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
        title: new material.Text(title),
        bottom: new material.TabBar(
          controller: _tabController,
          tabs: myTabs,
        ),
      ),
      body: new material.TabBarView(
        controller: _tabController,
        children: [new upload.Page(uploadPageModel, removeCard), new download.Page()],
      ),
      floatingActionButton: _getFloatingActionButton(),
    );
  }
}
