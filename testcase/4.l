var x = 4;
print("x=" + x + "\n");
{
    var x = 3;
    print("inner x=" + x + "\n");
}
print("outer x=" + x + "\n");

{
    x = 3;
    print("inner x=" + x + "\n");
}
print("outer x=" + x + "\n");