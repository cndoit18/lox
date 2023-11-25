var a = "global\n";
{
  func showA() {
    print a;
  }

  showA();
  var a = "block\n";
  showA();
  print a;
}