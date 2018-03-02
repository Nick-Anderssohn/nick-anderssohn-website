part of download;

class Page extends material.StatelessWidget {
  @override
  material.Widget build(material.BuildContext context) => new material.Container(
//          decoration: new material.BoxDecoration(
//              color: material.Colors.lightBlueAccent,
//              borderRadius: new material.BorderRadius.all(new material.Radius.circular(5.0))),
        padding: new material.EdgeInsets.all(20.0),
        margin: new material.EdgeInsets.symmetric(horizontal: 10.0),
        child: new material.Wrap(
          children: <material.Widget>[
            new material.Text(
              'Download Code:',
              style: new material.TextStyle(fontSize: 30.0, fontWeight: material.FontWeight.bold),
            ),
            new material.TextField(
              decoration: new material.InputDecoration(hintText: 'Enter Code'),
            ),
            new material.Container(
              margin: const material.EdgeInsets.only(top: 10.0),
              child: new material.Center(
                child: new material.RaisedButton(
                  onPressed: () => print('Ya clicked me'),
                  child: new material.Text(
                    'Download',
                    style: new material.TextStyle(color: material.Colors.white),
                  ),
                  color: material.Colors.blue,
                ),
              ),
            ),
          ],
        ),
      );
}
