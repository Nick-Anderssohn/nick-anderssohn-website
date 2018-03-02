import 'package:flutter/material.dart';
import 'home/home.dart' as home;

const appName = 'Full Share';

void main() => runApp(new MyApp());

class MyApp extends StatelessWidget {
  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return new MaterialApp(
        title: appName,
        theme: new ThemeData(
          primarySwatch: Colors.blue,
        ),
        home: new home.Page(title: appName));
  }
}
