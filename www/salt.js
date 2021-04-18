function SWANSalt(element) {
            
    const gridId = 'salt-grid';
    const imgIdPrefix = 'salt-';
    const colIdPrefix = 'salt-col-';
    const event = new Event('complete');

    var selected = [];
    var indicators = [];

    function stringValue() {
        var binary = "";
        var v = byteValue();
        var len = v.byteLength;
        for (var i = 0; i < len; i++) {
            binary += String.fromCharCode(v[i]);
        }
        return btoa(binary).replace(/\=/, '');
    }

    function byteValue() {
        if(selected.length == 4){
            var b1 = (selected[0] << 4) | selected[1];
            var b2 = (selected[2] << 4) | selected[3];

            return new Uint8Array([b1, b2]);
        }
    }

    function value() {
        return selected.join('-');
    }

    function onComplete(callBack) {
        element.addEventListener('complete', function(e) {
            if (typeof callBack == 'function') {
                callBack();
            }
        })
    }

    function complete() {
        var grid = document.getElementById(gridId);
        grid.classList.add('complete');
        element.dispatchEvent(event);
    }

    function reset(callBack) {
        var grid = document.getElementById(gridId);
        grid.classList.remove('complete');
        selected = [];
        indicators.forEach(i => {
            i.parentNode.removeChild(i);
        });
        indicators = [];
        if (typeof callBack == 'function') {
            callBack();
        }
    }

    function getIndicatorClass(value){
        switch (value) {
            case 1:
                return "top-left";
            case 2:
                return "top-right";
            case 3:
                return "bottom-left";
            case 4:
                return "bottom-right";
            default:
                throw `Invalid value: ${value}`;
        }
    }

    function updateIndicator(value) {
        var element = document.getElementById(colIdPrefix + value);
        var item = selected.length;
        
        var indicator = document.createElement('div');
        indicator.classList.add(getIndicatorClass(item));
        indicator.innerHTML = item;
        
        indicators.push(indicator);
        element.appendChild(indicator);

    }

    function add(event) {
        var value = parseInt(event.target.getAttribute("data-value"));
        
        if (selected.length < 4) {
            selected.push(value);
            updateIndicator(value);
        }
        if (selected.length == 4){
            complete();
        }
    }

    async function salt() {
        var animals = await fetch('/animals.json').then(res => res.json());

        var txt = `<div id="${gridId}" class="row">`;
        for (var i = 0; i < animals.length && i < 16; i++) {
            var animal = animals[i];

            txt += `<div id="${colIdPrefix + i}" class="col"><img id="${imgIdPrefix + i}" src="/${animal.Filename}" data-value="${i}"></div>`;

            if ((i + 1) % 4 == 0) {
                txt += '<div class="w-100"></div>';
            }
        }
        txt += '</div>';
        txt += `<small id="saltNote" class="form-text text-muted">
            This implementation needs to be modified to support screen readers and ARIA prior to production use. 
            It is provided for conceptual demonstration purposes only at this time.
        </small>
        <small id="saltNote" class="form-text text-muted">
            Icons provided by the Noun Project under creative commons licence.
        </small>
        `;
        element.innerHTML = txt;

        for (var i = 0; i < 16; i++) {
            var animal = document.getElementById(imgIdPrefix + i);

            animal.addEventListener('click', add);
        }
    }

    //#region public methods and getters

    Object.defineProperty(this, 'value', { 
        get: function() { return value(); } 
    });

    Object.defineProperty(this, 'byteValue', { 
        get: function() { return byteValue(); } 
    });

    Object.defineProperty(this, 'stringValue', { 
        get: function() { return stringValue(); } 
    });
    
    this.reset = function (callBack) {
        reset(callBack);
    }

    this.onComplete = function(callBack) {
        onComplete(callBack);
    }

    //#endregion

    // generate grid;
    salt();
}