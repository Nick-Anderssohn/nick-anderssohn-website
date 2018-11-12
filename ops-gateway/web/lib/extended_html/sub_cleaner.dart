part of extended_html;

class SubCleaner {
  List<StreamSubscription> subs = [];

  void cancelAllSubs() => subs.forEach((sub) => sub.cancel());
  void addSub(StreamSubscription sub) => subs.add(sub);
}
