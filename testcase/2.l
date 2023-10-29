var x = true;
var y = false;

if(x) {
    print "outer if" + "\n";
    if(y) {
        print "inner if" + "\n";
    } else {
        print "inner else" + "\n";
    }
} else {
    print "outer else" + "\n";
}

print "hi" or 2; // "hi".
print "\n";
print nil or "yes"; // "yes".