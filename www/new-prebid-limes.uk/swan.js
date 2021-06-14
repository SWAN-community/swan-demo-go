window.swan = function (window) {
  var swan_owids = {};
  return {
    addSignature: function(placement, owid) {
      console.log(placement, owid)
      if (!swan_owids[placement]) {
        swan_owids[placement] = [];
      }
      swan_owids[placement].push(owid)
    },
    adInfo: function(placement) {
      var query = []
      var owid = {OWID: window.swanId, Children: swan_owids[placement], Value: null};
      transform_owid = function(value) {
        value = value.replace(/=/g, '')
        return encodeURIComponent(value)
      }
      find_owids = function (tree) {
        if( tree !== null && typeof tree == "object" ) {
            Object.entries(tree).forEach(([key, value]) => {
                if (key === "OWID") {
                  value = transform_owid(value)
                  query.push(value);
                } else if (key === "Children" && value != null) {
                  value.forEach(e => find_owids(e));
                }
            });
        }
      }
      find_owids(owid)
      query.unshift(transform_owid(window.swanId))
      console.log(query)
      var result = [...new Set(query)];;
      var url = "http://cmp.swan-demo.uk/info?owid=" + result.join("&owid=");
      window.open(url, "_blank");
    }
  }
}(window)