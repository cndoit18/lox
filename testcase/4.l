var x = 4;
print("x=" + x);
{
    var x = 3;
    print("inner x=" + x);
}
print("outer x=" + x);

{
    x = 3;
    print("inner x=" + x);
}
print("outer x=" + x);