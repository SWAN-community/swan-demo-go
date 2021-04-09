var loaded = false;

stop = function(s, d, r) {
    loadOWID.then(() => {
        fetch("/stop?" +
            "host=" + encodeURIComponent(d) + "&" +
            "returnUrl=" + encodeURIComponent(r),
            { method: "GET", mode: "cors", cache: "no-cache" })
            .then(r => r.text() )
            .then(m => {
                console.log(m);
                window.location.href = m;
            })
            .catch(x => {
                console.log(x);
            });
    });
}

appendComplaintEmail = function(e, d, o, s, g) {
    loadOWID().then(() => {
        fetch((d ? "//" + d : "") + "/complain?" +
            "offerid=" + encodeURIComponent(o) + "&" +
            "swanowid=" + encodeURIComponent(s),
            { method: "GET", mode: "cors", cache: "no-cache" })
            .then(r => r.text() )
            .then(m => {
                var a = document.createElement("a");
                a.href = m;
                if (g) {
                    var i = document.createElement("img");
                    i.src = g;
                    i.style="width:32px"
                    a.appendChild(i);
                } else {
                    a.innerText = "?";
                }
                e.appendChild(a);
            }).catch(x => {
                console.log(x);
            });
    });
}

appendName = function(e, s) {
    loadOWID().then(() => {
        var supplier = new owid(s);

        fetch("//" + supplier.domain + "/owid/api/v1/creator", 
            { method: "GET", mode: "cors" })
            .then(r => r.json())
            .then(o => {
                var t = document.createTextNode(o.name);
                e.appendChild(t);
            }).catch(x => {
                console.log(x);
                var t = document.createTextNode(x);
                e.appendChild(t);
            });
    });
}

appendAuditMark = function(e, r, t) {

    // Append a failure HTML element.
    function returnFailed(e) {
        var t = document.createTextNode("Failed");
        e.appendChild(t);
    }

    // Append an audit HTML element.
    function addAuditMark(e, r) {
        var t = document.createElement("img");
        if (r) {
            t.src = "/green.svg";
        } else {
            t.src = "/red.svg";
        }
        e.appendChild(t);
    }

    function verify(r, t) {
        if (r === undefined || r === "") {
            var o = new owid(t);
            return o.verify();    
        } else {
            var parent = new owid(r);
            var o = new owid(t); 
            return o.verify(parent);
        }
    }

    // Valid the OWID against the creators public key OR if crypto not 
    // supported the well known end point for OWID creators. OWID providers
    // are not required to operate an end point for verifying OWIDs so these
    // calls might fail to return a result.
    loadOWID().then(() => {
        verify(r, t)
            .then(r => addAuditMark(e, r))
            .catch(x => {
                console.log(x);
                returnFailed(e);
            });
    });
}

loadOWID = function() {
    return new Promise((resolve, reject) => {
        if(loaded || typeof owid !== "undefined") {
            resolve();
        }
        var script = document.createElement('script');
        script.onload = function () {
            loaded = true;
            resolve();
        };
        script.src = "/owid-js/v1.js";
        document.head.appendChild(script);
    });
}