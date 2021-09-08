// JavaScript used only in the SWAN demo. Not for production use.

var loaded = false;

stop = function(s, d, r, a) {
    var data = new URLSearchParams();
    data.append("host", d);
    data.append("returnUrl", r);
    data.append("accessNode", s);
    fetch("/stop?",
        { 
            method: "POST", 
            mode: "cors", 
            cache: "no-cache",
            body: data 
        })
        .then(r => r.text() )
        .then(m => {
            console.log(m);
            window.location.href = m;
        })
        .catch(x => {
            console.log(x);
        });
}

appendComplaintEmail = function(e, d, o, s, g) {
    var data = new URLSearchParams();
    data.append("swanid", o);
    data.append("partyid", s);
    fetch((d ? "//" + d : "") + "/complain?",
        { 
            method: "POST", 
            mode: "cors", 
            cache: "no-cache",
            body: data 
        })
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
}

appendName = function(e, s) {
    var supplier = new owid(s);
    var url = "//" + supplier.domain + "/owid/api/v1/creator";
    fetch(url, 
        { method: "GET", mode: "cors" })
        .then(r => {
            if (r.ok) {
                return r.json();
            }
            throw "Bad response from '" + url + "'";
        })
        .then(o => {
            var t = document.createTextNode(o.name);
            e.appendChild(t);
        }).catch(x => {
            console.log(x);
            var t = document.createTextNode(x);
            e.appendChild(t);
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
    verify(r, t)
        .then(r => addAuditMark(e, r))
        .catch(x => {
            console.log(x);
            returnFailed(e);
        });
}