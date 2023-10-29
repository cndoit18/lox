var x = true;
var y = false;

if(x) {
    print "outer if";
    if(y) {
        print "inner if";
    } else {
        print "inner else";
    }
} else {
    print "outer else";
}