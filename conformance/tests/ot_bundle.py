


def test_bundle_filter(O):
    errors = []

    O.addVertex("srcVertex")

    edges = {}
    for i in range(100):
        O.addVertex("dstVertex%d" % i)
        edges["dstVertex%d" % i] = {"val" : i}

    O.addBundle("srcVertex", edges, "related")

    #print list(O.query().V("srcVertex").execute())
    #print list(O.query().V("srcVertex").outgoing("related").execute())
    #for i in O.query().V("srcVertex").outEdge("related").execute():
    #    print i
    #for i in O.query().V("srcVertex").groupBundle().execute():
    #    print i

    count = 0
    for i in O.query().V("srcVertex").outEdge("related").filter("function(x) { return x.val > 50; }").outgoing().execute():
        count += 1

    if count != 49:
        errors.append("Fail: Bundle Filter %s != %s" % (count, 49))
    return errors
