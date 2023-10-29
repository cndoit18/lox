func table() {
    for (var i = 1; i <= 9; i = i + 1) {
        for (var j = 1; j <= 9; j = j + 1) {
            print "" + i + "*" + j + "=" + i * j + "\t";
        }
        print "\n";
    }
}

print "table\n";
table();